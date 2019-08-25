package scheduler

import (
	"errors"
	logs "github.com/sirupsen/logrus"
	"github.com/yangqinjiang/mycrontab/worker/common"
	"github.com/yangqinjiang/mycrontab/worker/lib/command"
	"testing"
	"time"
)

//初始化任务调度器
func TestInitScheduler(t *testing.T) {
	//第一次初始化任务调度器
	err,s1 := scheduler.InitScheduler(nil)
	if err != nil {
		t.Fatal("初始化任务调度器 失败",err)
	}
	//第二次初始化
	err ,s2 := scheduler.InitScheduler(nil)
	if err != nil {
		t.Fatal("初始化任务调度器 失败",err)
	}
	if s1 != s2{
		t.Fatal("任务调度器  单例模式出错了")
	}
}

func TestScheduler_PushJobEvent(t *testing.T) {
	//第一次初始化任务调度器
	err,_ := scheduler.InitScheduler(nil)
	if err != nil {
		t.Fatal("初始化任务调度器 失败",err)
	}
	job := &common.Job{Name: "PushJobEvent"}
	jobEvent := common.BuildJobEvent(common.JOB_EVENT_KILL, job)
	scheduler.G_scheduler.PushEvent(jobEvent)
	scheduler.G_scheduler.PushEvent(jobEvent)
	if (scheduler.G_scheduler.JobEventChanLen() != 2){
		t.Fatal("PushJobEvent 失败,数量==2")
	}
	//不设置G_scheduler
	pusher := &scheduler.CustomJobEventReceiver{scheduler.JobEventReceiver: nil}
	b ,err := common.PackJob(job)
	if err != nil{
		t.Fatal("序列化job 出错")
	}
	pusher.PushSaveEventToScheduler("PushSaveEventToScheduler",b)
	pusher.PushKillEventToScheduler("PushKillEventToScheduler" )
	pusher.PushDeleteEventToScheduler("PushDeleteEventToScheduler" )
	if (scheduler.G_scheduler.JobEventChanLen() != 2){
		t.Fatal("这里不设置G_scheduler,不能出错")
	}

	//设置G_scheduler
	pusher = &scheduler.CustomJobEventReceiver{scheduler.JobEventReceiver: scheduler.G_scheduler}
	b ,err = common.PackJob(job)
	if err != nil{
		t.Fatal("序列化job 出错")
	}
	pusher.PushSaveEventToScheduler("PushSaveEventToScheduler",b)
	pusher.PushKillEventToScheduler("PushKillEventToScheduler" )
	pusher.PushDeleteEventToScheduler("PushDeleteEventToScheduler" )
	if (scheduler.G_scheduler.JobEventChanLen() != 5){
		t.Fatal("PushJobEvent 失败,数量==5")
	}
	b = append(b, 1)
	pusher.PushSaveEventToScheduler("PushSaveEventToScheduler",b)
	if (scheduler.G_scheduler.JobEventChanLen() != 5){
		t.Fatal("PushJobEvent 失败,数量==5")
	}
}

func TestScheduler_PushJobResult(t *testing.T) {
	//第一次初始化任务调度器
	err,_ := scheduler.InitScheduler(nil)
	if err != nil {
		t.Fatal("初始化任务调度器 失败",err)
	}
	//任务执行的结果
	result := &common.JobExecuteResult{
		ExecuteInfo: nil,
		Output:      make([]byte, 0),
		StartTime:   time.Now(),
	}
	scheduler.G_scheduler.PushResult(result)
	scheduler.G_scheduler.PushResult(result)
	scheduler.G_scheduler.PushResult(result)
	if (scheduler.G_scheduler.JobResultChanLen() != 3){
		t.Fatal("PushJobResult 失败,数量== 3")
	}
}


