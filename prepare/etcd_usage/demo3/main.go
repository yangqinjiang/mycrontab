package main

import (
	"fmt"
	"go.etcd.io/etcd/clientv3"
	"time"
	"context"
	"go.etcd.io/etcd/mvcc/mvccpb"
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

	//模拟etcd的kv变化
	go func() {
		for{
			kv.Put(context.TODO(),"/cron/jobs/job7","I am job7")
			kv.Delete(context.TODO(),"/cron/jobs/job7")
			time.Sleep(1*time.Second)
		}
	}()

	//先GET到当前的值,并监听后续变化
	getResp,err := kv.Get(context.TODO(),"/cron/jobs/job7")
	if err != nil{
		fmt.Println(err)
		return
	}
	//key是存在的
	if len(getResp.Kvs) != 0{
		fmt.Println("当前值:",string(getResp.Kvs[0].Value))

	}

	//当前etcd集群事务ID,单调递增的
	watchStartRevision := getResp.Header.Revision + 1
	//创建一个watcher
	watcher := clientv3.NewWatcher(client)
	//启动监听
	fmt.Println("从该版本向后监听",watchStartRevision)

	//5s后自动取消
	ctx,cancelFunc := context.WithCancel(context.TODO())
	time.AfterFunc(5*time.Second, func() {
		cancelFunc()
	})
	//
	watchRespChan := watcher.Watch(ctx,"/cron/jobs/job7",clientv3.WithRev(watchStartRevision))

	//处理kv变化事件
	for watchResp := range watchRespChan{
		for _,event :=range watchResp.Events{
			switch event.Type {
				//put
			case mvccpb.PUT:
				fmt.Println("修改为:",string(event.Kv.Value),"Revision:",event.Kv.CreateRevision,event.Kv.ModRevision)
			case mvccpb.DELETE:
				fmt.Println("删除了","Revision:",event.Kv.ModRevision)
			}
		}
	}



}
