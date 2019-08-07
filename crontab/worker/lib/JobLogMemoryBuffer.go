package lib

import (
	"errors"
	logs "github.com/sirupsen/logrus"
	"github.com/yangqinjiang/mycrontab/worker/common"
	"sync"
	"time"
)

//任务的执行日志 缓冲器
type JobLogMemoryBuffer struct {
	JobLoger                             //日志记录接口
	logChan        chan *common.JobLog   //日志队列
	autoCommitChan chan *common.LogBatch //提交日志的信息
}

var (
	//单例
	G_jobLogMemoryBuffer *JobLogMemoryBuffer
	oncelog              sync.Once
)

//初始化日志缓存器
func InitJobLogMemoryBuffer(jobLoger JobLoger) (err error) {
	if nil == jobLoger {
		return errors.New("必须传入common.Log的实现类")
	}
	oncelog.Do(func() {

		//选择db和collection
		G_jobLogMemoryBuffer = &JobLogMemoryBuffer{
			logChan:        make(chan *common.JobLog, 10000),   //日志队列
			autoCommitChan: make(chan *common.LogBatch, 1000), //提交日志的信息
			JobLoger:       jobLoger,
		}

		//批处理容量必须大于 0
		if G_config.JobLogBatchSize > 0 {
			logs.Info("启动一个日志处理协程...")
			//处理协程
			go G_jobLogMemoryBuffer.writeLoop()
		}

	})

	return
}

//返回内存中,日志的长度
func (logSink *JobLogMemoryBuffer) LogChanLength() int {
	return len(G_jobLogMemoryBuffer.logChan)
}

//发送日志
func (logSink *JobLogMemoryBuffer) Write(jobLog *common.LogBatch) (n int, err error) {
	logs.Debug("Batch write jobLog to logChan, len=", len(jobLog.Logs))
	for _, log := range jobLog.Logs {
		select {
		case logSink.logChan <- log: //未满
		default:
			//队列满了就丢弃
		}
	}

	return 0, nil
}

//---------------------private func------------
//批量写入日志
func (logSink *JobLogMemoryBuffer) saveLogs(batch *common.LogBatch) {
	logs.Info("LogSink批量写入日志 数量=", len(batch.Logs))

	if nil == logSink.JobLoger {
		logs.Error("logSink.JobLoger is nil")
		return
	}

	logs.Info("批量写入Logs----->")
	//批量写入
	_, err := logSink.JobLoger.Write(batch)
	logs.Info("写入完成<------")
	if err != nil {
		logs.Error("logSink.JobLoger Write some log Error", err)
		return
	}
}

//日志存储协程
func (logSink *JobLogMemoryBuffer) writeLoop() {
	var (
		log          *common.JobLog
		logBatch     *common.LogBatch //当前的批次
		commitTimer  *time.Timer
		timeoutBatch *common.LogBatch //超时批次
	)
	for {
		select {
		case log = <-logSink.logChan:

			//把这条log写到mongodb中
			//logSink.logCollection.inserOne
			if logBatch == nil {
				logBatch = &common.LogBatch{}
				//让这个批次超时自动提交(给1s的时间)
				logs.Debug("start commitTimer")
				//闭包的作用,防止logBatch被修改后,影响到chan
				commitTimer = time.AfterFunc(time.Duration(G_config.JobLogCommitTimeout)*time.Millisecond,
					func(batch *common.LogBatch) func() {
						return func() {
							logs.Info("让这个批次超时自动提交")
							//发出超时通知,不能直接提交batch
							logSink.autoCommitChan <- batch
						}
					}(logBatch),
				)
			}
			//把新日志append到当前批次中
			logBatch.Logs = append(logBatch.Logs, log)
			//如果批次满了,就发送
			if len(logBatch.Logs) >= G_config.JobLogBatchSize {
				logs.Info("如果批次满了,就发送 len=",len(logBatch.Logs)," , JobLogBatchSize= ",G_config.JobLogBatchSize)
				//发送日志
				logSink.saveLogs(logBatch)
				//清空
				logBatch = nil
				//取消定时器
				commitTimer.Stop()
			}
		case timeoutBatch = <-logSink.autoCommitChan: //过期的批次

			//判断过期批次是否仍旧是当前的批次
			if timeoutBatch != logBatch {
				continue //跳过已经被提交的批次
			}
			//把批次写入到mongodb
			logSink.saveLogs(timeoutBatch)
			//清空
			logBatch = nil

		}
	}

}
