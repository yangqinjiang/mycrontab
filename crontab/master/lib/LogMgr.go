package lib

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/yangqinjiang/mycrontab/master/common"
	"sync"
	"time"
)

type LogMgr struct {
	client        *mongo.Client
	logCollection *mongo.Collection
}

var (
	//单例
	G_logMgr   *LogMgr
	onceLogMgr sync.Once
)

//初始化mongodb的实例
func InitLogMgr() (err error) {
	onceLogMgr.Do(func() {

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
			return
		}
		//选择db和collection
		G_logMgr = &LogMgr{
			client:        client,
			logCollection: client.Database("cron").Collection("log"),
		}
	})
	return
}

//查询日志
func (logMgr *LogMgr) ListLog(name string, skip int, limit int) (logArr []*common.JobLog, err error) {
	var (
		filter *common.JobLogFilter
		jobLog *common.JobLog
	)
	//初始化
	logArr = make([]*common.JobLog, 0)
	//过滤条件
	filter = &common.JobLogFilter{JobName: name}
	//按照任务时间排序
	//logSort := &common.SortLogByStartTime{SortOrder: -1}
	cursor, err := logMgr.logCollection.Find(context.TODO(), filter)
	if err != nil {
		return
	}
	//延迟释放
	defer cursor.Close(context.TODO())
	for cursor.Next(context.TODO()) {
		jobLog = &common.JobLog{}

		//反序列化bson
		if err = cursor.Decode(jobLog); err != nil {
			continue //有日志不合法,
		}
		logArr = append(logArr, jobLog)
	}
	return
}
