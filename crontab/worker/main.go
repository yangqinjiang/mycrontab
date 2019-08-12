package main

import (
	"errors"
	"flag"
	logs "github.com/sirupsen/logrus"
	"runtime"
	// "strconv"
	"time"
	"github.com/yangqinjiang/mycrontab/worker/lib"
	// "github.com/yangqinjiang/mycrontab/worker/common"
)

var (
	confFile string //配置文件的路径
)
//TODO: 创建命令doctor,用于检查运行环境,例如连接etcd, 连接mongodb,等等
//解析命令行参数
//TODO:在 goland IDE里启动,需要替换working directory
///src/github.com/yangqinjiang/mycrontab/crontab/worker/main
func initArgs() {
	flag.StringVar(&confFile, "config", "./config/worker.json", "指定worker.json")
	flag.Parse()
}
func InitEnv() {
	//线程数==CPU数量
	runtime.GOMAXPROCS(runtime.NumCPU())
}
//初始化logs的配置
func initLogs(env_production bool) {

	// do something here to set environment depending on an environment variable
	// or command-line flag
	if env_production {
		//日志打印 代码调用的路径 
		logs.SetReportCaller(true)
		logs.SetFormatter(&logs.JSONFormatter{})
	} else {
	  // The TextFormatter is default, you don't actually have to do this.
		logs.SetFormatter(&logs.TextFormatter{})
	}
  }
func main() {
	logs.Info("Crontab Worker Starting...")
	var (
		err error
		jobEventPusher *lib.CustomJobEventReceiver
	)
	//初始化命令行参数
	initArgs()

	//初始化线程
	InitEnv()

	//加载配置
	err = lib.InitConfig(confFile)
	if err != nil {
		goto ERR
	}
	logs.Info("加载配置")

	initLogs(lib.G_config.LogsProduction)
	//-----------------------日志记录器的实现------------------------
	//var testWriter *lib.ConsoleLog
	// testWriter = &lib.ConsoleLog{}
	// logs.Info("init ConsoleLog")
	//err = lib.InitJobLogMemoryBuffer(testWriter)

	//使用 MongoDB 储存任务产生的日志
	err = lib.InitMongoDbLog()
	if err != nil {
		goto ERR
	}
	if nil == lib.G_MongoDbLog{
		err = errors.New("初始化mongodb的实例  数据库连接 [失败]")
		goto ERR
	}else{
		logs.Info("初始化mongodb的实例,储存任务产生的日志 [完成]")
	}
	
	err = lib.InitJobLogMemoryBuffer(lib.G_MongoDbLog)
	if err != nil {
		goto ERR
	}
	logs.Info("init JobLogMemoryBuffer")

	//启动异步任务执行器
	err = lib.InitGoroutineExecutor()
	if err != nil {
		goto ERR
	}
	logs.Info("启动异步任务执行器 InitGoroutineExecutor")
	//------------------任务管理器-----------------------------
	//启动  任务管理器 监听 etcd 的事件, 组装任务数据, 并推给 scheduler任务调度器
	err = lib.InitEtcdJobMgr()
	if err != nil {
		goto ERR
	}
	if nil == lib.G_EtcdJobMgr{
		err = errors.New("Etcd 任务存储 数据库连接 [失败]")
		goto ERR
	}else{
		logs.Info("Etcd 任务存储 数据库连接 [完成] ")
	}

	//启动任务调度器
	err,_ = lib.InitScheduler(nil)
	if err != nil {
		goto ERR
	}

	//设置 推送任务事件 的操作者
	jobEventPusher = &lib.CustomJobEventReceiver{JobEventReceiver: lib.G_scheduler}
	if nil == jobEventPusher {
		err = errors.New("jobEventPusher nil pointer")
		goto ERR
	}

	lib.G_EtcdJobMgr.SetJobEventPusher(jobEventPusher)
	//设置任务执行结果的接收器
	lib.G_GoroutineExecutor.SetJobResultReceiver(lib.G_scheduler)
	//----------------------任务调度器--------------------------
	// 使用 [ 任务管理器推给的任务数据 ],经过 [JobPlanManager调度时间排序] 得到最先应该执行的任务,
	// 再[同步或JobExecuter异步执行],最后 使用[JobLogger记录任务的执行日志]


	//设置任务调度器的日志记录器
	lib.G_scheduler.SetJobLogBuffer(lib.G_jobLogMemoryBuffer)
	//设置任务调度器的任务执行器  -> goroutine的任务执行器
	lib.G_scheduler.SetJobExecuter(lib.G_GoroutineExecutor)
	//设置 任务调度时间  的计算算法
	lib.G_scheduler.SetJobPlanManager(lib.NewJobPlanMinHeap(10000))
	//启动任务调度器的 调度协程,监听任务变化事件,任务执行结果
	lib.G_scheduler.Loop()
	logs.Info("启动任务调度器的 调度协程 [完成]")


	//---------------------服务注册管理器------------------
	//启动服务注册管理器, 保持在线
	err = lib.InitRegistr()
	if err != nil {
		goto ERR
	}
	logs.Info("启动服务注册管理器 [完成]")
	logs.Info("启动worker[完成]")
	logs.Info("running....")
	//正常退出
	for {
		time.Sleep(1 * time.Second)
	}
	return

	//异常退出
ERR:
	logs.Error("start worker  ERROR:", err)
}
