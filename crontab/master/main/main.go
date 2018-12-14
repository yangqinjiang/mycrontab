package main

import (
	"runtime"
	"github.com/yangqinjiang/mycrontab/crontab/master"
	"fmt"
)

func InitEnv() {
	//线程数==CPU数量
	runtime.GOMAXPROCS(runtime.NumCPU())
}
func main() {
	var (
		err error
	)
	//初始化线程
	InitEnv()

	//启动Api Http服务
	err = master.InitApiServer()
	if err != nil{
		goto ERR//启动出错,直接跳出
	}


	ERR:
		fmt.Println(err)
}
