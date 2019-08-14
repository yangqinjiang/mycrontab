package lib

import (
	"testing"
	"github.com/yangqinjiang/mycrontab/worker/lib"
)

func TestCommand(t *testing.T) {

	// 命令接收者
	receA := lib.NewCmdReceiver()

	//命令对象
	concomA := lib.NewConcreteCommand(*receA)

	invoker := lib.NewInvoker()
	//加载命令给调用者
	invoker.SetCommand(concomA)


	//调用者 执行 命令对象的execute函数
	invoker.ExecuteCommand(nil)

	cmd := lib.CommandFactory("sh")
	if nil == cmd{
		t.Fatal("应该存在 sh 命令对象")
	}
	cmd = lib.CommandFactory("")
	if nil == cmd{
		t.Fatal("应该存在 sh 命令对象")
	}

	cmd = lib.CommandFactory("not_bin")
	if nil == cmd{
		t.Fatal("应该存在 not_bin 命令对象")
	}

	cmd = lib.CommandFactory("error_cmd")
	if nil != cmd{
		t.Fatal("不应该存在 error_cmd 命令对象")
	}
}
