package worker

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"sync"
	"time"
)

//日志接口类
type Log interface {
	SaveLogs(batch *common.LogBatch)
}

//mongodb的日志模型
type MongoDbLog struct {
	client        *mongo.Client
	logCollection *mongo.Collection
}

func (mongodb *MongoDbLog) SaveLogs(batch *common.LogBatch) {
	fmt.Println("MongoDbLog批量写入日志")
	//不处理是否保存成功
	_, err := mongodb.logCollection.InsertMany(context.TODO(), batch.Logs)
	if err != nil {
		fmt.Println("写入日志出错了", err)
	}
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
