package worker

import (
	"testing"
)

func TestCommand(t *testing.T) {

	// 命令接收者
	receA := NewCmdReceiver()

	//命令对象
	concomA := NewConcreteCommand(*receA)

	invoker := NewInvoker()
	//加载命令给调用者
	invoker.SetCommand(concomA)


	//调用者 执行 命令对象的execute函数
	invoker.ExecuteCommand(nil)

	cmd := CommandFactory("sh")
	if nil == cmd{
		t.Fatal("应该存在 sh 命令对象")
	}
	cmd = CommandFactory("")
	if nil == cmd{
		t.Fatal("应该存在 sh 命令对象")
	}

	cmd = CommandFactory("not_bin")
	if nil == cmd{
		t.Fatal("应该存在 not_bin 命令对象")
	}

	cmd = CommandFactory("error_cmd")
	if nil != cmd{
		t.Fatal("不应该存在 error_cmd 命令对象")
	}
}
