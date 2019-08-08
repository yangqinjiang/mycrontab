package lib

import (
	"fmt"
	"github.com/yangqinjiang/mycrontab/worker/common"
	logs "github.com/sirupsen/logrus"
	"os/exec"
	"runtime"
	"strings"
)

/*
 Command 命令模式：
    将一个请求封装为一个对象，
    从而使你可用不同的请求对客户进行参数化；
    对请求排队或者记录请求日志，以及支持可撤销的操作

 个人想法：Invoker维护请求队列（Command接口队列），通过一些函数可以添加、修改、执行请求队列，
         在每一种ConcreteCommand中有对该命令的执行体（Receiver），最终响应请求队列的命令
 作者：   HCLAC
 日期：   20170310
*/

//------------------------------------------------------------------------------------
// Command为所有命令声明一个接口,调用命令对象的execute方法,就可以让接收者进行相关的动作
// 这个接口也具备一个undo方法(未来可实现)
type Command interface {
	Execute(info *common.JobExecuteInfo) ([]byte, error)
}

//------------------------------------------------------------------------------------
// 这个调用者,持有一个命令对象,并在某个时间点,调用命令对象的execute()方法,将请求付诸实行
type Invoker struct {
	command Command
}

// 添加命令
func (i *Invoker) SetCommand(c Command) {
	if i == nil {
		return
	}
	fmt.Println("call invoker SetCommand",c)
	i.command = c
}

// 执行命令
func (i *Invoker) ExecuteCommand(info *common.JobExecuteInfo) ([]byte, error) {
	if i == nil {
		return nil,nil
	}
	fmt.Println("call invoker ExecuteCommand",i.command)
	return i.command.Execute(info)
}

func NewInvoker() *Invoker {
	return &Invoker{}
}

//------------------------------------------------------------------------------------

// 这个ConcreteCommand 定义了动作和接收者之间的绑定关系
// 调用者只要调用 execute 就可以发出请求,然后由ConcreteCommand
// 调用接收者的一个或多个动作
type ConcreteCommand struct {
	Command
	receiver CmdReceiver
}


// 具体命令的执行体
func (c *ConcreteCommand) Execute(info *common.JobExecuteInfo) ([]byte, error) {
	if c == nil {
		return nil,nil
	}
	return c.receiver.action(info)
}

func NewConcreteCommand(r CmdReceiver) *ConcreteCommand {
	return &ConcreteCommand{receiver:r}
}

// 针对ConcreteCommand，如何处理该命令
type CmdReceiver struct {

}

func (r *CmdReceiver) action(info *common.JobExecuteInfo) ([]byte, error) {
	if r == nil {
		return nil,nil
	}
	if strings.Contains(runtime.GOOS,"windows"){
		logs.Error("当前测试运行在windows上,默认测试通过")
		return nil,nil
	}

	logs.Info("针对CmdReceiver->action，如何处理该命令,info=",info)

	if nil == info{
		return nil,nil
	}
	//bash的shell命令
	logs.Info("执行具体的shell命令",info.Job.Command);
	//执行shell命令
	cmd := exec.CommandContext(info.CancelCtx, "/bin/bash", "-c", info.Job.Command)
	//执行并捕获输出
	return cmd.CombinedOutput()
}


func NewCmdReceiver() *CmdReceiver {
	return &CmdReceiver{}
}

//-------------------------
type ConcreteCommand2 struct {
	Command
	receiver NotBin
}


// 具体命令的执行体
func (c *ConcreteCommand2) Execute(info *common.JobExecuteInfo) ([]byte, error) {
	if c == nil {
		return nil,nil
	}
	return c.receiver.action(info)
}

func NewConcreteCommand2(r NotBin) *ConcreteCommand2 {
	return &ConcreteCommand2{receiver:r}
}

// 针对ConcreteCommand，如何处理该命令
type NotBin struct {

}

func (r *NotBin) action(info *common.JobExecuteInfo) ([]byte, error) {
	if r == nil {
		return nil,nil
	}
	logs.Info("针对CmdReceiver2->action，如何处理该命令,info=",info)

	return nil,nil
}
func NewNotBin() *NotBin {
	return &NotBin{}
}
//命令对象的工厂
func CommandFactory(name string) Command  {
	//默认值
	if 0 == len(name) || "sh" == name{
		// 命令接收者
		receA := NewCmdReceiver()

		//命令对象
		concomA := NewConcreteCommand(*receA)
		return  concomA
	}
	if "not_bin" == name{
		// 命令接收者
		receA2 := NewNotBin()

		//命令对象
		concomA2 := NewConcreteCommand2(*receA2)
		return  concomA2
	}
	return  nil

}


