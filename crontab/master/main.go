package main

import (
	"errors"
	"flag"
	"fmt"
	logs "github.com/sirupsen/logrus"
	"github.com/yangqinjiang/mycrontab/master/lib"
	"os"
	"runtime"
	"time"
)

var (
	help     bool
	quiet    bool   //日志安静模式,不输出info级别的日志
	confFile string //配置文件的路径
	version  bool
)
var version_str string = "1.0"

//解析命令行参数
//TODO:在 goland IDE里启动,需要替换working directory
///src/github.com/yangqinjiang/mycrontab/crontab/master/main
func initFlag() {
	flag.BoolVar(&help, "h", false, "help ,github code source => https://github.com/yangqinjiang/mycrontab ")
	flag.StringVar(&confFile, "c", "./config/master.json", "指定master.json")
	flag.BoolVar(&quiet, "q", false, "quiet,Only log the warning severity or above.")
	flag.BoolVar(&version, "v", false, "Print version infomation and quit")
	flag.Usage = usage
}
func usage() {
	fmt.Fprintf(os.Stderr, fmt.Sprintf(`mycrontab master version: %s
		Usage: master [-hqv] [-c filename]
		Options:
		`, version_str))
	flag.PrintDefaults()
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

	//安静模式,只输出warn及以上的日志
	if quiet {
		// Only log the warning severity or above.
		logs.SetLevel(logs.WarnLevel)
	}
}
func InitEnv() {
	//线程数==CPU数量
	runtime.GOMAXPROCS(runtime.NumCPU())
}
func init() {
	//初始化命令行参数
	initFlag()
}
func main() {

	flag.Parse()
	//版本信息
	if version {
		showVersion()
		return
	}
	//帮助文档
	if help {
		flag.Usage()
		return
	}

	var (
		err error
	)

	//初始化线程
	InitEnv()

	//加载配置
	err = lib.InitConfig(confFile)
	if err != nil {
		goto ERR
	}
	initLogs(lib.G_config.LogsProduction)
	logs.Info("Crontab Master Running...")
	logs.Info("读取配置文件[成功]")

	//启动任务管理器
	err = lib.InitJobMgr()
	if err != nil {
		goto ERR
	}
	if nil == lib.G_jobMgr {
		err = errors.New(" 连接 ETCD 数据库出错,初始化 LogMgr实例 [失败]")
		goto ERR
	}
	logs.Info("启动任务管理器[成功]")
	//启动日志管理器
	err = lib.InitLogMgr()
	if err != nil {
		goto ERR
	}
	if nil == lib.G_jobMgr {
		err = errors.New(" 连接  mongodb 数据库出错,初始化 G_logMgr 实例 [失败]")
		goto ERR
	}
	logs.Info("启动日志管理器[成功]")
	//启动服务发现
	err = lib.InitWorkerMgr()
	if err != nil {
		goto ERR
	}
	if nil == lib.G_workerMgr {
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
	logs.Error("Master启动失败:" + err.Error())
}

func showVersion() {
	fmt.Println(fmt.Sprintf("v%s", version_str))
}
