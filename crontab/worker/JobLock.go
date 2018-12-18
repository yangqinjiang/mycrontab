package worker

import (
	"context"
	"go.etcd.io/etcd/clientv3"
	"github.com/yangqinjiang/mycrontab/crontab/common"
)

//分布式锁(TXN事务)
type JobLock struct {
	Kv         clientv3.KV
	Lease      clientv3.Lease
	JobName    string             //任务名
	cancelFunc context.CancelFunc //用于终止自动续租
	leaseId clientv3.LeaseID//租约ID
	isLocked bool//是否抢到锁
}

//初始化一把锁
func InitJobLock(jobName string, kv clientv3.KV, lease clientv3.Lease) (jobLock *JobLock) {
	jobLock = &JobLock{
		Kv:      kv,
		Lease:   lease,
		JobName: jobName,
	}
	return
}

//抢锁
func (jobLock *JobLock) TryLock() (err error) {
	var (
		leaseGrantResp *clientv3.LeaseGrantResponse
		cancelCtx      context.Context
		cancelFunc     context.CancelFunc
		keepRespChan   <-chan *clientv3.LeaseKeepAliveResponse
		txn            clientv3.Txn
		lockKey        string
		txnResp        *clientv3.TxnResponse
	)
	//1,创建租约 5s
	if leaseGrantResp, err = jobLock.Lease.Grant(context.TODO(), 5); err != nil {
		return
	}
	//2自动续租
	//contextCtl 用于取消自动续租
	cancelCtx, cancelFunc = context.WithCancel(context.TODO())
	//续租ID
	leaseId := leaseGrantResp.ID
	if keepRespChan, err = jobLock.Lease.KeepAlive(cancelCtx, leaseId); err != nil {
		goto FAIL
	}
	//3处理续租应答的协程
	go func() {
		var (
			keepResp *clientv3.LeaseKeepAliveResponse
		)
		for {
			select {
			case keepResp = <-keepRespChan: //自动续租的应答
				if keepResp == nil {
					goto END
				}
			}
		}
	END:
	}()

	//3创建事务txn
	txn = jobLock.Kv.Txn(context.TODO())
	//锁路径
	lockKey = common.JOB_LOCK_DIR + jobLock.JobName
	//4,事务抢锁
	txn.If(clientv3.Compare(clientv3.CreateRevision(lockKey), "=", 0)).
		Then(clientv3.OpPut(lockKey, "", clientv3.WithLease(leaseId))).
		Else(clientv3.OpGet(lockKey))
	//提交事务
	if txnResp, err = txn.Commit(); err != nil {
		goto FAIL
	}

	//5,成功返回,失败释放租约
	if !txnResp.Succeeded { //锁被占用
		err = common.ERR_LOCK_ALREADY_REQUIRED
		//
		goto FAIL
	}
	//抢锁成功
	jobLock.leaseId = leaseId
	jobLock.cancelFunc = cancelFunc
	jobLock.isLocked = true
	return
FAIL:
	cancelFunc()
	jobLock.Lease.Revoke(context.TODO(), leaseId)
	return
}

func (jobLock *JobLock) Unlock() {
	if jobLock.isLocked {
	jobLock.cancelFunc()//取消我们程序自动续租的协程
	jobLock.Lease.Revoke(context.TODO(),jobLock.leaseId)//释放租约
	}
}
