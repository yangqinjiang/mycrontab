package common

import "encoding/json"

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

