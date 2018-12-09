package main

import (
	"os/exec"
	"fmt"
)

func main() {

	//生成cmd
	cmd := exec.Command("/bin/sh","-c","sleep 3;echo helloworld")
	//cmd := exec.Command("ping","127.0.0.1")
	//cmd.Stdout = os.Stdout
	//cmd.Run()
	//

	//执行命令,捕获子进程的输出(pipe)
	output,err := cmd.CombinedOutput();
	if err != nil{
		fmt.Println("ERROR=",err)
		return
	}

	//正常运行,打印子进程的输出
	fmt.Println("output:",string(output))
}
