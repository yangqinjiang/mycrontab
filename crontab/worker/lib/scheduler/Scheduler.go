package scheduler

import (
	"errors"
	logs "github.com/sirupsen/logrus"
	"github.com/yangqinjiang/mycrontab/worker/lib/job_executor"
	"github.com/yangqinjiang/mycrontab/worker/common"
	//"github.com/yangqinjiang/mycrontab/worker/lib/job_build"
	"github.com/yangqinjiang/mycrontab/worker/lib/job_plan"
	"github.com/yangqinjiang/mycrontab/worker/lib/log"
	"sync"
	"time"
)

/**
调度器,遍历所有任务列表, 找出最近一个要过期的任务
*/
type Scheduler struct {
	//job_build.JobEventReceiver                                //推送任务事件的接口
	//job_build.SetJobResultPusher                               //推送任务执行结果 	的接口
	jobEventChan      chan *common.JobEvent             //etcd任务事件队列
	jobResultChan     chan *common.JobExecuteResult     //任务执行结果队列
	jobExecutingTable map[string]*common.JobExecuteInfo //任务执行表
	jobLogger         log.JobLoger                      //日志记录器
	jobExecuter       job_executor.JobExecuter          //任务执行器
	jobPlanManager    job_plan.JobPlanManager           //任务调度计划表内存里的任务计划管理
}
/**
日志记录器
*/
func (scheduler *Scheduler) SetJobLogBuffer(jobLogger log.JobLoger) {
	scheduler.jobLogger = jobLogger
}
/**
设置任务的执行器
*/
func (scheduler *Scheduler) SetJobExecuter(jobExecuter job_executor.JobExecuter) {
	scheduler.jobExecuter = jobExecuter
}

/**
设置任务计划的管理者
*/
func (scheduler *Scheduler) SetJobPlanManager(jobPlanManager job_plan.JobPlanManager) {
	scheduler.jobPlanManager = jobPlanManager
}

//处理任务事件
func (scheduler *Scheduler) handleJobEvent(jobEvent *common.JobEvent) {

	if nil == scheduler.jobPlanManager{
		logs.Info("没设置jobPlanManager对象")
		return
	}
	var (
		jobSchedulePlan *common.JobSchedulePlan
		err             error
		jobExecuteInfo  *common.JobExecuteInfo //执行中的任务
		jobExecting     bool
	)

	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE: //保存任务事件
		if jobSchedulePlan, err = common.BuildJobSchedulePlan(jobEvent.Job); err != nil {
			logs.Error("构建任务执行计划出错:", err.Error())
			return
		}
		logs.Info("保存任务:", jobEvent.Job.Name)
		if err := scheduler.jobPlanManager.Insert(jobSchedulePlan);err != nil{
			logs.Error("保存任务出错:", err.Error())
		}
	case common.JOB_EVENT_DELETE: //删除任务事件
		logs.Warn("删除任务:", jobEvent.Job)
		if err := scheduler.jobPlanManager.Remove(jobEvent.Job.Name,jobEvent.Job);err != nil{
			logs.Error("删除任务出错:", err.Error())
		}
	case common.JOB_EVENT_KILL: //强杀任务事件
		//取消command的执行
		if jobExecuteInfo, jobExecting = scheduler.jobExecutingTable[jobEvent.Job.Name]; jobExecting {
			logs.Info("强杀任务:", jobEvent.Job.Name)
			jobExecuteInfo.CancelFunc() //触发command杀死shell
		}
	}
}

//尝试执行任务
func (scheduler *Scheduler) TryStartJob(jobPlan *common.JobSchedulePlan) (err error) {
	//调度
	if nil == scheduler.jobExecuter  {
		return errors.New("还没有设置任务的执行器")
	}
	if nil == jobPlan{
		return errors.New("参数jobPlan不应该为空")
	}
	//执行的任务可能运行很久,1分钟会调度60次,但是只能执行1次,防止并发
	//如果任务正在执行,跳过本次调度
	if _, jobExecuting := scheduler.jobExecutingTable[jobPlan.Job.Name]; jobExecuting {
		return errors.New("尚未退出,跳过执行")
	}

	//不存在,则构建一个
	jobExecuteInfo := common.BuildJobExecuteInfo(jobPlan)
	//1,记录执行状态信息
	scheduler.jobExecutingTable[jobPlan.Job.Name] = jobExecuteInfo

	//2,执行任务, 同步或异步, 取决于 jobExecuter的实现
	err = scheduler.jobExecuter.Exec(jobExecuteInfo)
	//TODO:3,等待任务执行结果的通知,使用回调函数?还是填写数据到jobResultChan?

	return err

}

