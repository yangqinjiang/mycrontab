package worker

import (
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"time"
)

//日志接口类
type JobLoger interface {
	Write(jobLog *common.LogBatch) (n int, err error)
}
//任务的执行器 的接口
type JobExecuter interface {
	//设置调用者
	SetCommand(c Command)
	//执行
	Exec(info *common.JobExecuteInfo)(error)
}
//任务计划 接口
type JobPlanManager interface {
	Size() int
	//插入一个任务
	Insert(info *common.JobSchedulePlan)(error)
	// 使用key 删除一个任务
	Remove(key string)(error)
	//找出最早
	ExtractEarliest(func (jobPlan *common.JobSchedulePlan)(err error)) (time.Duration,error)
}
//推送任务事件的管理者
type JobEventReceiver interface {
	PushEvent(jobEvent *common.JobEvent)
}
//推送任务执行结果  事件的管理者
type JobResultReceiver interface {
	PushResult(jobResult *common.JobExecuteResult)
}
