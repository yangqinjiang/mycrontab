package master

import (
	"go.etcd.io/etcd/clientv3"
	"time"
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"encoding/json"
	"context"
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
func (jobMgr *JobMgr)SaveJob(job *common.Job) (oldJobObj common.Job,err error)  {
	//把任务保存到/cron/jobs/任务名 -> json
	var (
		jobKey string
		jobValue []byte
	)

	//etcd保存的key
	jobKey = "/cron/jobs/"+job.Name
	//任务信息json
	jobValue,err = json.Marshal(job)
	if err != nil{
		return
	}
	//保存到etcd
	//返回旧值
	putResp,err := jobMgr.kv.Put(context.TODO(),jobKey,string(jobValue),clientv3.WithPrevKV())
	if err != nil{
		return
	}
	//保存成功
	//如果是更新,那么返回旧值
	if putResp.PrevKv != nil{
		//对旧值反序列化
		err := json.Unmarshal(putResp.PrevKv.Value,&oldJobObj)
		if err != nil{
			err = nil //如果反序列化出错.跳过
			return
		}
	}

	return
}