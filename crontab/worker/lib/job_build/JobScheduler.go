package job_build

import (
	"github.com/yangqinjiang/mycrontab/worker/common"
)
//推送任务事件的管理者
type JobEventReceiver interface {
	PushEvent(jobEvent *common.JobEvent)
}
//推送任务执行结果  事件的管理者
type JobResultReceiver interface {
	PushResult(jobResult *common.JobExecuteResult)
}
