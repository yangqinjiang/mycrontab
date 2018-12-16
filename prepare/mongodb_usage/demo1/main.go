package main

import (
	"context"
	"fmt"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
	"github.com/mongodb/mongo-go-driver/mongo"
	"github.com/mongodb/mongo-go-driver/mongo/clientopt"
	"github.com/mongodb/mongo-go-driver/mongo/findopt"
	"time"
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

//查询
type FindByJobName struct {
	JobName string `bson:"jobName"`
}

//startTime 小于某时间
//{"$lt":timestamp}
type TimeBeforeCond struct {
	Before int64 `bson:"$lt"`
}

//{"timePoint.startTime":{"$lt":timestamp}}
type DeleteCond struct {
	BeforeCond TimeBeforeCond `bson:"timePoint.startTime"`
}

func main() {
	fmt.Println("使用mongodb")
	//1,建立连接

	client, err := mongo.Connect(context.TODO(), "mongodb://106.12.37.124:27017",
		clientopt.ConnectTimeout(5*time.Second),
		clientopt.Auth(clientopt.Credential{
			Username: "root",
			Password: "123456",
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
	if err != nil {
		fmt.Println("InsertOne error=", err)
		return
	}
	//id 默认生成一个全局唯一ID,objectId, 12字节的二进制

	docid := ok.InsertedID.(objectid.ObjectID)
	//打印成十六进制
	fmt.Println("插入成功:", docid.Hex())

	//查找
	//按照FindByJobName字段过滤,找出JobName=10的记录
	cond := &FindByJobName{
		JobName: "job10", //{"jobName":"job10"}
	}
	//查询(过滤+分页)
	cursor, err := collection.Find(context.TODO(), cond, findopt.Skip(0), findopt.Limit(2))
	if err != nil {
		fmt.Println(err)
		return
	}
	//延迟释放游标
	defer cursor.Close(context.TODO())

	//遍历结果集
	for cursor.Next(context.TODO()) {
		//定义一个日志对象
		record := &LogRecord{}

		//反序列化bson到对象
		err := cursor.Decode(record)

		if err != nil {
			fmt.Println(err)
			return
		}
		//把日志行打印出来
		fmt.Println(*record)

	}

	//删除
	//删除开始时间早于当前时间的所有日志($lt是less than)
	//delete({"timePoint.startTime":{"$lt":当前时间-30s}})

	delCond := &DeleteCond{BeforeCond: TimeBeforeCond{Before: time.Now().Unix() - 30}}

	delResp, err := collection.DeleteMany(context.TODO(), delCond)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("删除的行数:", delResp.DeletedCount)

}
