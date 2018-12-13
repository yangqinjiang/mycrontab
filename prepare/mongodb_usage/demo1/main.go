package main

import (
	"fmt"
	"github.com/mongodb/mongo-go-driver/mongo"
	"context"
	"github.com/mongodb/mongo-go-driver/mongo/options"
	"time"
)

func main() {
	fmt.Println("使用mongodb")
	clientopt :=options.Client();
	//1,建立连接
	clientopt.SetConnectTimeout(5*time.Second).SetAuth(options.Credential{Username:"root",Password:"123456"})
	client,err := mongo.Connect(context.TODO(),"mongodb://106.12.37.124:27017",clientopt)
	if err != nil{
		fmt.Println(err)
		return
	}
	fmt.Println("链接成功mongodb")
	//2选择数据库my_db
	database := client.Database("my_db")

	//3选择表my_collection
	collection := database.Collection("my_my_collection")
	fmt.Println("Name:",collection.Name())

	collection = collection

}
