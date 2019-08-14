package lib

import (
	"github.com/yangqinjiang/mycrontab/worker/common"
	"github.com/yangqinjiang/mycrontab/worker/lib"
	"strconv"
	"testing"
	"time"
)
var config_file_path = "../config/worker.json"


func TestLogSinkOnlyPrint(t *testing.T) {

	err := lib.InitConfig(config_file_path)
	if err != nil {
		t.Fatal("Config error", err)
	}
	lib.G_config.JobLogCommitTimeout = 10000 //日志自动提交超时
	lib.G_config.JobLogBatchSize = 10        //日志批次大小

	w := &lib.TestWriter{}

	err = lib.InitJobLogMemoryBuffer(w)
	if err != nil {
		t.Fatal("InitJobLogMemoryBuffer ERROR", err)
	}
	//避免单例模式的影响
	//G_jobLogMemoryBuffer.JobLoger = w
	//注意整除的影响
	FOR_SIZE := 1002

	dd := make([]*common.JobLog,FOR_SIZE)

	for i := 0; i < FOR_SIZE; i++ {

		d := &common.JobLog{
			JobName: "JobName is " + strconv.Itoa(i),
			Err:     "This is Error " + strconv.Itoa(i),
		}
		dd[i] = d
	}
	lib.G_jobLogMemoryBuffer.Write(&common.LogBatch{dd})
	if lib.G_config.JobLogBatchSize > 0{
		time.Sleep(time.Duration(lib.G_config.JobLogCommitTimeout+2000)*time.Millisecond)

	}

	if err != nil {
		t.Fatal("ERROR", err)
	}
	if lib.G_config.JobLogBatchSize == 0 && lib.G_jobLogMemoryBuffer.LogChanLength() != FOR_SIZE {
		t.Fatal("ERROR, 写入内存的日志size不正确", err)
	}
	t.Log("OK")

}
func TestLogSinkOnlyPrintWithTimeout(t *testing.T) {

	err := lib.InitConfig(config_file_path)
	if err != nil {
		t.Fatal("Config error", err)
	}
	lib.G_config.JobLogCommitTimeout = 1 //日志自动提交超时
	lib.G_config.JobLogBatchSize = 100        //日志批次大小

	w := &lib.TestWriter{}

	err = lib.InitJobLogMemoryBuffer(w)
	if err != nil {
		t.Fatal("InitJobLogMemoryBuffer ERROR", err)
	}
	//避免单例模式的影响
	lib.G_jobLogMemoryBuffer.JobLoger = w
	//注意整除的影响
	FOR_SIZE := 5

	dd := make([]*common.JobLog,FOR_SIZE)
	for i := 0; i < FOR_SIZE; i++ {
		//time.Sleep(2*time.Second)
		d:= &common.JobLog{
			JobName: "JobName is " + strconv.Itoa(i),
			Err:     "This is Error " + strconv.Itoa(i),
		}
		dd[i] = d
	}
	lib.G_jobLogMemoryBuffer.Write(&common.LogBatch{dd})

	time.Sleep(time.Duration(lib.G_config.JobLogCommitTimeout+2000)*time.Millisecond)

	if err != nil {
		t.Fatal("ERROR", err)
	}

}

//测试mongodb的保存日志
func TestLogSinkToMongoDb(t *testing.T) {
	err := lib.InitConfig(config_file_path)
	if err != nil {
		t.Fatal("Config error", err)
	}
	lib.G_config.JobLogCommitTimeout = 10000 //日志自动提交超时
	lib.G_config.JobLogBatchSize = 100        //日志批次大小

	//第一次实例化
	err = lib.InitMongoDbLog()
	if err != nil {
		t.Fatal("MongoDbLog error", err)
	}
	backG_MongoDbLog := lib.G_MongoDbLog
	//第二次实例化
	err = lib.InitMongoDbLog()
	if err != nil {
		t.Fatal("MongoDbLog error", err)
	}
	if backG_MongoDbLog != lib.G_MongoDbLog{
		t.Fatal("InitMongoDbLog error,单例模式错误")
	}

	err = lib.InitJobLogMemoryBuffer(lib.G_MongoDbLog)
	if err != nil {
		t.Fatal("InitJobLogMemoryBuffer ERROR", err)
	}
	//避免单例模式的影响
	lib.G_jobLogMemoryBuffer.JobLoger = lib.G_MongoDbLog
	//注意整除的影响
	FOR_SIZE := 10
	dd := make([]*common.JobLog,FOR_SIZE)
	for i := 0; i < FOR_SIZE; i++ {
		d := &common.JobLog{
			JobName: "JobName : " + strconv.Itoa(i),
			Err:     "This is Error :" + strconv.Itoa(i),
		}
		dd[i] = d
	}
	lib.G_jobLogMemoryBuffer.Write(&common.LogBatch{dd})
	time.Sleep(time.Duration(lib.G_config.JobLogCommitTimeout+2000)*time.Millisecond)
	if err != nil {
		t.Fatal("ERROR", err)
	}

}

func TestLogSinkNoWriter(t *testing.T) {

	err := lib.InitConfig(config_file_path)
	if err != nil {
		t.Fatal("Config error", err)
	}
	lib.G_config.JobLogCommitTimeout = 10000 //日志自动提交超时
	lib.G_config.JobLogBatchSize = 10        //日志批次大小

	err = lib.InitJobLogMemoryBuffer(nil)
	if err == nil{
		t.Fatal("错误,必须传入common.Log的实现类")
	}

}

