package worker

import (
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/astaxie/beego/logs"
	"context"
	"sync"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
)

//mongodb的日志模型
type MongoDbLog struct {
	client        *mongo.Client
	logCollection *mongo.Collection
}


func (mongodb *MongoDbLog) Write(p []byte) (n int, err error) {
	logs.Info("MongoDbLog批量写入日志",len(p))

	documents := make([]interface{}, len(p))
	for i, s := range p {
		documents[i] = s
	}
	_, err = mongodb.logCollection.InsertMany(context.TODO(),documents)
	n = 0
	if err != nil {
		logs.Error("写入日志出错了", err)
	}
	return
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
		if client, err = mongo.Connect(context.TODO(), G_config.MongodbUri,
			clientopt.ConnectTimeout(time.Duration(G_config.MongodbConnectTimeout)*time.Millisecond),
			clientopt.Auth(clientopt.Credential{
				Username: G_config.MongodbUsername,
				Password: G_config.MongodbPassword,
			})); err != nil {
			return
		}

		//选择db和collection
		G_MongoDbLog = &MongoDbLog{
			client:        client,
			logCollection: client.Database("cron").Collection("log"),
		}
	})
	return
}
