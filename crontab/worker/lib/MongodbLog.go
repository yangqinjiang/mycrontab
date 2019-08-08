package lib

import (
	"context"
	logs "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"sync"
	"time"

	"github.com/yangqinjiang/mycrontab/worker/common"
)

//mongodb的日志模型
type MongoDbLog struct {
	JobLoger //实现Log接口方法
	client        *mongo.Client
	logCollection *mongo.Collection
}

//MongoDbLog批量写入日志
func (mongodb *MongoDbLog) Write(jobLog *common.LogBatch) (n int, err error) {
	logs.Info("MongoDbLog批量写入日志" )
	doc := make([]interface{}, len(jobLog.Logs))
	for key,value := range jobLog.Logs {
		doc[key] = value
	}

	logs.Debug("写入日志到mongoDb,日志数量len=",len(doc))
	_, err = mongodb.logCollection.InsertMany(context.TODO(),doc)

	if err != nil {
		logs.Error("日志写入到mongoDb,失败", err)
		return  0,err
	}
	logs.Debug("日志写入到mongoDb, 成功")
	return len(doc),nil
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
			options.Client().SetConnectTimeout(time.Duration(G_config.MongodbConnectTimeout)*time.Millisecond),
			//连接URL
			options.Client().ApplyURI(G_config.MongodbUri),
			//连接认证的用户信息
			options.Client().SetAuth(options.Credential{
				Username:G_config.MongodbUsername,
				Password:G_config.MongodbPassword})); err != nil {
					logs.Error("连接mongoDb失败")
			return
		}
		logs.Debug("连接mongoDb成功")

		//选择db和collection
		G_MongoDbLog = &MongoDbLog{
			client:        client,
			logCollection: client.Database("cron").Collection("log"),
		}
	})
	return
}
