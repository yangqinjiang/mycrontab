package master

import (
	"encoding/json"
	"fmt"
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"net"
	"net/http"
	"strconv"
	"time"
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
func handleJobSave(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("保存任务接口")
	var (
		err     error
		postJob string
		job     common.Job
		oldJob  *common.Job
		bytes   []byte
	)
	err = req.ParseForm()
	if err != nil {
		goto ERR
	}
	//取表单的job字段
	postJob = req.PostForm.Get("job")
	//反序列化
	err = json.Unmarshal([]byte(postJob), &job)
	if err != nil {
		goto ERR
	}
	oldJob, err = G_jobMgr.SaveJob(&job)
	if err != nil {
		goto ERR
	}
	//返回正常应答{{"errno":0,"msg":"","data":{...}}}
	bytes, err = common.BuildResponse(0, "success", oldJob)
	if err != nil {
		goto ERR
	}
	resp.Write(bytes)

	return

ERR:
	//返回异常应答
	bytes, _ = common.BuildResponse(-1, err.Error(), nil)
	resp.Write(bytes)
}

//删除etcd的任务
// post /job/delete name=job1
func handleJobDelete(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("删除etcd的任务")
	var (
		err    error
		name   string
		oldJob *common.Job
		bytes  []byte
	)
	err = req.ParseForm()
	if err != nil {
		goto ERR
	}

	//删除任务名
	name = req.PostForm.Get("name")

	oldJob, err = G_jobMgr.DeleteJob(name)
	if err != nil {
		goto ERR
	}
	//返回正常应答{{"errno":0,"msg":"","data":{...}}}
	bytes, err = common.BuildResponse(0, "success", oldJob)
	if err != nil {
		goto ERR
	}
	resp.Write(bytes)
	return

ERR:
	//返回异常应答
	bytes, _ = common.BuildResponse(-1, err.Error(), nil)
	resp.Write(bytes)
}

//列出任务
func handleJobList(resp http.ResponseWriter, req *http.Request) {
	fmt.Println("列出任务")
	var (
		jobList []*common.Job
		err     error
		bytes   []byte
	)
	jobList, err = G_jobMgr.ListJobs()
	if err != nil {
		goto ERR
	}
	//返回正常应答{{"errno":0,"msg":"","data":{...}}}
	bytes, err = common.BuildResponse(0, "success", jobList)
	if err != nil {
		goto ERR
	}
	resp.Write(bytes)
	return
ERR:
	//返回异常应答
	bytes, _ = common.BuildResponse(-1, err.Error(), nil)
	resp.Write(bytes)
}

//强杀任务
// post /job/kill name=job1
func handleJobKill(resp http.ResponseWriter, req *http.Request) {

	var (
		err   error
		name  string
		bytes []byte
	)
	err = req.ParseForm()
	if err != nil {
		goto ERR
	}
	//要杀死的任务名称
	name = req.PostForm.Get("name")
	err = G_jobMgr.KillJob(name)
	if err != nil {
		goto ERR
	}
	//返回正常应答{{"errno":0,"msg":"","data":{...}}}
	bytes, err = common.BuildResponse(0, "success", nil)
	if err != nil {
		goto ERR
	}
	resp.Write(bytes)
	return

ERR:
	//返回异常应答
	bytes, _ = common.BuildResponse(-1, err.Error(), nil)
	resp.Write(bytes)
}

//初始化服务
func InitApiServer() (err error) {
	var (
		mux      *http.ServeMux
		listener net.Listener
		staticDir http.Dir  //静态文件根目录
		staticHandler http.Handler //静态文件的HTTP回调
	)
	//配置路由
	mux = http.NewServeMux()
	mux.HandleFunc("/job/save", handleJobSave)
	mux.HandleFunc("/job/delete", handleJobDelete)
	mux.HandleFunc("/job/list", handleJobList)
	mux.HandleFunc("/job/kill", handleJobKill)

	//静态文件目录
	staticDir = http.Dir(G_config.WebRoot)
	staticHandler = http.FileServer(staticDir)
	// /index.html -> index.html  -> ./webroot/index.htmlß
	mux.Handle("/",http.StripPrefix("/",staticHandler)) //匹配最长的 pattern


	//启动TCP监听
	listener, err = net.Listen("tcp", ":"+strconv.Itoa(G_config.ApiPort))
	if err != nil {
		return
	}

	//创建一个HTTP服务
	httpServer := &http.Server{
		ReadTimeout:  time.Duration(G_config.ApiReadTimeout) * time.Second, //超时
		WriteTimeout: time.Duration(G_config.ApiWriteTimeout) * time.Second,
		Handler:      mux,
	}
	G_apiServer = &ApiServer{
		httpServer: httpServer,
	}
	//启动服务端
	go httpServer.Serve(listener)
	return
}
