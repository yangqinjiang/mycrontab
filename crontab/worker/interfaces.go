package worker

import (
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"io"
)

//日志接口类
type Log interface {
	io.Writer
}

//任务日志的接口
type JobLogger interface {
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
}
