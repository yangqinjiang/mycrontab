package worker

import (
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"testing"
)

func TestJobPlanArray(t *testing.T) {
	j  := NewJobPlanArray()
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
		t.Fatal("Remove 失败,数量 != 2")
	}
}