//重新计算任务调度状态
func (scheduler *Scheduler) TrySchedule() (scheduleAfter time.Duration) {

	//如果任务表为空,随便睡眠多久
	if nil == scheduler.jobPlanManager || 0 == scheduler.jobPlanManager.Size() {
		//fmt.Println("无任务被调度")
		scheduleAfter = 1 * time.Second
		return
	}

	//查找最早的任务,并传入  scheduler.TryStartJob 执行
	scheduleAfter ,_= scheduler.jobPlanManager.ExtractEarliest(scheduler.TryStartJob);
	return

}

//调度协程
func (scheduler *Scheduler) scheduleLoop() {
	var (
		jobEvent      *common.JobEvent
		scheduleAfter time.Duration
		scheduleTimer *time.Timer
		jobResult     *common.JobExecuteResult
	)
	//初始化一次(1s)
	scheduleAfter = scheduler.TrySchedule()
	//调度的延迟定时器
	scheduleTimer = time.NewTimer(scheduleAfter)
	//定时任务common.Job
	for {
		select {
		case jobEvent = <-scheduler.jobEventChan: //监听任务变化事件
			scheduler.handleJobEvent(jobEvent)
		case <-scheduleTimer.C: //最近的任务到期了
		case jobResult = <-scheduler.jobResultChan: //监听任务执行结果
			scheduler.HandleJobResult(jobResult)
		}
		//调度一次任务
		scheduleAfter = scheduler.TrySchedule()
		//重置调度间隔
		scheduleTimer.Reset(scheduleAfter)
	}
}

//推送etcd任务变化事件
func (scheduler *Scheduler) PushEvent(jobEvent *common.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}

//回传任务执行结果
func (scheduler *Scheduler) PushResult(jobResult *common.JobExecuteResult) {
	scheduler.jobResultChan <- jobResult
}

//JobEventChan etcd任务事件队列的数量
func (scheduler *Scheduler) JobEventChanLen() int {
	return len(scheduler.jobEventChan)
}

//回传任务执行结果 队列的数量
func (scheduler *Scheduler) JobResultChanLen() int {
	return len(scheduler.jobResultChan)
}

//启动协程
func (scheduler *Scheduler) Loop() {
	//启动协程
	logs.Info("启动任务调度协程")
	go G_scheduler.scheduleLoop()
}

//处理任务结果,记录任务的执行时间,计划时间,输出结果
func (scheduler *Scheduler) HandleJobResult(result *common.JobExecuteResult)(n int, err error) {

	if nil == result || result.ExecuteInfo == nil {
		return 0, errors.New("日志对象不能为空")
	}
	var (
		jobLog *common.JobLog
	)
	//删除执行状态
	delete(scheduler.jobExecutingTable, result.ExecuteInfo.Job.Name)

	result.PrintSuccessLog()
	//生成执行日志
	jobLog = result.ParseJobLog()

	if nil != scheduler.jobLogger && nil != jobLog{
		//发送给日志记录器 
		dd := []*common.JobLog{jobLog}
		logs.Warn("发送jobLog给日志记录器",&jobLog," ,jobLog[]=",&dd)
		return scheduler.jobLogger.Write(&common.LogBatch{dd})
	}else{
		return 0, errors.New("没设置日志记录器")
	}

}

//单例
var (
	G_scheduler   *Scheduler
	oncescheduler sync.Once
)

//初始化调度器
func InitScheduler(jobLogger log.JobLoger) (err error, scheduler *Scheduler) {
	oncescheduler.Do(func() {

		G_scheduler = &Scheduler{
			jobEventChan:      make(chan *common.JobEvent, 1000), //有缓冲区?
			jobExecutingTable: make(map[string]*common.JobExecuteInfo),
			jobResultChan:     make(chan *common.JobExecuteResult, 1000),
		}

	})
	G_scheduler.jobLogger = jobLogger

	scheduler = G_scheduler
	return
}
