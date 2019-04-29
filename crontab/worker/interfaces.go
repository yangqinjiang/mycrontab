package worker

import (
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"io"
	"time"
)

//日志接口类
type Log interface {
	io.Writer
}

//任务日志缓冲器的接口
type JobLogBuffer interface {
	Write(jobLog *common.JobLog)
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
