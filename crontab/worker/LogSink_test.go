package worker

import (
	"fmt"
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"strconv"
	"sync"
	"testing"
)

var wg = &sync.WaitGroup{}

type TestWriter struct {
}

func (w *TestWriter) Write(p []byte) (n int, err error) {
	wg.Done()
	fmt.Println(len(p))
	return 0, nil
}

func TestLogSink1(t *testing.T) {

	err := InitConfig("./main/worker.json")
	if err != nil {
		t.Fatal("Config error", err)
	}
	G_config.JobLogCommitTimeout = 10000 //日志自动提交超时
	G_config.JobLogBatchSize = 10        //日志批次大小

	fmt.Println("JobLogCommitTimeout=", G_config.JobLogCommitTimeout)
	fmt.Println("JobLogBatchSize=", G_config.JobLogBatchSize)
	fmt.Println("MAIN=", &wg)
	w := &TestWriter{}

	err = InitLogSink(w)
	if err != nil {
		t.Fatal("InitLogSink ERROR", err)
	}
	//注意整除的影响
	FOR_SIZE := 100
	bsize := FOR_SIZE / G_config.JobLogBatchSize
	for i := 0; i < FOR_SIZE; i++ {

		if i%bsize == 0 { //
			wg.Add(1)
		}

		G_logSink.Append(&common.JobLog{
			JobName: "JobName is" + strconv.Itoa(i),
			Err:     "This is Error " + strconv.Itoa(i),
		})

	}
	wg.Wait()

	if err != nil {
		t.Fatal("ERROR", err)
	}

}
