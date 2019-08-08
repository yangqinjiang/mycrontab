package lib

import (
	"fmt"
	"github.com/yangqinjiang/mycrontab/worker/common"
)
//打印日志到控制台
type ConsoleLog struct {
	JobLoger
}

func (w *ConsoleLog) Write(jobLog *common.LogBatch) (n int, err error) {

	fmt.Println("ConsoleLog: wirte =",jobLog) //只打印 p的长度
	return 0, nil
}
