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

	//申请一个lease(租约)
	lease := clientv3.NewLease(client)
	//10s的租约
	leaseGrantResp,err := lease.Grant(context.TODO(),10)
	if err != nil{
		fmt.Println(err)
		return
	}
	leaseId := leaseGrantResp.ID

	//超时
	//ctx,_ := context.WithTimeout(context.TODO(),5*time.Second)

	//自动续租约
	keepAliveChan,err := lease.KeepAlive(context.TODO(),leaseId)
	if err != nil{
		fmt.Println(err)
		return
	}
	//启动一个协程,接收keepAlive的信息
	go func() {
		for{
			select {
				case keepResp := <- keepAliveChan:
					if keepAliveChan == nil{
						fmt.Println("租约已经失效了")
						goto END
					}else{
						//每秒会续租一次,所以就会受到一次应答
						fmt.Println("收到自动续租应答",keepResp.ID)
					}

			}
		}
		END:
			fmt.Println("结束接收应答的协程")
	}()

	//用于读写etcd的键值对
	kv:=clientv3.NewKV(client)
	//put一个kv,让它与租约关联起来,从而实现10s后自动过期
	//opt with-lease
	putResp,err := kv.Put(context.TODO(),"/cron/lock/job1","",clientv3.WithLease(leaseId))
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println("写入成功:",putResp.Header.Revision)

	//定时看一下key过期与否
	for{
		getResp,err := kv.Get(context.TODO(),"/cron/lock/job1")
		if err != nil{
			return
		}
		if getResp.Count == 0{
			fmt.Println("kv过期了")
			break
		}
		fmt.Println("kv还没过期",getResp.Kvs)
		time.Sleep(2*time.Second)
	}


}
