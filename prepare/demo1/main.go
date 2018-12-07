package main

import (
	"os/exec"
	"fmt"
)

func main() {

	//生成cmd
	cmd := exec.Command("C:\\cygwin64\\bin\\bash.exe","-c","echo helloworld")

	//执行命令,捕获子进程的输出(pipe)
	output,err := cmd.CombinedOutput();
	if err != nil{
		fmt.Println("ERROR=",err)
		return
	}

	//正常运行,打印子进程的输出
	fmt.Println("output:",string(output))
}
