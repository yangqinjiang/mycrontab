package log

import (
	"context"
	logs "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	// "go.mongodb.org/mongo-driver/bson"
	"sync"
	"time"
	"github.com/yangqinjiang/mycrontab/worker/common"
	"github.com/yangqinjiang/mycrontab/worker/lib/config"
)

//mongodb的日志模型
//实现 JobLoger 接口方法
type MongoDbLog struct {
	client        *mongo.Client
	logCollection *mongo.Collection
}

//MongoDbLog批量写入日志
func (mongodb *MongoDbLog) Write(jobLog *common.LogBatch) (n int, err error) {
	logs.Info("MongoDbLog批量写入日志,len=",len(jobLog.Logs) )
	docs := make([]interface{}, len(jobLog.Logs))
	for key,value := range jobLog.Logs {
		if nil != value{
			logs.Debug(value)
			docs[key] = value
		}else{
			logs.Error("MongoDbLog write one log of jobLog.Logs  is nil")
		}
	}
 
	logs.Info("写入日志到mongoDb,docs数量=",len(docs))
	_, err = mongodb.logCollection.InsertMany(context.TODO(),docs)

	if err != nil {
		logs.Error("日志写入到mongoDb,失败", err)
		return  0,err
	}
	logs.Info("日志写入到mongoDb, 成功")
	return len(docs),nil
}

var (
	//单例
	G_MongoDbLog *MongoDbLog
	oncelog1     sync.Once
)

//mongodb的日志模型
func InitMongoDbLog() (err error) {
	oncelog1.Do(func() {
		var (
			client *mongo.Client
		)


		//建立mongodb链接
		if client, err = mongo.Connect(context.TODO(),
			//连接超时
			options.Client().SetConnectTimeout(time.Duration(config.G_config.MongodbConnectTimeout)*time.Millisecond),
			//连接URL
			options.Client().ApplyURI(config.G_config.MongodbUri),
			//连接认证的用户信息
			options.Client().SetAuth(options.Credential{
				Username:config.G_config.MongodbUsername,
				Password:config.G_config.MongodbPassword})); err != nil {
					logs.Error("连接mongoDb失败")
			return
		}
		logs.Info("连接mongoDb成功")

		//选择db和collection
		G_MongoDbLog = &MongoDbLog{
			client:        client,
			logCollection: client.Database("cron").Collection("log"),
		}
	})
	return
}
