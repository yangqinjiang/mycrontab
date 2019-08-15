package main

import (
	"flag"
	"github.com/yangqinjiang/mycrontab/master/lib"
	"runtime"
	"time"
	"errors"
	logs "github.com/sirupsen/logrus"
)

var (
	confFile string //配置文件的路径
)

//解析命令行参数
//TODO:在 goland IDE里启动,需要替换working directory
///src/github.com/yangqinjiang/mycrontab/crontab/master/main
func initArgs() {
	flag.StringVar(&confFile, "config", "./config/master.json", "指定master.json")
	flag.Parse()
}
func InitEnv() {
	//线程数==CPU数量
	runtime.GOMAXPROCS(runtime.NumCPU())
}
func main() {
	
	logs.Info("Crontab Master Running...")

	var (
		err error
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
	logs.Info("读取配置文件[成功]")

	//启动任务管理器
	err = lib.InitJobMgr()
	if err != nil {
		goto ERR
	}
	if nil == lib.G_jobMgr{
		err = errors.New(" 连接 ETCD 数据库出错,初始化 LogMgr实例 [失败]")
		goto ERR
	}
	logs.Info("启动任务管理器[成功]")
	//启动日志管理器
	err = lib.InitLogMgr()
	if err != nil {
		goto ERR
	}
	if nil == lib.G_jobMgr{
		err = errors.New(" 连接  mongodb 数据库出错,初始化 G_logMgr 实例 [失败]")
		goto ERR
	}
	logs.Info("启动日志管理器[成功]")
	//启动服务发现
	err = lib.InitWorkerMgr()
	if err != nil {
		goto ERR
	}
	if nil == lib.G_workerMgr{
		err = errors.New(" 连接 ETCD 数据库出错,初始化 G_workerMgr 实例 [失败]")
		goto ERR
	}
	logs.Info("启动服务发现[成功]")


	//启动Api Http服务
	err = lib.InitApiServer()
	if err != nil {
		goto ERR //启动出错,直接跳出
	}
	logs.Info("启动Api Http服务[成功]\nMaster启动完成.正常待机")
	//休息一秒
	for {
		time.Sleep(1 * time.Second)
	}
	return

	//异常退出
ERR:
	logs.Error("Master启动失败:"+err.Error())
}
