package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"time"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

//任务的执行时间点
type TimePoint struct {
	StartTime int64 `bson:"startTime"`
	EndTime   int64 `bson:"endTime"`
}

//一条日志
type LogRecord struct {
	JobName   string    `bson:"jobName"`   //任务名
	Command   string    `bson:"command"`   //shell命令
	Err       string    `bson:"err"`       //脚本错误
	Content   string    `bson:"content"`   //脚本输出
	TimePoint TimePoint `bson:"timePoint"` //
}

func main() {
	fmt.Println("使用mongodb")
	//1,建立连接

	client, err := mongo.Connect(context.TODO(), "mongodb://106.12.37.124:27017",
		clientopt.ConnectTimeout(5*time.Second),
		clientopt.Auth(clientopt.Credential{
			Username:"root",
			Password:"123456",
		}))
	if err != nil {
		fmt.Println(err)
		return
	}
	//2选择数据库my_db
	database := client.Database("my_db")

	//3选择表my_collection
	collection := database.Collection("my_collection")
	fmt.Println("Name:", collection.Name())

	collection = collection

	//插入一条记录(bson)
	recode := &LogRecord{
		JobName: "job10",
		Command: "echo hello",
		Err:     "",
		Content: "Hello",
		TimePoint: TimePoint{StartTime: time.Now().Unix(),
			EndTime: time.Now().Unix() + 10},
	}

	ok, err := collection.InsertOne(context.TODO(), recode)
	if err != nil{
		fmt.Println("InsertOne error=",err)
		return
	}
	//id 默认生成一个全局唯一ID,objectId, 12字节的二进制

	docid:=ok.InsertedID.(objectid.ObjectID)
	//打印成十六进制
	fmt.Println("插入成功:",docid.Hex())

}
