package job_build

import (
	"github.com/yangqinjiang/mycrontab/worker/common"
)

//推送任务执行结果  事件的管理者
type JobResultPusher interface {
	PushResult(jobResult *common.JobExecuteResult)
}
