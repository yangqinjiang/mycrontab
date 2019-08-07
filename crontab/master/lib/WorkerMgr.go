package lib

import (
	"context"
	"github.com/yangqinjiang/mycrontab/master/common"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"sync"
	"time"
)

type WorkerMgr struct {
	client *clientv3.Client
	kv     clientv3.KV
	lease  clientv3.Lease
}

var (
	G_workerMgr   *WorkerMgr
	onceWorkerMgr sync.Once
)

func InitWorkerMgr() (err error) {
	onceWorkerMgr.Do(func() {

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
		G_workerMgr = &WorkerMgr{
			client: client,
			kv:     kv,
			lease:  lease,
		}
	})
	return
}

//获取在线worker列表
func (workerMgr *WorkerMgr) ListWorker() (workerArr []string, err error) {
	var (
		getResp *clientv3.GetResponse
		kv      *mvccpb.KeyValue
	)
	//初始化数组
	workerArr = make([]string, 0)
	//获取目录下的所有kv
	if getResp, err = workerMgr.kv.Get(context.TODO(), common.JOB_WORKER_DIR, clientv3.WithPrefix()); err != nil {
		return
	}
	//解析每个节点的IP
	for _, kv = range getResp.Kvs {
		//kv.key = /cron/workers/192.168.1.x
		worderIP := common.ExtractWorkerIP(string(kv.Key))
		workerArr = append(workerArr, worderIP)
	}
	return
}
