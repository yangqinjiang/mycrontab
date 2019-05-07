package worker

import (
	"bytes"
	"github.com/astaxie/beego/logs"
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"strconv"
	"testing"
	"time"
)

func TestJobPlanHeap(t *testing.T) {

	SIZE := 3
	j  := NewJobPlanMinHeap(SIZE)
	job_1 := &common.Job{Name: "job_1",CronExpr:"* * * * *"}
	j_1_0,_ := common.BuildJobSchedulePlan(job_1)
	j_1_1,_ := common.BuildJobSchedulePlan(job_1)
	//插入同一个job
	j.Insert(j_1_0)
	j.Insert(j_1_1)
	if 1 != j.Size(){
		t.Fatal("Insert 失败,数量 != 1")
	}
	//插入另外两个job
	job_2 := &common.Job{Name: "job_2",CronExpr:"* * * * *"}
	job_3 := &common.Job{Name: "job_3",CronExpr:"* * * * *"}
	jj_1,_ := common.BuildJobSchedulePlan(job_1)
	jj_2,_ := common.BuildJobSchedulePlan(job_2)
	jj_3,_ := common.BuildJobSchedulePlan(job_3)
	j.Insert(jj_1)
	j.Insert(jj_2)
	j.Insert(jj_3)
	if 3 != j.Size(){
		t.Fatal("Insert 失败,数量 != 3")
	}

	//删除
	j.Remove("job_2")
	j.Remove("job_2")
	j.Remove("job_2")
	if 2 != j.Size(){
		t.Fatal("Remove 失败,数量 != 2 实际是:",j.Size())
	}
}
//找出最早的
func TestExtractEarliestHeap(t *testing.T) {

	// 简单的性能测试,如下:
	// 1w, 耗时 1 ms~2 ms.
	// 10w, 耗时 10 ms~20 ms.
	// 100w, 耗时 100 ms以上.
	SIZE := 10
	j  := NewJobPlanMinHeap(SIZE)
	for i:=1; i<=SIZE;i++  {
		istr := strconv.Itoa(i)
		i_60 := i
		if i>= 60{ //大于 60,求余数
			i_60 = i%60
		}
		if i_60 <= 0{
			i_60  = 1
		}

		i_60_str := strconv.Itoa(i_60)
		one_job := &common.Job{Name: "job_"+istr,CronExpr:"*/"+i_60_str+" * * * * * *"}
		jj,err := common.BuildJobSchedulePlan(one_job)
		if err == nil{
			err = j.Insert(jj)
			if err != nil{
				t.Error(err.Error())
			}
		}else{
			t.Error(err.Error())
		}

	}
	if SIZE != j.Size(){
		t.Fatal("Insert 失败,数量 != ",SIZE)
	}
	go func() {
		for  {
			j.ExtractEarliest(func(jobPlan *common.JobSchedulePlan) (err error) {
				//t.Log("执行任务",jobPlan.Job.Name," ,cronExpr=",jobPlan.Job.CronExpr,"下次执行时间=",jobPlan.NextTime)
				//不存在,则构建一个
				jobExecuteInfo := common.BuildJobExecuteInfo(jobPlan)
				//任务执行的结果
				result := &common.JobExecuteResult{
					ExecuteInfo: jobExecuteInfo,
					Output:      make([]byte, 0),
					StartTime:   time.Now(),
				}

				result.StartTime = time.Now()
				// 模拟执行....

				time.Sleep(100*time.Millisecond)
				//记录结束时间
				result.EndTime = time.Now()
				result.Output = bytes.NewBufferString("test output").Bytes()
				result.Err = err
				logs.Debug("模拟执行: ",jobExecuteInfo.Job.Name," 开始时间:",result.StartTime ," 结束时间:",result.EndTime , " 耗时:",result.EndTime.Sub(result.StartTime))
				return
			});

			//logs.Debug("sleep...",t.Seconds())
			time.Sleep(10*time.Millisecond)

		}
	}()


	time.Sleep(6000*time.Second)
	t.Log("run over...")

}
