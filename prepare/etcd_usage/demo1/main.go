package main

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
	"context"
)

func main() {
	//在服务器启动etcd
	// nohup ./etcd --listen-client-urls 'http://0.0.0.0:2379' --advertise-client-urls 'http://0.0.0.0:2379' &
	config:=clientv3.Config{
		Endpoints:[]string{"127.0.0.1:2379"},//集群列表
		DialTimeout:5*time.Second,
	}

	//建立一个客户端
	client,err := clientv3.New(config)
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println("链接成功")
	//用于读写etcd的键值对
	kv:=clientv3.NewKV(client)
	putResp,err := kv.Put(context.TODO(),"/cron/jobs/job1","hi0",clientv3.WithPrevKV())
	if err != nil{
		fmt.Println(err)
		return
	}else{
		fmt.Println("Revision:",putResp.Header.Revision)

		if putResp.PrevKv !=nil{//put的时候,必须 clientv3.WithPrevKV()
			fmt.Println("PrevKv",string(putResp.PrevKv.Value))
		}

		//读取
		fmt.Println("读取...")
		getResp ,err := kv.Get(context.TODO(),"/cron/jobs/job1")
		if err !=nil{
			fmt.Println(err)
			return
		}else{
			fmt.Println(getResp.Kvs)
		}
	}




}
