package lib
import (
	logs "github.com/sirupsen/logrus"
	"github.com/yangqinjiang/mycrontab/worker/common"
)
type TestWriter struct {
}

func (w *TestWriter) Write(jobLog *common.LogBatch) (n int, err error) {

	logs.Info("call TestWriter ,print =>") //只打印 p的长度
	for _,log := range jobLog.Logs{
		if nil != log{
			logs.Info("one log.Name=",log.JobName)
		}

	}
	return 0, nil
}