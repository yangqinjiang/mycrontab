package main

import (
	"context"
	"fmt"
	"os/exec"
	"time"
)

//命令运行结果
type result struct {
	output []byte
	err    error
}

func main() {

	//执行一个cmd,让它在一个协程里去执行,让它执行2秒,sleep 2;echo hello;
	//1秒的时候,我们杀死cmd
	var (
		ctx        context.Context
		cancelFunc context.CancelFunc
		resultChan chan *result
	)

	resultChan = make(chan *result, 1000)
	//生成cmd
	ctx, cancelFunc = context.WithCancel(context.TODO())
	go func() {
		cmd := exec.CommandContext(ctx, "/bin/sh", "-c", "sleep 3;echo helloworld;")
		//执行命令,捕获子进程的输出(pipe)
		output, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Println("ERROR=", err)
			return
		}

		//正常运行,打印子协程的输出
		fmt.Println("子协程 output:", string(output))
		resultChan <- &result{
			err:    err,
			output: output,
		}

	}()
	//继续往下走
	time.Sleep(1 * time.Second)
	//模拟取消任务的执行
	if false {
		cancelFunc()
	}

	//在main协程里,等待子协程的退出,并打印任务执行结果
	fmt.Println("等待中...")
	res := <-resultChan
	fmt.Println("res:", string(res.output))

}
