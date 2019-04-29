package main

import (
	"flag"
	"github.com/astaxie/beego/logs"
	"github.com/yangqinjiang/mycrontab/crontab/worker"
	"runtime"
	"time"
)

var (
	confFile string //配置文件的路径
)

//解析命令行参数
//TODO:在 goland IDE里启动,需要替换working directory
///src/github.com/yangqinjiang/mycrontab/crontab/worker/main
func initArgs() {
	flag.StringVar(&confFile, "config", "./worker.json", "指定worker.json")
	flag.Parse()
}
func InitEnv() {
	//线程数==CPU数量
	runtime.GOMAXPROCS(runtime.NumCPU())
}
func main() {
	logs.Info("crontab worker running...")
	var (
		err error
		testWriter *worker.ConsoleLog
	)
	//初始化命令行参数
	initArgs()

	//初始化线程
	InitEnv()

	//加载配置
	err = worker.InitConfig(confFile)
	if err != nil {
		goto ERR
	}
	logs.Info("读取配置文件")
	//-----------------------日志记录器的实现------------------------
	//TODO:暂时不使用 MongoDB
	//err = worker.InitMongoDbLog()
	//if err != nil {
	//	goto ERR
	//}
	//logs.Info("初始化mongodb的实例")
	testWriter = &worker.ConsoleLog{}

	err = worker.InitLogSink(testWriter)
	if err != nil {
		goto ERR
	}
	logs.Info("初始化LogSink的实例")
	//启动任务执行器
	err = worker.InitExecutor()
	if err != nil {
		goto ERR
	}
	logs.Info("启动任务执行器")
	//------------------任务管理器-----------------------------
	//启动  任务管理器 监听 etcd 的事件, 组装任务数据, 并推给 scheduler任务调度器
	err = worker.InitJobMgr()
	if err != nil {
		goto ERR
	}
	logs.Info("启动任务管理器")

	//启动任务调度器
	err,_ = worker.InitScheduler(nil)
	if err != nil {
		goto ERR
	}
	//----------------------任务调度器--------------------------
	// 使用 [ 任务管理器推给的任务数据 ],经过 [JobPlanManager调度时间排序] 得到最先应该执行的任务,
	// 再[同步或JobExecuter异步执行],最后 使用[JobLogger记录任务的执行日志]


	//设置任务调度器的日志记录器
	worker.G_scheduler.SetJobLogger(worker.G_jobLogBuffer)
	//设置任务调度器的任务执行器  -> goroutine的任务执行器
	worker.G_scheduler.SetJobExecuter(worker.G_GoroutineExecutor)
	//设置 任务调度时间  的计算算法
	worker.G_scheduler.SetJobPlanManager(worker.NewJobPlanArray())
	//启动任务调度器的 调度协程,监听任务变化事件,任务执行结果
	worker.G_scheduler.Loop()
	logs.Info("启动任务调度器")


	//---------------------服务注册管理器------------------
	//启动服务注册管理器, 保持在线
	err = worker.InitRegistr()
	if err != nil {
		goto ERR
	}
	logs.Info("启动服务注册管理器")
	logs.Info("worker启动完成")
	//正常退出
	for {
		time.Sleep(1 * time.Second)
	}
	return

	//异常退出
ERR:
	logs.Error("worker启动失败:", err)
}
