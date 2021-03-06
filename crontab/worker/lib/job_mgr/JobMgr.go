package job_mgr

import (
	"context"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	logs "github.com/sirupsen/logrus"
	"github.com/yangqinjiang/mycrontab/worker/lib/job_build"
	"github.com/yangqinjiang/mycrontab/worker/common"
	"github.com/yangqinjiang/mycrontab/worker/lib/config"
	"sync"
	"time"
)
/**
	etcd任务管理器, 监听 etcd 的事件, 组装任务数据, 并推给 scheduler任务调度器
 */
type EtcdJobMgr struct {
	client           *clientv3.Client
	kv               clientv3.KV
	lease            clientv3.Lease
	watcher          clientv3.Watcher
	jobEventPusher   *job_build.CustomJobEventPusher //推送任务事件的类
}

var (
	//单例
	G_EtcdJobMgr *EtcdJobMgr
	onceJobMgr   sync.Once
)

func (jobMgr *EtcdJobMgr) SetJobEventPusher(jobEventPusher *job_build.CustomJobEventPusher) {
	jobMgr.jobEventPusher = jobEventPusher
}

//初始化管理器
func InitEtcdJobMgr() (err error) {
	logs.Info("Before onceJobMgr")
	onceJobMgr.Do(func() {

		logs.Info("EtcdEndpoints = ",config.G_config.EtcdEndpoints)
		//初始化配置
		//读取配置文件
		config := clientv3.Config{
			Endpoints:   config.G_config.EtcdEndpoints, //集群地址
			DialTimeout: time.Duration(config.G_config.EtcdDialTimeout) * time.Microsecond,
		}
		//建立连接
		client, err := clientv3.New(config)

		if err != nil {
			logs.Error("----InitEtcdJobMgr Err:----",err)
			return
		}
		//得到Kv和Lease的API子集
		kv := clientv3.NewKV(client)
		lease := clientv3.NewLease(client)
		watcher := clientv3.NewWatcher(client)
		//赋值单例
		G_EtcdJobMgr = &EtcdJobMgr{
			client:  client,
			kv:      kv,
			lease:   lease,
			watcher: watcher,
		}

		//启动监听任务
		G_EtcdJobMgr.watchJobs()

		//启动监听killer
		G_EtcdJobMgr.watchKiller()
	})
	return
}
//创建任务执行锁
func (jobMgr *EtcdJobMgr) CreateJobLock(jobName string) (jobLock JobLocker) {
	//返回 一把锁
	jobLock = InitJobLock(jobName, jobMgr.kv, jobMgr.lease)
	return
}
//监听任务的变化
func (jobMgr *EtcdJobMgr) watchJobs() (err error) {

	var (
		getResp            *clientv3.GetResponse
		kvpair             *mvccpb.KeyValue
		watchStartRevision int64
		watchChan          clientv3.WatchChan
		watchResp          clientv3.WatchResponse
		watchEvent         *clientv3.Event
	)
	//1,get一下/cron/jobs目录下的所有任务,并且获知当前集群的revision
	if getResp, err = jobMgr.kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {
		return
	}
	//当前有哪些任务
	for _, kvpair = range getResp.Kvs {
		//反序列化json,得到job,并推送保存任务的事件到Scheduler
		jobMgr.jobEventPusher.PushSaveEventToScheduler(string(kvpair.Key),kvpair.Value)
	}

	//2,从该revision向后监听变化事件
	//监听协程
	go func() {
		//从GET时刻的后续版本,开始监听
		watchStartRevision = getResp.Header.Revision + 1 // 监听下一次的
		//监听/cron/jobs/目录的后续变化
		//with-rev, with-prefix
		watchChan = jobMgr.watcher.Watch(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithRev(watchStartRevision), clientv3.WithPrefix())
		//处理监听事件
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT: //任务保存
					jobMgr.jobEventPusher.PushSaveEventToScheduler(string(watchEvent.Kv.Key), watchEvent.Kv.Value)
				case mvccpb.DELETE: //任务被删除了
					jobMgr.jobEventPusher.PushDeleteEventToScheduler(string(watchEvent.Kv.Key))
				default:
					//ingore
				}
			}
		}

	}()

	return
}





//监听强杀任务通知
func (jobMgr *EtcdJobMgr) watchKiller() (err error) {

	var (
		watchChan  clientv3.WatchChan
		watchResp  clientv3.WatchResponse
		watchEvent *clientv3.Event
	)

	//2,从该revision向后监听变化事件
	//监听/cron/killer协程
	go func() {
		//监听/cron/killer目录的后续变化
		//  with-prefix
		watchChan = jobMgr.watcher.Watch(context.TODO(), common.JOB_KILLER_DIR, clientv3.WithPrefix())
		//处理监听事件
		for watchResp = range watchChan {
			for _, watchEvent = range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT: //杀死任务事件
					jobMgr.jobEventPusher.PushKillEventToScheduler(string(watchEvent.Kv.Key))
				case mvccpb.DELETE: //killer标记过期,被自动删除
					//不关心此操作
				}

			}
		}

	}()

	return
}

