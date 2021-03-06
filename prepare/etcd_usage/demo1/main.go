package main

import (
	"context"
	"fmt"
	"github.com/coreos/etcd/clientv3"
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
	putResp, err := kv.Put(context.TODO(), "/cron/jobs/job1", "hi1", clientv3.WithPrevKV())
	if err != nil {
		fmt.Println(err)
		return
	} else {
		fmt.Println("Revision:", putResp.Header.Revision)

		if putResp.PrevKv != nil { //put的时候,必须 clientv3.WithPrevKV()
			fmt.Println("PrevKv", string(putResp.PrevKv.Value))
		}

		//读取
		fmt.Println("读取...")
		getResp, err := kv.Get(context.TODO(), "/cron/jobs/job1")
		if err != nil {
			fmt.Println(err)
			return
		} else {
			fmt.Println(getResp.Kvs)
		}
	}

	//4
	kv.Put(context.TODO(), "/cron/jobs/job2", "hi2", clientv3.WithPrevKV())
	kv.Put(context.TODO(), "/cron/jobs/job3", "hi3", clientv3.WithPrevKV())
	kv.Put(context.TODO(), "/cron/jobs/job4", "hi4", clientv3.WithPrevKV())
	//读取/cron/jobs/为前缀的所有key
	fmt.Println("withPrefix")
	getResp, err := kv.Get(context.TODO(), "/cron/jobs", clientv3.WithPrefix())
	if err != nil {
		fmt.Println(err)
		return
	} else {
		//获取成功,遍历所有的Kvs
		fmt.Println(getResp.Kvs)
	}

	//5,删除kv
	//还有更多删除的选项
	delResp, err := kv.Delete(context.TODO(), "/cron/jobs/job1", clientv3.WithPrevKV())
	//
	if err != nil {
		fmt.Println(err)
		return
	} else {
		//被删除之前的PreKv是什么
		if len(delResp.PrevKvs) != 0 {
			for _, kvpair := range delResp.PrevKvs {
				fmt.Println("删除了:", string(kvpair.Key), string(kvpair.Value))
			}
		}
	}

}
