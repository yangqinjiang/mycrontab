package lib

import (
	logs "github.com/sirupsen/logrus"
	"github.com/yangqinjiang/mycrontab/worker/common"
)

type CustomJobEventReceiver struct {
	JobEventReceiver JobEventReceiver
}

//推送保存任务的事件到Scheduler
func (this *CustomJobEventReceiver) PushSaveEventToScheduler(jobKey string, value []byte) {
	logs.Info("推送保存任务的事件到Scheduler,jobKey=",jobKey)
	var job *common.Job
	var err error
	if job, err = common.UnpackJob(value); err != nil {
		return
	}
	//反序列化成功
	jobName := common.ExtractJobName(jobKey)
	//构建一个更新event事件
	logs.Warn("推送 [ 保存 ] 任务的事件到Scheduler ,JobName = ", jobName)
	jobEvent := common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
	//推送给scheduler
	this.PushToScheduler(jobEvent)
}

//推送删除任务的事件到Scheduler
func (this *CustomJobEventReceiver) PushDeleteEventToScheduler(jobKey string) {
	// Delete /cron/jobs/job10
	jobName := common.ExtractJobName(jobKey)
	job := &common.Job{
		Name: jobName,
	}
	logs.Warn("推送 [ 删除 ] 任务的事件到Scheduler ,JobName = ", jobName)
	//构建一个删除event
	jobEvent := common.BuildJobEvent(common.JOB_EVENT_DELETE, job)
	//推送给scheduler
	this.PushToScheduler(jobEvent)
}

//推送强杀任务的事件到Scheduler
func (this *CustomJobEventReceiver) PushKillEventToScheduler(jobKey string) {
	jobName := common.ExtractKillerName(jobKey)
	logs.Warn("推送 [ 强杀 ] 任务的事件到Scheduler ,JobName = ", jobName)
	job := &common.Job{Name: jobName}
	jobEvent := common.BuildJobEvent(common.JOB_EVENT_KILL, job)
	//推送给scheduler
	this.PushToScheduler(jobEvent)
}
//推送给scheduler
func (this *CustomJobEventReceiver) PushToScheduler(jobEvent *common.JobEvent) {
	if nil != this.JobEventReceiver {
		this.JobEventReceiver.PushEvent(jobEvent)
	} else {
		logs.Error("没设置JobEventReceiver对象")
	}
}