//尝试执行任务
func TestScheduler_TryStartJobSuccess(t *testing.T) {
	//第一次初始化任务调度器
	err,_ := scheduler.InitScheduler(nil)
	if err != nil {
		t.Fatal("初始化任务调度器 失败",err)
	}
	 

	//调用者 执行 命令对象的execute函数
	invoker := &TestJobExecSuccessInvoker{}
	//加载命令给调用者
	invoker.SetCommand(scheduler.CommandFactory("sh"))


	//invoker.SetCommand()
	scheduler.G_scheduler.SetJobExecuter(invoker)

	//FAIL,cron表达式是错误的,
	job_fail := &common.Job{Name: "TryStartJob",CronExpr:"error cron", command.Command:"echo hello"}
	jobEvent_fail := common.BuildJobEvent(common.JOB_EVENT_KILL, job_fail)

	if _, err = common.BuildJobSchedulePlan(jobEvent_fail.Job); err == nil {
		t.Fatal("构建任务调度计划应该是失败,因为cronExpr是错误的",err)
		return
	}
	//cron表达式是正确的
	job := &common.Job{Name: "TryStartJobSuccess",CronExpr:"* * * * * *", command.Command:"echo hello"}
	jobEvent := common.BuildJobEvent(common.JOB_EVENT_KILL, job)
	var jobSchedulePlan *common.JobSchedulePlan
	if jobSchedulePlan, err = common.BuildJobSchedulePlan(jobEvent.Job); err != nil {
		t.Fatal("构建任务调度计划失败",err)
		return
	}
	err = scheduler.G_scheduler.TryStartJob(nil)
	if err == nil{
		t.Fatal("参数jobPlan为空",err)
	}
	err = scheduler.G_scheduler.TryStartJob(jobSchedulePlan)
	if err != nil{
		t.Fatal("执行任务应该是正确的",err)
	}
	//处理任务执行的结果
	_,err = scheduler.G_scheduler.HandleJobResult(nil)
	if (nil == err){
		t.Fatal("处理任务执行的结果应该是失败的",err)
	}
	jobExecuteInfo := common.BuildJobExecuteInfo(jobSchedulePlan)
	//任务执行的结果
	result := &common.JobExecuteResult{
		ExecuteInfo: jobExecuteInfo,
		Output:      make([]byte, 0),
		StartTime:   time.Now(),
	}
	//处理任务执行的结果
	_,err = scheduler.G_scheduler.HandleJobResult(result)
	if (nil == err){
		t.Fatal("处理任务执行的结果应该是失败的",err)
	}

	//设置日志记录器
	w := &scheduler.TestWriter{}
	scheduler.G_scheduler.SetJobLogBuffer(w)
	_,err = scheduler.G_scheduler.HandleJobResult(result)
	if (nil != err){
		t.Fatal("处理任务执行的结果应该是成功的",err)
	}

}


//尝试执行任务
func TestScheduler_TryStartJobFail(t *testing.T) {
	//第一次初始化任务调度器
	err,_ := scheduler.InitScheduler(nil)
	if err != nil {
		t.Fatal("初始化任务调度器 失败",err)
	}
	logs.Info("设置空的任务的执行器 nil")
	scheduler.G_scheduler.SetJobExecuter(nil)//设置为nil值
	// no JobExec
	err = scheduler.G_scheduler.TryStartJob(nil)
	if err == nil{
		t.Fatal("执行器应该是nil值",err)
	}
	logs.Info("设置任务的执行器")
	scheduler.G_scheduler.SetJobExecuter(&TestJobExecFailInvoker{})
	err = scheduler.G_scheduler.TryStartJob(nil)
	if err == nil{
		t.Fatal("参数jobPlan为空",err)
	}
	//cron表达式是正确的
	job := &common.Job{Name: "TryStartJobFail",CronExpr:"* * * * * *", command.Command:"echo hello"}
	jobEvent := common.BuildJobEvent(common.JOB_EVENT_KILL, job)
	var jobSchedulePlan *common.JobSchedulePlan
	if jobSchedulePlan, err = common.BuildJobSchedulePlan(jobEvent.Job); err != nil {
		t.Fatal("构建任务调度计划失败",err)
		return
	}
	logs.Info("执行任务")
	err = scheduler.G_scheduler.TryStartJob(jobSchedulePlan)
	if err.Error() != "执行任务失败"{
		t.Fatal("执行任务应该是出错的",err)
	}
	err = scheduler.G_scheduler.TryStartJob(jobSchedulePlan)
	if err.Error() != "尚未退出,跳过执行"{
		t.Fatal("执行任务应该是出错的",err)
	}
}


//---------------------
//执行任务成功的
type TestJobExecSuccessInvoker struct {
	scheduler.JobExecuter
	c scheduler.Command
}

func (t *TestJobExecSuccessInvoker)Exec(info *common.JobExecuteInfo)(err error)  {
	logs.Info("TestJobExecSuccessInvoker 正在执行一个成功的任务:",info.Job.Name,info.Job.Command,info.Job.CronExpr,)
	time.Sleep(time.Second)
	_, err =  t.c.Execute(info)
	return
}
//设置命令对象
func (t *TestJobExecSuccessInvoker)SetCommand(c scheduler.Command)  {
	logs.Info("call TestJobExecSuccessInvoker SetCommand")
	t.c = c
}
//执行任务失败的
type TestJobExecFailInvoker struct {
	scheduler.JobExecuter
	c scheduler.Command
}

func (je *TestJobExecFailInvoker)Exec(info *common.JobExecuteInfo)(err error)  {
	logs.Info("正在执行失败的任务:",info.Job.Name,info.Job.Command,info.Job.CronExpr,)
	time.Sleep(2*time.Second)
	return  errors.New("执行任务失败")
}
//设置命令对象
func (t *TestJobExecFailInvoker)SetCommand(c scheduler.Command)  {
	t.c = c
}