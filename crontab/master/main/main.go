package main

import (
	"runtime"
	"github.com/yangqinjiang/mycrontab/crontab/master"
	"fmt"
	"flag"
)

var (
	confFile string //配置文件的路径
)
//解析命令行参数
//TODO:在 goland IDE里启动,需要替换working directory
func initArgs()  {
	flag.StringVar(&confFile,"config","./master.json","指定master.json")
	flag.Parse()
}
func InitEnv() {
	//线程数==CPU数量
	runtime.GOMAXPROCS(runtime.NumCPU())
}
func main() {
	var (
		err error
	)
	//初始化命令行参数
	initArgs()

	//初始化线程
	InitEnv()

	//加载配置
	err = master.InitConfig(confFile)
	if err != nil{
		goto ERR
	}

	//启动任务管理器
	err = master.InitJobMgr()
	if err != nil{
		goto ERR
	}

	//启动Api Http服务
	err = master.InitApiServer()
	if err != nil{
		goto ERR//启动出错,直接跳出
	}

	//正常退出
	return

	//异常退出
	ERR:
		fmt.Println(err)
}