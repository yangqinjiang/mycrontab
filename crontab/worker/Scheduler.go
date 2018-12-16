package worker

import (
	"fmt"
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"time"
)

type Scheduler struct {
	jobEventChan chan *common.JobEvent              //etcd任务事件队列
	jobPlanTable map[string]*common.JobSchedulePlan //任务调度计划表内存里的任务计划表,
	jobExecutingTable map[string]*common.JobExecuteInfo//任务执行表
}

//单例
var (
	G_scheduler *Scheduler
)

//处理任务事件
func (scheduler *Scheduler) handleJobEvent(jobEvent *common.JobEvent) {
	var (
		jobSchedulePlan *common.JobSchedulePlan
		jobExisted      bool //是否存在任务
		err             error
	)

	switch jobEvent.EventType {
	case common.JOB_EVENT_SAVE: //保存任务事件
		if jobSchedulePlan, err = common.BuildJobSchedulePlan(jobEvent.Job); err != nil {
			return
		}
		scheduler.jobPlanTable[jobEvent.Job.Name] = jobSchedulePlan
	case common.JOB_EVENT_DELETE: //删除任务事件
		//首先检查是否已存在此任务 (内存)
		if jobSchedulePlan, jobExisted = scheduler.jobPlanTable[jobEvent.Job.Name]; jobExisted {
			delete(scheduler.jobPlanTable, jobEvent.Job.Name)
		}
	}
}
//尝试执行任务
func (scheduler *Scheduler) TryStartJob(jobPlan *common.JobSchedulePlan)  {
	//调度和执行是2件事件
	var (
		jobExecuteInfo *common.JobExecuteInfo
		jobExecuting bool
	)
	//执行的任务可能运行很久,1分钟会调度60次,但是只能执行1次,防止并发
	//如果任务正在执行,跳过本次调度
	if _,jobExecuting = scheduler.jobExecutingTable[jobPlan.Job.Name];jobExecuting{
		fmt.Println("尚未退出,跳过执行")
		return
	}
	//不存在,则构建一个
	jobExecuteInfo = common.BuildJobExecuteInfo(jobPlan)

	//构建执行状态信息
	scheduler.jobExecutingTable[jobPlan.Job.Name] = jobExecuteInfo

	//TODO:执行任务

	fmt.Println("正式执行任务:",jobExecuteInfo.Job.Name," P=",jobExecuteInfo.PlanTime," R=",jobExecuteInfo.RealTime)


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
	//1遍历所有任务
	for _, jobPlan = range scheduler.jobPlanTable {
		//是否过期
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			//TODO:尝试执行任务
			scheduler.TryStartJob(jobPlan)
			fmt.Println("执行任务:", jobPlan.Job.Name, " @ ", now.Minute(), ':', now.Second())
			jobPlan.NextTime = jobPlan.Expr.Next(now) //更新下次执行时间

		}
		//统计最近一个要过期的任务时间
		if nearTime == nil || jobPlan.NextTime.Before(*nearTime) {
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

		}
		//调度一次任务
		scheduleAfter = scheduler.TrySchedule()
		//重置调度间隔
		scheduleTimer.Reset(scheduleAfter)
	}
}

//推送任务变化事件
func (scheduler *Scheduler) PUshJobEvent(jobEvent *common.JobEvent) {
	scheduler.jobEventChan <- jobEvent
}

//初始化调度器
func InitScheduler() (err error) {
	G_scheduler = &Scheduler{
		jobEventChan: make(chan *common.JobEvent, 1000),              //有缓冲区?
		jobPlanTable: make(map[string]*common.JobSchedulePlan, 1000), //内存里的任务计划表,
		jobExecutingTable:make(map[string]*common.JobExecuteInfo),
	}
	//启动协程
	go G_scheduler.scheduleLoop()
	return
}
