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
	Exec(info *common.JobExecuteInfo,f func(info *common.JobExecuteInfo) ([]byte, error))(error)
}
