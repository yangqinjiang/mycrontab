package worker

import (
	"go.etcd.io/etcd/clientv3"
	"time"
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
