package common

import (
	"context"
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"strings"
	"time"
	"github.com/astaxie/beego/logs"
)

//定时任务
type Job struct {
	Name     string `json:"name"`     //任务名
	Command  string `json:"command"`  //shell命令
	ShellName  string `json:"shellName"`  //shell命令
	CronExpr string `json:"cronExpr"` //cron表达式
}

//任务调度计划
type JobSchedulePlan struct {
	Job      *Job                 //要调度的任务信息
	Expr     *cronexpr.Expression //解释好的cronexpr表达式
	NextTime time.Time            //下次调度时间
	Del bool
}

//任务执行状态
type JobExecuteInfo struct {
	Job        *Job               //任务信息
	PlanTime   time.Time          //理论上的调度时间
	RealTime   time.Time          //实际的调度时间
	CancelCtx  context.Context    //取消command的context
	CancelFunc context.CancelFunc //取消command的函数
}

func (this *JobExecuteResult)PrintSuccessLog()  {
	logs.Info("任务执行完成:", this.ExecuteInfo.Job.Name, " Es=",
		this.EndTime.Sub(this.StartTime),
		string(this.Output), " Err=", this.Err)
}
//使用JobExecuteResult构建JobLog对象
func (this *JobExecuteResult)ParseJobLog() (*JobLog) {
	//过滤无用 的日志
	if this.Err != ERR_LOCK_ALREADY_REQUIRED {
		return nil
	}
	jobLog := &JobLog{
		JobName:      this.ExecuteInfo.Job.Name,
		Command:      this.ExecuteInfo.Job.Command,
		Output:       string(this.Output),
		PlanTime:     this.ExecuteInfo.PlanTime.UnixNano() / 1000000,
		ScheduleTime: this.ExecuteInfo.RealTime.UnixNano() / 1000000,
		StartTime:    this.StartTime.UnixNano() / 1000000,
		EndTime:      this.EndTime.UnixNano() / 1000000,
	}
	if this.Err != nil {
		jobLog.Err = this.Err.Error()
	} else {
		jobLog.Err = ""
	}
	return  jobLog
}

func (j *JobExecuteInfo)PrintStatus()  {
	logs.Info("正式执行任务:", j.Job.Name, " P=", j.PlanTime, " R=", j.RealTime)
}

type Response struct {
	Errno int         `json:"errno"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data"`
}

//变化事件
type JobEvent struct {
	EventType int //save delete
	Job       *Job
}

//任务执行结果
type JobExecuteResult struct {
	ExecuteInfo *JobExecuteInfo //执行状态
	Output      []byte          //脚本输出
	Err         error           //脚本错误原因
	StartTime   time.Time       //启动时间
	EndTime     time.Time       //结束时间
}

//
type JobLog struct {
	JobName      string `json:"jobName" bson:"jobName"`      //任务名字
	Command      string `json:"command" bson:"command"`      //脚本命令
	Err          string `json:"err" bson:"err"`          //错误原因
	Output       string `json:"output" bson:"output"`       //脚本输出
	PlanTime     int64  `json:"planTime" bson:"planTime"`     //计划开始时间
	ScheduleTime int64  `json:"scheduleTime" bson:"scheduleTime"` //实际调度时间
	StartTime    int64  `json:"startTime" bson:"startTime"`    //任务执行开始时间
	EndTime      int64  `json:"endTime" bson:"endTime"`      //任务执行结束时间
}

//日志批次
type LogBatch struct {
	Logs []*JobLog //多条日志
}
//任务日志过滤条件
type JobLogFilter struct {
	JobName string `bson:"jobName"`
}
//任务日志排序规则
type SortLogByStartTime struct {
	SortOrder int `bson:"startTime"` //{startTime:-1}
}

//应答方法
func BuildResponse(errno int, msg string, data interface{}) (resp []byte, err error) {
	//定义一个response
	response := &Response{
		Errno: errno,
		Msg:   msg,
		Data:  data,
	}
	//序列化json
	resp, err = json.Marshal(response)
	return
}

//反序列化job
func UnpackJob(value []byte) (ret *Job, err error) {
	job := &Job{}
	err = json.Unmarshal(value, job)
	if err != nil {
		return
	}
	ret = job
	return
}
//序列化job
func PackJob(job *Job) (b []byte, err error) {

	b,err = json.Marshal(job)
	if err != nil {
		return nil,err
	}
	return b,nil
}

//从etcd的key中提取任务名称
// /cron/jobs/job10 => job10
func ExtractJobName(jobKey string) string {
	return strings.TrimPrefix(jobKey, JOB_SAVE_DIR)
}

//从etcd的key中提取任务名称
// /cron/killer/job10 => job10
func ExtractKillerName(jobKey string) string {
	return strings.TrimPrefix(jobKey, JOB_KILLER_DIR)
}

//提取worker的ip
func ExtractWorkerIP(jobKey string) string {
	return strings.TrimPrefix(jobKey, JOB_WORKER_DIR)
}

//任务变化事件有2种, 1,更新任务, 2,删除任务
func BuildJobEvent(eventType int, job *Job) (jobEvent *JobEvent) {
	return &JobEvent{
		EventType: eventType,
		Job:       job,
	}
}

//构建任务执行计划
func BuildJobSchedulePlan(job *Job) (jobSchedulePlan *JobSchedulePlan, err error) {

	var (
		expr *cronexpr.Expression
	)
	//解析JOB的crontab表达式
	if expr, err = cronexpr.Parse(job.CronExpr); err != nil {
		logs.Error("解析JOB的crontab表达式 出错")
		return
	}
	//生成任务调度计划对象
	jobSchedulePlan = &JobSchedulePlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()), //根据当前时间,计算下次时间
		Del:false,
	}
	return
}

//构造执行状态信息
func BuildJobExecuteInfo(jobSchedulePlan *JobSchedulePlan) (jobExecuteInfo *JobExecuteInfo) {
	jobExecuteInfo = &JobExecuteInfo{
		Job:      jobSchedulePlan.Job,
		PlanTime: jobSchedulePlan.NextTime, //计算调度时间
		RealTime: time.Now(),               //真实调度时间
	}
	jobExecuteInfo.CancelCtx, jobExecuteInfo.CancelFunc = context.WithCancel(context.TODO())
	return
}
