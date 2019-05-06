package worker

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/yangqinjiang/mycrontab/crontab/common"
)

type JobEventPusher struct {
	JobEventReceiver JobEventReceiver
}

//推送保存任务的事件到Scheduler
func (this *JobEventPusher) PushSaveEventToScheduler(jobKey string, value []byte) {
	var job *common.Job
	var err error
	if job, err = common.UnpackJob(value); err != nil {
		return
	}
	//反序列化成功
	jobName := common.ExtractJobName(jobKey)
	//构建一个更新event事件
	fmt.Println("推送保存任务的事件到Scheduler", jobName)
	jobEvent := common.BuildJobEvent(common.JOB_EVENT_SAVE, job)
	//推送给scheduler
	this.PushToScheduler(jobEvent)
}

//推送删除任务的事件到Scheduler
func (this *JobEventPusher) PushDeleteEventToScheduler(jobKey string) {
	// Delete /cron/jobs/job10
	jobName := common.ExtractJobName(jobKey)
	job := &common.Job{
		Name: jobName,
	}
	fmt.Println("推送删除任务的事件到Scheduler", jobName)
	//构建一个删除event
	jobEvent := common.BuildJobEvent(common.JOB_EVENT_DELETE, job)
	//推送给scheduler
	this.PushToScheduler(jobEvent)
}

//推送强杀任务的事件到Scheduler
func (this *JobEventPusher) PushKillEventToScheduler(jobKey string) {
	jobName := common.ExtractKillerName(jobKey)
	fmt.Println("推送强杀任务的事件到Scheduler", jobName)
	job := &common.Job{Name: jobName}
	jobEvent := common.BuildJobEvent(common.JOB_EVENT_KILL, job)
	//推送给scheduler
	this.PushToScheduler(jobEvent)
}
//推送给scheduler
func (this *JobEventPusher) PushToScheduler(jobEvent *common.JobEvent) {
	if nil != this.JobEventReceiver {
		this.JobEventReceiver.PushEvent(jobEvent)
	} else {
		logs.Error("没设置JobEventReceiver对象")
	}
}
