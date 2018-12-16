package main

import (
	"fmt"
	"github.com/gorhill/cronexpr"
	"time"
)

func main() {
	//每5s执行一次
	expr, err := cronexpr.Parse("*/5 * * * * * *")
	if err != nil {
		fmt.Println(err)
	}
	//  /5
	// 0 5 10 15
	//当前时间
	now := time.Now()
	//下次调度时间
	nextTime := expr.Next(now)
	fmt.Println("下次调度时间:", nextTime)

	//等待这个定时器超时
	time.AfterFunc(nextTime.Sub(now), func() {
		fmt.Println("被调度了", time.Now())
	})
	select {
	case <-time.NewTimer(5 * time.Second).C:
		fmt.Println("定时器到期")
	}
	//time.Sleep(5 * time.Second)
	fmt.Println("主协程退出", time.Now())

}
