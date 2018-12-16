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

	//lease实现锁自动过期
	//op 操作
	//txn事务 if else then

	//1,上锁(创建租约,自动续租,拿着租约去抢占一个key

	//申请一个lease(租约)
	lease := clientv3.NewLease(client)
	//5s的租约
	leaseGrantResp, err := lease.Grant(context.TODO(), 5)
	if err != nil {
		fmt.Println(err)
		return
	}
	leaseId := leaseGrantResp.ID

	//准备一个用于取消自动续租的context
	ctx, cancelFunc := context.WithCancel(context.TODO())

	defer cancelFunc()                          //确保函数退出后,自动续租停止
	defer lease.Revoke(context.TODO(), leaseId) //撤回租约

	//自动续租约
	keepAliveChan, err := lease.KeepAlive(ctx, leaseId)
	if err != nil {
		fmt.Println(err)
		return
	}
	//启动一个协程,接收keepAlive的信息
	go func() {
		for {
			select {
			case keepResp := <-keepAliveChan:
				if keepAliveChan == nil {
					fmt.Println("租约已经失效了")
					goto END
				} else if keepResp != nil {
					//每秒会续租一次,所以就会受到一次应答
					fmt.Println("收到自动续租应答", keepResp.ID)
				}

			}
		}
	END:
		fmt.Println("结束接收应答的协程")
	}()

	//if 不存在key, then 设置它,else 抢锁失败
	kv := clientv3.NewKV(client)
	//创建事务
	txn := kv.Txn(context.TODO())
	//定义事务

	//如果key不存在
	txn.If(clientv3.Compare(clientv3.CreateRevision("/cron/lock/job9"), "=", 0)).
		Then(clientv3.OpPut("/cron/lock/job9", "xxx", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet("/cron/lock/job9")) //否则抢锁失败

	txnResp, err := txn.Commit()
	if err != nil {
		fmt.Println(err)
		return //提交失败
	}
	//判断是否报到了锁
	if !txnResp.Succeeded {
		fmt.Println("锁被占用:", string(txnResp.Responses[0].GetResponseRange().Kvs[0].Value))
		return
	}

	//2,处理业务
	//在锁内,很安全

	fmt.Println("处理任务")
	time.Sleep(5 * time.Second)
	//3,释放锁(取消自动续租).释放租约
	//defer会把租约自动取消

}
