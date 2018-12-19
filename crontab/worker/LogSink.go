package worker

import (
	"context"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"time"
	"fmt"
)

type LogSink struct {
	client         *mongo.Client
	logCollection  *mongo.Collection
	logChan        chan *common.JobLog
	autoCommitChan chan *common.LogBatch
}

var (
	//单例
	G_logSink *LogSink
)

//初始化mongodb的实例
func InitLogSink() (err error) {
	var (
		client *mongo.Client
	)
	//建立mongodb链接
	if client, err = mongo.Connect(context.TODO(), G_config.MongodbUri,
		clientopt.ConnectTimeout(time.Duration(G_config.MongodbConnectTimeout)*time.Millisecond),
		clientopt.Auth(clientopt.Credential{
			Username: G_config.MongodbUsername,
			Password: G_config.MongodbPassword,
		})); err != nil {
		return
	}
	//选择db和collection
	G_logSink = &LogSink{
		client:         client,
		logCollection:  client.Database("cron").Collection("log"),
		logChan:        make(chan *common.JobLog, 1000),   //日志队列
		autoCommitChan: make(chan *common.LogBatch, 1000), //提交日志的信息
	}
	//启动一个mongodb处理协程
	go G_logSink.writeLoop()
	return
}

//批量写入日志
func (logSink *LogSink) saveLogs(batch *common.LogBatch) {
	fmt.Println("批量写入日志")
	//不处理是否保存成功
	_,err := logSink.logCollection.InsertMany(context.TODO(), batch.Logs)
	if err != nil{
		fmt.Println("写入日志出错了",err)
	}
}

//日志存储协程
func (logSink *LogSink) writeLoop() {
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

				//闭包的作用,防止logBatch被修改后,影响到chan
				commitTimer = time.AfterFunc(time.Duration(G_config.JobLogCommitTimeout)*time.Millisecond,
					func(batch *common.LogBatch) func() {
						return func() {
							fmt.Println("让这个批次超时自动提交")
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
				fmt.Println("如果批次满了,就发送")
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

//发送日志
func (logSink *LogSink) Append(jobLog *common.JobLog) {
	select {
	case logSink.logChan <- jobLog: //未满
	default:
		//队列满了就丢弃
	}

}
