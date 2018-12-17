package worker

import "go.etcd.io/etcd/clientv3"

//分布式锁(TXN事务)
type JobLock struct {
	Kv clientv3.KV
	Lease clientv3.Lease
	JobName string //任务名
}

//初始化一把锁
func InitJobLock(jobName string,kv clientv3.KV,lease clientv3.Lease) (jobLock *JobLock)  {
	jobLock = &JobLock{
		Kv:kv,
		Lease:lease,
		JobName:jobName,
	}
	return
}