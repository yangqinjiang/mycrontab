package main

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
)

func main() {
	//在服务器启动etcd
	// nohup ./etcd --listen-client-urls 'http://0.0.0.0:2379' --advertise-client-urls 'http://0.0.0.0:2379' &
	config := clientv3.Config{
		Endpoints:   []string{"127.0.0.1:2379"}, //集群列表
		DialTimeout: 5 * time.Second,
	}

	//建立一个客户端
	client, err := clientv3.New(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("链接成功")

	//用于读写etcd的键值对
	kv := clientv3.NewKV(client)
	kv = kv

	putOp := clientv3.OpPut("/cron/jobs/job8", "123123123")

	opResp, err := kv.Do(context.TODO(), putOp)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("写入Revision:", opResp.Put().Header.Revision)

	//创建Op
	getOp := clientv3.OpGet("/cron/jobs/job8")

	//执行op
	getOpRsp, err := kv.Do(context.TODO(), getOp)
	if err != nil {
		fmt.Println(err)
		return
	}
	//打印
	fmt.Println("数据Revision:", getOpRsp.Get().Kvs[0].ModRevision)
	fmt.Println("数据value:", string(getOpRsp.Get().Kvs[0].Value))

}
