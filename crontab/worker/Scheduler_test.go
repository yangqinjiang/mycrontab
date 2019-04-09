package worker

import (
	"github.com/astaxie/beego/logs"
	"testing"
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"time"
	"errors"
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


//尝试执行任务
func TestScheduler_TryStartJobSuccess(t *testing.T) {
	//第一次初始化任务调度器
	err,_ := InitScheduler(nil)
	if err != nil {
		t.Fatal("初始化任务调度器 失败",err)
	}
	/**
		// 命令接收者
		receA := NewReceiverA()

		//命令对象
		concomA := NewConcreteCommandA(*receA)

		invoker := NewInvoker()
		//加载命令给调用者
		invoker.SetCommand(concomA)


		//调用者 执行 命令对象的execute函数
		invoker.ExecuteCommand(nil)
	 */
	// 命令接收者
	//receA := NewCmdReceiver()

	//命令对象
	//concomA := NewConcreteCommand(*receA)

	//调用者 执行 命令对象的execute函数
	invoker := &TestJobExecSuccessInvoker{}
	//加载命令给调用者
	invoker.SetCommand(CommandFactory("sh"))


	//invoker.SetCommand()
	G_scheduler.SetJobExecuter(invoker)

	//FAIL,cron表达式是错误的,
	job_fail := &common.Job{Name: "TryStartJob",CronExpr:"error cron",Command:"echo hello"}
	jobEvent_fail := common.BuildJobEvent(common.JOB_EVENT_KILL, job_fail)

	if _, err = common.BuildJobSchedulePlan(jobEvent_fail.Job); err == nil {
		t.Fatal("构建任务调度计划应该是失败,因为cronExpr是错误的",err)
		return
	}
	//cron表达式是正确的
	job := &common.Job{Name: "TryStartJobSuccess",CronExpr:"* * * * * *",Command:"echo hello"}
	jobEvent := common.BuildJobEvent(common.JOB_EVENT_KILL, job)
	var jobSchedulePlan *common.JobSchedulePlan
	if jobSchedulePlan, err = common.BuildJobSchedulePlan(jobEvent.Job); err != nil {
		t.Fatal("构建任务调度计划失败",err)
		return
	}
	err = G_scheduler.TryStartJob(jobSchedulePlan)
	if err != nil{
		t.Fatal("执行任务应该是正确的",err)
	}
}


//尝试执行任务
func TestScheduler_TryStartJobFail(t *testing.T) {
	//第一次初始化任务调度器
	err,_ := InitScheduler(nil)
	if err != nil {
		t.Fatal("初始化任务调度器 失败",err)
	}
	logs.Info("设置空的任务的执行器 nil")
	G_scheduler.SetJobExecuter(nil)//设置为nil值
	// no JobExec
	err = G_scheduler.TryStartJob(nil)
	if err == nil{
		t.Fatal("执行器应该是nil值",err)
	}
	logs.Info("设置任务的执行器")
	G_scheduler.SetJobExecuter(&TestJobExecFailInvoker{})
	//cron表达式是正确的
	job := &common.Job{Name: "TryStartJobFail",CronExpr:"* * * * * *",Command:"echo hello"}
	jobEvent := common.BuildJobEvent(common.JOB_EVENT_KILL, job)
	var jobSchedulePlan *common.JobSchedulePlan
	if jobSchedulePlan, err = common.BuildJobSchedulePlan(jobEvent.Job); err != nil {
		t.Fatal("构建任务调度计划失败",err)
		return
	}
	logs.Info("执行任务")
	err = G_scheduler.TryStartJob(jobSchedulePlan)
	if err.Error() != "执行任务失败"{
		t.Fatal("执行任务应该是出错的",err)
	}
	err = G_scheduler.TryStartJob(jobSchedulePlan)
	if err.Error() != "尚未退出,跳过执行"{
		t.Fatal("执行任务应该是出错的",err)
	}
}


//---------------------
//执行任务成功的
type TestJobExecSuccessInvoker struct {
	c Command
}

func (t *TestJobExecSuccessInvoker)Exec(info *common.JobExecuteInfo)(err error)  {
	logs.Info("TestJobExecSuccessInvoker 正在执行一个成功的任务:",info.Job.Name,info.Job.Command,info.Job.CronExpr,)
	time.Sleep(2*time.Second)
	_, err =  t.c.Execute(info)
	return
}
//设置命令对象
func (t *TestJobExecSuccessInvoker)SetCommand(c Command)  {
	logs.Info("call TestJobExecSuccessInvoker SetCommand")
	t.c = c
}
//执行任务失败的
type TestJobExecFailInvoker struct {
	c Command
}

func (je *TestJobExecFailInvoker)Exec(info *common.JobExecuteInfo)(err error)  {
	logs.Info("正在执行失败的任务:",info.Job.Name,info.Job.Command,info.Job.CronExpr,)
	time.Sleep(2*time.Second)
	return  errors.New("执行任务失败")
}
//设置命令对象
func (t *TestJobExecFailInvoker)SetCommand(c Command)  {
	t.c = c
}