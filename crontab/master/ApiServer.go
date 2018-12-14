package master

import (
	"net/http"
	"net"
	"time"
	"strconv"
	"encoding/json"
	"github.com/yangqinjiang/mycrontab/crontab/common"
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
func handleJobSave(w http.ResponseWriter,req *http.Request) {

	var (
		err error
		postJob string
		job common.Job
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

	return

	ERR:
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