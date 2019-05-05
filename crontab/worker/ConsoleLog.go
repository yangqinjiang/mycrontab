package worker

import (
	"fmt"
	"github.com/yangqinjiang/mycrontab/crontab/common"
)
//打印日志到控制台
type ConsoleLog struct {
	JobLoger
}

func (w *ConsoleLog) Write(jobLog *common.JobLog) (n int, err error) {

	fmt.Println("ConsoleLog wirte =",jobLog) //只打印 p的长度
	return 0, nil
}
