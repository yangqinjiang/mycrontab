package lib

import (
	"github.com/yangqinjiang/mycrontab/worker/lib/command"
	"github.com/yangqinjiang/mycrontab/worker/common"
)


//任务的执行器 的接口
type JobExecuter interface {
	//设置调用者
	SetCommand(c command.Command)
	//执行
	Exec(info *common.JobExecuteInfo)(error)
}

//推送任务事件的管理者
type JobEventReceiver interface {
	PushEvent(jobEvent *common.JobEvent)
}
//推送任务执行结果  事件的管理者
type JobResultReceiver interface {
	PushResult(jobResult *common.JobExecuteResult)
}
//任务锁 接口,
type JobLocker interface {
	TryLock() (err error)//抢锁
	Unlock()//释放锁
}
