package worker

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestJobPlanHeapIsJobPlanManager(t *testing.T)  {

	SIZE := 1
	j := NewJobPlanMinHeap(SIZE)
	_,ok:=interface{}(j).(JobPlanManager)

	if !ok{
		t.Fatal("JobPlanMinHeap没有实现JobPlanManager接口的方法")
	}
}
func TestJobPlanHeap(t *testing.T) {

	SIZE := 3
	j := NewJobPlanMinHeap(SIZE)
	job_1 := &common.Job{Name: "job_1", CronExpr: "* * * * *"}
	j_1_0, _ := common.BuildJobSchedulePlan(job_1)
	j_1_1, _ := common.BuildJobSchedulePlan(job_1)
	//插入同一个job
	j.Insert(j_1_0)
	j.Insert(j_1_1)
	if 1 != j.Size() {
		t.Fatal("Insert 失败,数量 != 1")
	}
	//插入另外两个job
	job_2 := &common.Job{Name: "job_2", CronExpr: "* * * * *"}
	job_3 := &common.Job{Name: "job_3", CronExpr: "* * * * *"}
	jj_1, _ := common.BuildJobSchedulePlan(job_1)
	jj_2, _ := common.BuildJobSchedulePlan(job_2)
	jj_3, _ := common.BuildJobSchedulePlan(job_3)
	j.Insert(jj_1)
	j.Insert(jj_2)
	j.Insert(jj_3)
	if 3 != j.Size() {
		t.Fatal("Insert 失败,数量 != 3")
	}
	if !j.SizeTrue(){
		t.Fatal("三种容器的len不一致")
	}
}

//找出最早的
func TestHeapSort(t *testing.T) {
	logs.SetLevel(logs.LevelInfo)
	// 简单的性能测试,如下:
	// 1w, 耗时 0 ms.
	// 10w, 耗时 0 ms.
	// 100w, 耗时 0 ms.
	SIZE := 1*10000 //
	j := NewJobPlanMinHeap(SIZE)
	for i := 1; i <= SIZE; i++ {
		istr := strconv.Itoa(i)
		i_60 := i
		if i >= 60 { //大于 60,求余数
			i_60 = i % 60
		}
		if i_60 <= 0 {
			i_60 = 1
		}

		i_60_str := strconv.Itoa(i_60)
		one_job := &common.Job{Name: "job_" + istr, CronExpr: "*/" + i_60_str + " * * * * * *"}
		jj, err := common.BuildJobSchedulePlan(one_job)
		if err == nil {
			startTime := time.Now()
			err = j.Insert(jj)
			elapsed := time.Since(startTime)
			logs.Debug("插入一条数据,并排序:",one_job.Name," took :",  elapsed)
			if err != nil {
				t.Error(err.Error())
			}
		} else {
			t.Error(err.Error())
		}

	}
	if SIZE != j.Size() {
		t.Fatal("Insert 失败,数量 != ", SIZE)
	}
	if !j.SizeTrue(){
		t.Fatal("三种容器的len不一致")
	}

	logs.Info("排序...")
	startTime := time.Now()
	//检查是否排序
	if j.IsSorted(){
		t.Log("插入元素后,排序正常")
	}else{
		t.Fatal("插入元素后,排序不对")
	}
	elapsed1 := time.Since(startTime)
	fmt.Printf(" Sort Test ,took %s%s",  elapsed1, "  \n")
	logs.Info("排序 over")

}
//找出最早的
func TestExtractEarliestHeap(t *testing.T) {
	logs.SetLevel(logs.LevelInfo)
	// 简单的性能测试,如下:
	// 1w, 耗时 0 ms.
	// 10w, 耗时 0 ms.
	// 100w, 耗时 0 ms.
	SIZE := 60
	j := NewJobPlanMinHeap(SIZE)
	for i := 1; i <= SIZE; i++ {
		istr := strconv.Itoa(i)
		i_60 := i
		if i >= 60 { //大于 60,求余数
			i_60 = i % 60
		}
		if i_60 <= 0 {
			i_60 = 1
		}

		//i_60_str := "1"
		i_60_str := strconv.Itoa(i_60)
		one_job := &common.Job{Name: "job_" + istr, CronExpr: "*/" + i_60_str + " * * * * * *"}
		jj, err := common.BuildJobSchedulePlan(one_job)
		if err == nil {
			startTime := time.Now()
			err = j.Insert(jj)
			elapsed := time.Since(startTime)
			logs.Debug("插入一条数据,并排序:",one_job.Name," took :",  elapsed)
			if err != nil {
				t.Error(err.Error())
			}
		} else {
			t.Error(err.Error())
		}

	}
	if SIZE != j.Size() {
		t.Fatal("Insert 失败,数量 != ", SIZE)
	}

	go func() {
		scheduleTimer := time.NewTimer(1*time.Second)
		for {
			select {
				case <-scheduleTimer.C: //最近的任务到期了
			}
			logs.Info("")
			logs.Debug("for...",os.Getpid())
			miniTime, err1 := j.ExtractEarliest(func(jobPlan *common.JobSchedulePlan) (err error) {
				logs.Info("执行任务", jobPlan.Job.Name, " ,本次执行时间=", jobPlan.NextTime)
				return nil
			})
			if err1 != nil {
				t.Error("ExtractEarliest err=", err1)
			}


			logs.Info("sleep ", miniTime.Seconds(), "s","...end")

			scheduleTimer.Reset(miniTime)
			/*
			if miniTime.Seconds() <= 0{
				time.Sleep(100*time.Millisecond)
			}else{
				time.Sleep(miniTime)
			}
			*/



		}
	}()

	time.Sleep(60*2 * time.Second)
	//<- ending
	t.Log("run over...")

}


