package job_mgr
//任务锁 接口,
type JobLocker interface {
	TryLock() (err error)//抢锁
	Unlock()//释放锁
}
