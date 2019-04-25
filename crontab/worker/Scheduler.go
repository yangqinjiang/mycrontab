package worker

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"sync"
	"time"
)

/**
 调度器,遍历所有任务列表, 找出最近一个要过期的任务
 */
type Scheduler struct {
	jobEventChan      chan *common.JobEvent              //etcd任务事件队列
	jobPlanTable      map[string]*common.JobSchedulePlan //任务调度计划表内存里的任务计划表,
	jobExecutingTable map[string]*common.JobExecuteInfo  //任务执行表
	jobResultChan     chan *common.JobExecuteResult      //任务执行结果队列
	jobLogger JobLogger
	jobExecuter JobExecuter
}
/**
设置任务的执行器
 */
func (scheduler *Scheduler)SetJobExecuter(jobExecuter JobExecuter)  {
	scheduler.jobExecuter = jobExecuter
}



//处理任务事件
func (scheduler *Scheduler) handleJobEvent(jobEvent *common.JobEvent) {
	var (
		jobSchedulePlan *common.JobSchedulePlan
		jobExisted      bool //是否存在任务
		err             error
		jobExecuteInfo  *common.JobExecuteInfo //执行中的任务
		jobExecting     bool
	)

	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE: //保存任务事件
		if jobSchedulePlan, err = common.BuildJobSchedulePlan(jobEvent.Job); err != nil {
			return
		}
		logs.Info("保存任务:", jobEvent.Job.Name)
		scheduler.jobPlanTable[jobEvent.Job.Name] = jobSchedulePlan
	case common.JOB_EVENT_DELETE: //删除任务事件
		//首先检查是否已存在此任务 (内存)
		if jobSchedulePlan, jobExisted = scheduler.jobPlanTable[jobEvent.Job.Name]; jobExisted {
			logs.Info("删除任务:", jobEvent.Job.Name)
			delete(scheduler.jobPlanTable, jobEvent.Job.Name)
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
func (scheduler *Scheduler) TryStartJob(jobPlan *common.JobSchedulePlan)(err error) {
	//调度
	if scheduler.jobExecuter == nil{
		return errors.New("还没有设置任务的执行器")
	}
	//执行的任务可能运行很久,1分钟会调度60次,但是只能执行1次,防止并发
	//如果任务正在执行,跳过本次调度
	if _, jobExecuting := scheduler.jobExecutingTable[jobPlan.Job.Name]; jobExecuting {
		return errors.New("尚未退出,跳过执行")
	}

	//不存在,则构建一个
	jobExecuteInfo := common.BuildJobExecuteInfo(jobPlan)
	//记录执行状态信息
	scheduler.jobExecutingTable[jobPlan.Job.Name] = jobExecuteInfo

	//执行任务
	err =scheduler.jobExecuter.Exec(jobExecuteInfo)

	return err

}


//重新计算任务调度状态
func (scheduler *Scheduler) TrySchedule() (scheduleAfter time.Duration) {
	var (
		jobPlan  *common.JobSchedulePlan
		now      time.Time
		nearTime *time.Time //指针
	)
	//如果任务表为空,随便睡眠多久
	if len(scheduler.jobPlanTable) == 0 {
		//fmt.Println("无任务被调度")
		scheduleAfter = 1 * time.Second
		return
	}
	//当前时间
	now = time.Now()
	//TODO: 使用最小堆， 动态排序任务
	//  参考资料
	//  https://github.com/liuyubobobo/Play-with-Algorithms/blob/master/04-Heap/Course%20Code%20(C%2B%2B)/Optional-2-Min-Heap/MinHeap.h
	//1遍历所有任务
	for _, jobPlan = range scheduler.jobPlanTable {
		//是否过期,小于或都等于当前时间
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			//尝试执行任务
			scheduler.TryStartJob(jobPlan)
			jobPlan.NextTime = jobPlan.Expr.Next(now) //执行后,更新下次执行时间的值

		}
		//统计最近一个要过期的任务时间
		if nearTime == nil{ // 刚开始for第一个,则设置一个值
			nearTime = &jobPlan.NextTime
		}
		//判断第二个及以后,如果是更往后的时刻,则更新它,
		// 即 找出更晚的时刻
		if jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}
	//下次调度间隔,(最近要执行的任务调度时间 - 当前时间)
	scheduleAfter = (*nearTime).Sub(now)
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
			scheduler.handleJobResult(jobResult)
		}
		//调度一次任务
		scheduleAfter = scheduler.TrySchedule()
		//重置调度间隔
		scheduleTimer.Reset(scheduleAfter)
	}
}

//推送etcd任务变化事件
func (scheduler *Scheduler) PushJobEvent(jobEvent *common.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}
//JobEventChan etcd任务事件队列的数量
func (scheduler *Scheduler)JobEventChanLen() int  {
	return len(scheduler.jobEventChan)
}


//回传任务执行结果
func (scheduler *Scheduler) PushJobResult(jobResult *common.JobExecuteResult) {
	scheduler.jobResultChan <- jobResult
}
//回传任务执行结果 队列的数量
func (scheduler *Scheduler)JobResultChanLen() int  {
	return len(scheduler.jobResultChan)
}
//启动协程
func (scheduler *Scheduler)Loop()  {
	//启动协程
	logs.Info("启动调度协程")
	go G_scheduler.scheduleLoop()
}
//处理任务结果,记录任务的执行时间,计划时间,输出结果
func (scheduler *Scheduler) handleJobResult(result *common.JobExecuteResult) {

	if result.ExecuteInfo == nil{
		return
	}
	var (
		jobLog *common.JobLog
	)
	//删除执行状态
	delete(scheduler.jobExecutingTable, result.ExecuteInfo.Job.Name)
	logs.Info("任务执行完成:", result.ExecuteInfo.Job.Name, " Es=", result.EndTime.Sub(result.StartTime), string(result.Output), " Err=", result.Err)
	//生成执行日志
	if result.Err != common.ERR_LOCK_ALREADY_REQUIRED {
		jobLog = &common.JobLog{
			JobName:      result.ExecuteInfo.Job.Name,
			Command:      result.ExecuteInfo.Job.Command,
			Output:       string(result.Output),
			PlanTime:     result.ExecuteInfo.PlanTime.UnixNano() / 1000000,
			ScheduleTime: result.ExecuteInfo.RealTime.UnixNano() / 1000000,
			StartTime:    result.StartTime.UnixNano() / 1000000,
			EndTime:      result.EndTime.UnixNano() / 1000000,
		}
		if result.Err != nil {
			jobLog.Err = result.Err.Error()
		} else {
			jobLog.Err = ""
		}
		if scheduler.jobLogger != nil{
			//发送给日志记录器
			scheduler.jobLogger.Write(jobLog)
		}

	}
}


//单例
var (
	G_scheduler *Scheduler
	oncescheduler        sync.Once
)


//初始化调度器
func InitScheduler(jobLogger JobLogger) (err error,scheduler *Scheduler) {
	oncescheduler.Do(func() {

		G_scheduler = &Scheduler{
			jobEventChan:      make(chan *common.JobEvent, 1000),              //有缓冲区?
			jobPlanTable:      make(map[string]*common.JobSchedulePlan, 1000), //内存里的任务计划表,
			jobExecutingTable: make(map[string]*common.JobExecuteInfo),
			jobResultChan:     make(chan *common.JobExecuteResult,1000),
		}

	})
	G_scheduler.jobLogger = jobLogger
	scheduler = G_scheduler
	return
}
