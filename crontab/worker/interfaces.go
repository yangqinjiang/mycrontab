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
