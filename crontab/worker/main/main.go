package main

import (
	"runtime"
	"fmt"
	"flag"
	"time"
	"github.com/yangqinjiang/mycrontab/crontab/worker"
)

var (
	confFile string //配置文件的路径
)
//解析命令行参数
//TODO:在 goland IDE里启动,需要替换working directory
///src/github.com/yangqinjiang/mycrontab/crontab/worker/main
func initArgs()  {
	flag.StringVar(&confFile,"config","./worker.json","指定worker.json")
	flag.Parse()
}
func InitEnv() {
	//线程数==CPU数量
	runtime.GOMAXPROCS(runtime.NumCPU())
}
func main() {
	fmt.Println("worker running...")
	var (
		err error
	)
	//初始化命令行参数
	initArgs()

	//初始化线程
	InitEnv()

	//加载配置
	err = worker.InitConfig(confFile)
	if err != nil{
		goto ERR
	}

	//启动任务管理器
	err = worker.InitJobMgr()
	if err != nil{
		goto ERR
	}

	//正常退出
	for{
		time.Sleep(1*time.Second)
	}
	return

	//异常退出
	ERR:
		fmt.Println(err)
}
