package worker

import (
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"time"
)

//日志接口类
type JobLoger interface {
	Write(jobLog *common.JobLog) (n int, err error)
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
	ExtractEarliest(func (jobPlan *common.JobSchedulePlan)(err error)) (time.Duration)
}
