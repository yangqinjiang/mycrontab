package main

import (
	"flag"
	"fmt"
	"github.com/yangqinjiang/mycrontab/crontab/master"
	"runtime"
	"time"
)

var (
	confFile string //配置文件的路径
)

//解析命令行参数
//TODO:在 goland IDE里启动,需要替换working directory
///src/github.com/yangqinjiang/mycrontab/crontab/master/main
func initArgs() {
	flag.StringVar(&confFile, "config", "./master.json", "指定master.json")
	flag.Parse()
}
func InitEnv() {
	//线程数==CPU数量
	runtime.GOMAXPROCS(runtime.NumCPU())
}
func main() {
	fmt.Println("crontab master running...")
	var (
		err error
	)
	//初始化命令行参数
	initArgs()

	//初始化线程
	InitEnv()

	//加载配置
	err = master.InitConfig(confFile)
	if err != nil {
		goto ERR
	}
	fmt.Println("读取配置文件")

	//启动任务管理器
	err = master.InitJobMgr()
	if err != nil {
		goto ERR
	}
	fmt.Println("启动任务管理器")
	//启动日志管理器
	err = master.InitLogMgr()
	if err != nil {
		goto ERR
	}
	fmt.Println("启动日志管理器")
	//启动服务发现
	err = master.InitWorkerMgr()
	if err != nil {
		goto ERR
	}
	fmt.Println("启动服务发现")


	//启动Api Http服务
	err = master.InitApiServer()
	if err != nil {
		goto ERR //启动出错,直接跳出
	}
	fmt.Println("启动Api Http服务")
	fmt.Println("master启动完成")
	//正常退出
	for {
		time.Sleep(1 * time.Second)
	}
	return

	//异常退出
ERR:
	fmt.Println("master启动失败:",err)
}
