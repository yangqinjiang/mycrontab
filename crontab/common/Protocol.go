package common

import (
	"encoding/json"
	"strings"
)

//定时任务
type Job struct {
	Name string `json:"name"`//任务名
	Command string `json:"command"` //shell命令
	CronExpr string `json:"cronExpr"`//cron表达式
}

type Response struct {
	Errno  int `json:"errno"`
	Msg string `json:"msg"`
	Data interface{} `json:"data"`
}

//变化事件
type JobEvent struct {
	EventType int//save delete
	Job *Job
}

//应答方法
func BuildResponse(errno int,msg string,data interface{}) (resp []byte,err error)  {
	//定义一个response
	response := &Response{
		Errno:errno,
		Msg:msg,
		Data:data,
	}
	//序列化json
	resp,err = json.Marshal(response)
	return
}

//反序列化job
func UnpackJob(value []byte)(ret *Job,err error)  {
	job := &Job{}
	err = json.Unmarshal(value,job)
	if err != nil{
		return
	}
	ret = job
	return
}

//从etcd的key中提取任务名称
// /cron/jobs/job10 => job10
func ExtractJobName(jobKey string) (string)  {
	return strings.TrimPrefix(jobKey,JOB_SAVE_DIR)
}

//任务变化事件有2种, 1,更新任务, 2,删除任务
func BuildJobEvent(eventType int,job *Job)(jobEvent *JobEvent)  {
	return &JobEvent{
		EventType:eventType,
		Job:job,
	}
}
