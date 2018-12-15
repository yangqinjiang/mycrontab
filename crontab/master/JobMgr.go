package master

import (
	"go.etcd.io/etcd/clientv3"
	"time"
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"encoding/json"
	"context"
	"fmt"
)

type JobMgr struct {
	client *clientv3.Client
	kv clientv3.KV
	lease clientv3.Lease
}

var(
	//单例
	G_jobMgr *JobMgr
)

//初始化管理器
func InitJobMgr()(err error)  {
	//初始化配置
	//读取配置文件
	config := clientv3.Config{
		Endpoints:G_config.EtcdEndpoints,//集群地址
		DialTimeout:time.Duration(G_config.EtcdDialTimeout)*time.Microsecond,
	}
	//建立连接
	client,err := clientv3.New(config)

	if err != nil{
		return
	}
	//得到Kv和Lease的API子集
	kv := clientv3.NewKV(client)
	lease := clientv3.NewLease(client)
	//赋值单例
	G_jobMgr = &JobMgr{
		client:client,
		kv:kv,
		lease:lease,
	}
	return
}

//保存任务
func (jobMgr *JobMgr)SaveJob(job *common.Job) (oldJob *common.Job,err error)  {
	//把任务保存到/cron/jobs/任务名 -> json
	var (
		jobKey string
		jobValue []byte
		oldJobObj common.Job
		putResp *clientv3.PutResponse
	)

	//etcd保存的key
	jobKey = common.JOB_SAVE_DIR+job.Name
	//任务信息json
	jobValue,err = json.Marshal(job)
	if err != nil{
		return
	}
	//保存到etcd
	//返回旧值
	putResp,err = jobMgr.kv.Put(context.TODO(),jobKey,string(jobValue),clientv3.WithPrevKV())
	if err != nil{
		return
	}
	//保存成功
	//如果是更新,那么返回旧值
	if putResp.PrevKv != nil{
		//对旧值反序列化
		err = json.Unmarshal(putResp.PrevKv.Value,&oldJobObj)
		if err != nil{
			err = nil //如果反序列化出错.跳过
			return
		}
		oldJob = &oldJobObj
	}

	return
}

//删除任务
func (jobMgr *JobMgr)DeleteJob(name string )(oldJob *common.Job,err error)  {
	var (
		jobKey string
	)
	//etcd保存任务的key
	jobKey = common.JOB_SAVE_DIR + name
	//从etcd删除它
	delResp,err := jobMgr.kv.Delete(context.TODO(),jobKey,clientv3.WithPrevKV())
	if err != nil{
		return nil,err
	}
	//返回被删除的任务信息
	if len(delResp.PrevKvs) != 0{
		oldJobObj := common.Job{}
		//解析并返回
		err := json.Unmarshal(delResp.PrevKvs[0].Value,&oldJobObj)
		if err != nil{
			err = nil
			return nil,err
		}
		oldJob = &oldJobObj
	}
	return
}

//列出任务
func (jobMgr *JobMgr)ListJobs()(jobList []*common.Job,err error)  {
	var (
		dirkey string
	)
	//初始化数组空间
	jobList = make([]*common.Job,0) //不会返回nil
	dirkey = common.JOB_SAVE_DIR
	getResp,err := jobMgr.kv.Get(context.TODO(),dirkey,clientv3.WithPrefix())
	if err != nil{
		return jobList, err
	}

	//成功,foreach,反序列化
	for _,kvPair := range getResp.Kvs{
		job := &common.Job{}
		err := json.Unmarshal(kvPair.Value,job)
		if err!= nil{
			continue
		}
		jobList = append(jobList,job)
	}


	return jobList,nil
}

//杀死任务
func (jobMgr *JobMgr)KillJob(name string)(err error)  {
	//更新一下key=/cron/killer 任务名
	var(
		killerKey string
		leaseGrantResp *clientv3.LeaseGrantResponse
		leaseId clientv3.LeaseID
	)
	//通知worker杀死对应任务
	killerKey = common.JOB_KILLER_DIR + name
	fmt.Println(killerKey)

	//让worker监听到一次put操作,创建一个租约让其稍后自动过期即可
	leaseGrantResp,err = jobMgr.lease.Grant(context.TODO(),30)
	if err != nil{
		return
	}

	//租约ID
	leaseId = leaseGrantResp.ID

	//设置killer标记
	_,err = jobMgr.kv.Put(context.TODO(),killerKey,"",clientv3.WithLease(leaseId))
	if err != nil{
		return
	}






	return
}