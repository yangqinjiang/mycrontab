package worker

import (
	"testing"
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"time"
)

//初始化任务调度器
func TestInitScheduler(t *testing.T) {
	//第一次初始化任务调度器
	err,s1 := InitScheduler(nil)
	if err != nil {
		t.Fatal("初始化任务调度器 失败",err)
	}
	//第二次初始化
	err ,s2 := InitScheduler(nil)
	if err != nil {
		t.Fatal("初始化任务调度器 失败",err)
	}
	if s1 != s2{
		t.Fatal("任务调度器  单例模式出错了")
	}
}

func TestScheduler_PushJobEvent(t *testing.T) {
	//第一次初始化任务调度器
	err,_ := InitScheduler(nil)
	if err != nil {
		t.Fatal("初始化任务调度器 失败",err)
	}
	job := &common.Job{Name: "PushJobEvent"}
	jobEvent := common.BuildJobEvent(common.JOB_EVENT_KILL, job)
	G_scheduler.PushJobEvent(jobEvent)
	G_scheduler.PushJobEvent(jobEvent)
	if (G_scheduler.JobEventChanLen() != 2){
		t.Fatal("PushJobEvent 失败,数量==2")
	}
}

func TestScheduler_PushJobResult(t *testing.T) {
	//第一次初始化任务调度器
	err,_ := InitScheduler(nil)
	if err != nil {
		t.Fatal("初始化任务调度器 失败",err)
	}
	//任务执行的结果
	result := &common.JobExecuteResult{
		ExecuteInfo: nil,
		Output:      make([]byte, 0),
		StartTime:   time.Now(),
	}
	G_scheduler.PushJobResult(result)
	G_scheduler.PushJobResult(result)
	G_scheduler.PushJobResult(result)
	if (G_scheduler.JobResultChanLen() != 3){
		t.Fatal("PushJobResult 失败,数量== 3")
	}
}

type TestJobExec struct {
}

func (je *TestJobExec)Exec(info *common.JobExecuteInfo)(err error)  {
	return
}

//尝试执行任务
func TestScheduler_TryStartJob(t *testing.T) {
	//第一次初始化任务调度器
	err,_ := InitScheduler(nil)
	if err != nil {
		t.Fatal("初始化任务调度器 失败",err)
	}
	G_scheduler.SetJobExecuter(&TestJobExec{})

	//FAIL,cron表达式是错误的,
	job_fail := &common.Job{Name: "TryStartJob",CronExpr:"error cron",Command:"echo hello"}
	jobEvent_fail := common.BuildJobEvent(common.JOB_EVENT_KILL, job_fail)

	if _, err = common.BuildJobSchedulePlan(jobEvent_fail.Job); err == nil {
		t.Fatal("构建任务调度计划应该是失败,因为cronExpr是错误的",err)
		return
	}
	//cron表达式是正确的
	job := &common.Job{Name: "TryStartJob",CronExpr:"* * * * * *",Command:"echo hello"}
	jobEvent := common.BuildJobEvent(common.JOB_EVENT_KILL, job)
	var jobSchedulePlan *common.JobSchedulePlan
	if jobSchedulePlan, err = common.BuildJobSchedulePlan(jobEvent.Job); err != nil {
		t.Fatal("构建任务调度计划失败",err)
		return
	}
	err = G_scheduler.TryStartJob(jobSchedulePlan)
	if err != nil{
		t.Fatal("执行任务出错")
	}
}
