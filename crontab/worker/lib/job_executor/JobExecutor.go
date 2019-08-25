package job_executor

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