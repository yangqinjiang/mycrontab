package master

import (
	"net/http"
	"net"
	"time"
	"strconv"
	"encoding/json"
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"fmt"
)

//任务的HTTP接口
type ApiServer struct {
	httpServer *http.Server
}

var (
	//单例对象
	G_apiServer *ApiServer
)

//保存任务接口
//保存任务到ETCD
//post job={"name":"job1",""command:"echo hello","cronExpr":"* * * * *"}
func handleJobSave(resp http.ResponseWriter,req *http.Request) {
	fmt.Println("保存任务接口")
	var (
		err error
		postJob string
		job common.Job
		oldJob *common.Job
		bytes []byte
	)
	err = req.ParseForm()
	if err != nil{
		goto ERR
	}
	//取表单的job字段
	postJob = req.PostForm.Get("job")
	//反序列化
	err= json.Unmarshal([]byte(postJob),&job)
	if err != nil{
		goto ERR
	}
	oldJob ,err = G_jobMgr.SaveJob(&job)
	if err != nil{
		goto ERR
	}
	//返回正常应答{{"errno":0,"msg":"","data":{...}}}
	bytes ,err  =common.BuildResponse(0,"success",oldJob)
	if err != nil{
		goto ERR
	}
	resp.Write(bytes)

	return

	ERR:
		//返回异常应答
		bytes ,_  =common.BuildResponse(0,err.Error(),nil)
		resp.Write(bytes)
}
//初始化服务
func InitApiServer()(err error)  {
	var (
		mux *http.ServeMux
		listener net.Listener
	)
	//配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save",handleJobSave)

	//启动TCP监听
	listener ,err = net.Listen("tcp",":"+strconv.Itoa(G_config.ApiPort))
	if err != nil{
		return
	}

	//创建一个HTTP服务
	httpServer := &http.Server{
		ReadTimeout:time.Duration(G_config.ApiReadTimeout)*time.Second,//超时
		WriteTimeout:time.Duration(G_config.ApiWriteTimeout)*time.Second,
		Handler:mux,
	}
	G_apiServer = &ApiServer{
		httpServer:httpServer,
	}
	//启动服务端
	go httpServer.Serve(listener)
	return
}