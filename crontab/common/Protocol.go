package common

import (
	"encoding/json"
	"github.com/gorhill/cronexpr"
	"strings"
	"time"
)

//定时任务
type Job struct {
	Name     string `json:"name"`     //任务名
	Command  string `json:"command"`  //shell命令
	CronExpr string `json:"cronExpr"` //cron表达式
}

//任务调度计划
type JobSchedulePlan struct {
	Job      *Job                 //要调度的任务信息
	Expr     *cronexpr.Expression //解释好的cronexpr表达式
	NextTime time.Time            //下次调度时间
}

//任务执行状态
type JobExecuteInfo struct {
	Job *Job //任务信息
	PlanTime time.Time //理论上的调度时间
	RealTime time.Time //实际的调度时间
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

//从etcd的key中提取任务名称
// /cron/jobs/job10 => job10
func ExtractJobName(jobKey string) string {
	return strings.TrimPrefix(jobKey, JOB_SAVE_DIR)
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
		return
	}
	//生成任务调度计划对象
	jobSchedulePlan = &JobSchedulePlan{
		Job:      job,
		Expr:     expr,
		NextTime: expr.Next(time.Now()), //根据当前时间,计算下次时间
	}
	return
}

//构造执行状态信息
func BuildJobExecuteInfo(jobSchedulePlan *JobSchedulePlan)(jobExecuteInfo *JobExecuteInfo)  {
	jobExecuteInfo = &JobExecuteInfo{
		Job:jobSchedulePlan.Job,
		PlanTime:jobSchedulePlan.NextTime,//计算调度时间
		RealTime:time.Now(),//真实调度时间
	}
	return
}
