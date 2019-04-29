package worker

import "fmt"
//打印日志到控制台
type ConsoleLog struct {
	Log
}

func (w *ConsoleLog) Write(p []byte) (n int, err error) {

	fmt.Println("ConsoleLog wirte ,len=",len(p)) //只打印 p的长度
	return 0, nil
}
