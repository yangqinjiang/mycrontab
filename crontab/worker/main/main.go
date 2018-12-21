package main

import (
	"flag"
	"fmt"
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
	fmt.Println("crontab worker running...")
	var (
		err error
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
	fmt.Println("读取配置文件")

	err = worker.InitLogSink()
	if err != nil {
		goto ERR
	}
	fmt.Println("初始化mongodb的实例")
	//启动任务执行器
	err = worker.InitExecutor()
	if err != nil {
		goto ERR
	}
	fmt.Println("启动任务执行器")
	//启动任务调度器
	err = worker.InitScheduler()
	if err != nil {
		goto ERR
	}
	fmt.Println("启动任务调度器")
	//启动任务管理器
	err = worker.InitJobMgr()
	if err != nil {
		goto ERR
	}
	fmt.Println("启动任务管理器")
	//启动服务注册管理器
	err = worker.InitRegistr()
	if err != nil {
		goto ERR
	}
	fmt.Println("启动服务注册管理器")
	fmt.Println("worker启动完成")
	//正常退出
	for {
		time.Sleep(1 * time.Second)
	}
	return

	//异常退出
ERR:
	fmt.Println("worker启动失败:",err)
}
