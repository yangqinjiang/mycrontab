package worker

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"math/rand"
	"os"
	"strconv"
	"testing"
	"time"
)

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
	SIZE := 100*10000 //
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
func TestHeapSortRemoveSort(t *testing.T) {
	//logs.SetLevel(logs.LevelInfo)
	// 简单的性能测试,如下:
	// 1w, 耗时 0 ms.
	// 10w, 耗时 0 ms.
	// 100w, 耗时 0 ms.
	SIZE := 1*100 //
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
			logs.Info("插入一条数据,并排序:",one_job.Name," took :",  elapsed)
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
	//测试删除
	//for i := 1; i <= SIZE/2; i++  {
	//	istr := strconv.Itoa(i)
	//	j.Remove("job_"+istr) //删除
	//}
	//检查是否排序
	if j.IsSorted(){
		t.Log("删除后,排序正常")
	}else{
		t.Fatal("删除后,排序不对")
	}

}
//找出最早的
func TestExtractEarliestHeap(t *testing.T) {

	// 简单的性能测试,如下:
	// 1w, 耗时 0 ms.
	// 10w, 耗时 0 ms.
	// 100w, 耗时 0 ms.
	SIZE := 10
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
			logs.Info("插入一条数据,并排序:",one_job.Name," took :",  elapsed)
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
		for {
			logs.Info("")
			logs.Info("for...",os.Getpid())
			//随机删除一些数据
			//remove_err := j.Remove("job_"+strconv.Itoa(randInt(1,j.Size())))
			//if remove_err != nil{
			//	logs.Error(remove_err)
			//}
			miniTime, err1 := j.ExtractEarliest(func(jobPlan common.JobSchedulePlan) (err error) {
				logs.Info("执行任务", jobPlan.Job.Name," ,after ", jobPlan.NextTime.Sub(time.Now()))
				return nil
			})
			if err1 != nil {
				t.Error("ExtractEarliest err=", err1)
			}


			logs.Info("sleep ", miniTime.Seconds(), "s","...end")

			time.Sleep(miniTime)


		}
	}()

	time.Sleep(60 * time.Second)
	//<- ending
	t.Log("run over...")

}
func randInt(min int , max int) int {
	rand.Seed( time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

