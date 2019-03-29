package worker

import "testing"

func TestCommand(t *testing.T) {
	invoker := NewInvoker()
	concomA := NewConcreteCommandA()
	receA := NewReceiverA()

	concomA.SetReceiver(*receA)
	invoker.SetCommand(concomA)


	invoker.ExecuteCommand(nil)
}
