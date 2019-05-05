package worker

import (
	"fmt"
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"strconv"
	"testing"
	"time"
)

type TestWriter struct {
}

func (w *TestWriter) Write(jobLog *common.JobLog) (n int, err error) {

	fmt.Println("call TestWriter ,print =>",jobLog) //只打印 p的长度
	return 0, nil
}

func TestLogSinkOnlyPrint(t *testing.T) {

	err := InitConfig("./main/worker.json")
	if err != nil {
		t.Fatal("Config error", err)
	}
	G_config.JobLogCommitTimeout = 10000 //日志自动提交超时
	G_config.JobLogBatchSize = 10        //日志批次大小

	w := &TestWriter{}

	err = InitJobLogMemoryBuffer(w)
	if err != nil {
		t.Fatal("InitJobLogMemoryBuffer ERROR", err)
	}
	//避免单例模式的影响
	//G_jobLogMemoryBuffer.JobLoger = w
	//注意整除的影响
	FOR_SIZE := 100

	for i := 0; i < FOR_SIZE; i++ {

		G_jobLogMemoryBuffer.Write(&common.JobLog{
			JobName: "JobName is " + strconv.Itoa(i),
			Err:     "This is Error " + strconv.Itoa(i),
		})

	}
	if G_config.JobLogBatchSize > 0{
		time.Sleep(10*time.Second)
	}

	if err != nil {
		t.Fatal("ERROR", err)
	}
	if G_config.JobLogBatchSize == 0 && G_jobLogMemoryBuffer.LogChanLength() != FOR_SIZE {
		t.Fatal("ERROR, 写入内存的日志size不正确", err)
	}
	t.Log("OK")

}
func TestLogSinkOnlyPrintWithTimeout(t *testing.T) {

	err := InitConfig("./main/worker.json")
	if err != nil {
		t.Fatal("Config error", err)
	}
	G_config.JobLogCommitTimeout = 1 //日志自动提交超时
	G_config.JobLogBatchSize = 100        //日志批次大小

	w := &TestWriter{}

	err = InitJobLogMemoryBuffer(w)
	if err != nil {
		t.Fatal("InitJobLogMemoryBuffer ERROR", err)
	}
	//避免单例模式的影响
	G_jobLogMemoryBuffer.JobLoger = w
	//注意整除的影响
	FOR_SIZE := 5

	for i := 0; i < FOR_SIZE; i++ {
		time.Sleep(2*time.Second)
		G_jobLogMemoryBuffer.Write(&common.JobLog{
			JobName: "JobName is" + strconv.Itoa(i),
			Err:     "This is Error " + strconv.Itoa(i),
		})

	}
	time.Sleep(10*time.Second)
	if err != nil {
		t.Fatal("ERROR", err)
	}

}


func TestLogSinkToMongoDb(t *testing.T) {
	err := InitConfig("./main/worker.json")
	if err != nil {
		t.Fatal("Config error", err)
	}
	G_config.JobLogCommitTimeout = 10000 //日志自动提交超时
	G_config.JobLogBatchSize = 10        //日志批次大小

	//第一次实例化
	err = InitMongoDbLog()
	if err != nil {
		t.Fatal("MongoDbLog error", err)
	}
	backG_MongoDbLog := G_MongoDbLog
	//第二次实例化
	err = InitMongoDbLog()
	if err != nil {
		t.Fatal("MongoDbLog error", err)
	}
	if backG_MongoDbLog != G_MongoDbLog{
		t.Fatal("InitMongoDbLog error,单例模式错误")
	}

	err = InitJobLogMemoryBuffer(G_MongoDbLog)
	if err != nil {
		t.Fatal("InitJobLogMemoryBuffer ERROR", err)
	}
	//避免单例模式的影响
	G_jobLogMemoryBuffer.JobLoger = G_MongoDbLog
	//注意整除的影响
	FOR_SIZE := 100

	for i := 0; i < FOR_SIZE; i++ {

		G_jobLogMemoryBuffer.Write(&common.JobLog{
			JobName: "JobName is" + strconv.Itoa(i),
			Err:     "This is Error " + strconv.Itoa(i),
		})

	}
	time.Sleep(10*time.Second)
	if err != nil {
		t.Fatal("ERROR", err)
	}

}

func TestLogSinkNoWriter(t *testing.T) {

	err := InitConfig("./main/worker.json")
	if err != nil {
		t.Fatal("Config error", err)
	}
	G_config.JobLogCommitTimeout = 10000 //日志自动提交超时
	G_config.JobLogBatchSize = 10        //日志批次大小

	err = InitJobLogMemoryBuffer(nil)
	if err == nil{
		t.Fatal("错误,必须传入common.Log的实现类")
	}

}

