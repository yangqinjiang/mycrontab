package worker

import (
	"context"
	"fmt"
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"go.etcd.io/etcd/clientv3"
	"go.etcd.io/etcd/mvcc/mvccpb"
	"sync"
	"time"
)

type JobMgr struct {
	client  *clientv3.Client
	kv      clientv3.KV
	lease   clientv3.Lease
	watcher clientv3.Watcher
}

var (
	//单例
	G_jobMgr *JobMgr
	onceJobMgr     sync.Once
)

//初始化管理器
func InitJobMgr() (err error) {
	onceJobMgr.Do(func() {

		//初始化配置
		//读取配置文件
		config := clientv3.Config{
			Endpoints:   G_config.EtcdEndpoints, //集群地址
			DialTimeout: time.Duration(G_config.EtcdDialTimeout) * time.Microsecond,
		}
		//建立连接
		client, err := clientv3.New(config)

		if err != nil {
			return
		}
		//得到Kv和Lease的API子集
		kv := clientv3.NewKV(client)
		lease := clientv3.NewLease(client)
		watcher := clientv3.NewWatcher(client)
		//赋值单例
		G_jobMgr = &JobMgr{
			client:  client,
			kv:      kv,
			lease:   lease,
			watcher: watcher,
		}

		//启动监听任务
		G_jobMgr.watchJobs()

		//启动监听killer
		G_jobMgr.watchKiller()
	})
	return
}

//监听任务的变化
func (jobMgr *JobMgr) watchJobs() (err error) {

	var (
		getResp            *clientv3.GetResponse
		kvpair             *mvccpb.KeyValue
		job                *common.Job
		watchStartRevision int64
		watchChan          clientv3.WatchChan
		watchResp          clientv3.WatchResponse
		watchEvent         *clientv3.Event
		jobEvent           *common.JobEvent
	)
	//1,get一下/cron/jobs目录下的所有任务,并且获知当前集群的revision
	if getResp, err = jobMgr.kv.Get(context.TODO(), common.JOB_SAVE_DIR, clientv3.WithPrefix()); err != nil {
		return
	}
	//当前有哪些任务
	for _, kvpair = range getResp.Kvs {
		//反序列化json,得到job
		if job, err = common.UnpackJob(kvpair.Value); err == nil {
			//反序列化成功

			jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
			//TODO:是把这个job同步给scheduler(调度协程)
			fmt.Println("当前任务=", jobEvent.Job.Name, jobEvent.Job.Command)
			G_scheduler.PushJobEvent(jobEvent)
		}

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
					if job, err = common.UnpackJob(watchEvent.Kv.Value); err != nil {
						continue
					}
					jobName := common.ExtractJobName(string(watchEvent.Kv.Key))
					//构建一个更新event事件
					fmt.Println("更新任务", jobName)
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_SAVE, job)

				case mvccpb.DELETE: //任务被删除了
					// Delete /cron/jobs/job10
					jobName := common.ExtractJobName(string(watchEvent.Kv.Key))
					job = &common.Job{
						Name: jobName,
					}
					fmt.Println("删除任务", jobName)
					//构建一个删除event
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_DELETE, job)

				}
				//推送给scheduler
				G_scheduler.PushJobEvent(jobEvent)
			}
		}

	}()

	return
}

//创建任务执行锁
func (jobMgr *JobMgr) CreateJobLock(jobName string) (jobLock *JobLock) {
	//返回 一把锁
	jobLock = InitJobLock(jobName, jobMgr.kv, jobMgr.lease)
	return
}

//监听强杀任务通知
func (jobMgr *JobMgr) watchKiller() (err error) {

	var (
		job        *common.Job
		watchChan  clientv3.WatchChan
		watchResp  clientv3.WatchResponse
		watchEvent *clientv3.Event
		jobEvent   *common.JobEvent
		jobName    string
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
					jobName = common.ExtractKillerName(string(watchEvent.Kv.Key))
					job = &common.Job{Name: jobName}
					jobEvent = common.BuildJobEvent(common.JOB_EVENT_KILL, job)
					//推送给scheduler
					G_scheduler.PushJobEvent(jobEvent)
				case mvccpb.DELETE: //killer标记过期,被自动删除
					//不关心此操作
				}

			}
		}

	}()

	return
}
