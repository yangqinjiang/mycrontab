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
}
