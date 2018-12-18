package worker

import (
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"context"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
)

type LogSink struct {
	client *mongo.Client
	logCollection *mongo.Collection
	logChan chan *common.JobLog
}
var (
	//单例
	G_logSink *LogSink
)

//初始化mongodb的实例
func InitLogSink()(err error)  {
	var (
		client *mongo.Client
	)
	//建立mongodb链接
	if client,err = mongo.Connect(context.TODO(),G_config.MongodbUri,
		clientopt.ConnectTimeout(time.Duration(G_config.MongodbConnectTimeout)*time.Microsecond),
		clientopt.Auth(clientopt.Credential{
			Username: G_config.MongodbUsername,
			Password: G_config.MongodbPassword,
		}));err != nil{
		return
	}
	//选择db和collection
	G_logSink = &LogSink{
		client:client,
		logCollection:client.Database("cron").Collection("log"),
		logChan:make(chan *common.JobLog,1000),//日志队列
	}
	//启动一个mongodb处理协程
	go G_logSink.writeLoop()
	return
}

//日志存储协程
func (logSink *LogSink)writeLoop()  {
	var (
		log *common.JobLog
	)
	for{
		select {
			case log = <- logSink.logChan:
				//TODO:
				//把这条log写到mongodb中
				//logSink.logCollection.inserOne
		}
	}

}