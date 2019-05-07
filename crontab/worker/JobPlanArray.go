package worker

import (
	"github.com/astaxie/beego/logs"
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"time"
)

//任务计划表 ,使用最小堆实现
type JobPlanArray struct {
	JobPlanManager                                  //任务计划管理
	jobPlanTable   map[string]*common.JobSchedulePlan //任务调度计划表内存里的任务计划表,
}

func (j *JobPlanArray) Size() int {
	return len(j.jobPlanTable)
}

//插入一个任务
func (j *JobPlanArray) Insert(info *common.JobSchedulePlan) error {
	j.jobPlanTable[info.Job.Name] = info
	return nil
}

// 使用key 删除一个任务
func (j *JobPlanArray) Remove(key string) error {
	//首先检查是否已存在此任务 (内存)
	if _, jobExisted := j.jobPlanTable[key]; jobExisted {
		delete(j.jobPlanTable, key)
	}
	return nil
}

//找出最早的
func (j *JobPlanArray) ExtractEarliest(tryStartJob func(jobPlan *common.JobSchedulePlan) (err error)) time.Duration {
	now := time.Now()
	var nearTime *time.Time
	//计算 一次foreach的计算时间
	startTime := time.Now()
	for _, jobPlan := range j.jobPlanTable {
		//是否过期,小于或都等于当前时间
		if jobPlan.NextTime.Before(now) || jobPlan.NextTime.Equal(now) {
			if nil != tryStartJob{
				//尝试执行任务
				//tryStartJob(jobPlan)
				jobPlan.NextTime = jobPlan.Expr.Next(now) //执行后,更新下次执行时间的值
			}


		}
		//统计最近一个要过期的任务时间
		if nearTime == nil { // 刚开始for第一个,则设置一个值
			nearTime = &jobPlan.NextTime
		}
		//判断第二个及以后,如果是更往后的时刻,则更新它,
		// 即 找出更晚的时刻
		if jobPlan.NextTime.Before(*nearTime) {
			nearTime = &jobPlan.NextTime
		}
	}

	logs.Debug("JobPlanArray ForEach 遍历耗时: ",time.Since(startTime))
	//下次调度间隔,(最近要执行的任务调度时间 - 当前时间)
	return (*nearTime).Sub(now)
}
func NewJobPlanArray() *JobPlanArray {
	logs.Debug("NewJobPlanArray");
	return &JobPlanArray{
		jobPlanTable: make(map[string]*common.JobSchedulePlan, 1000), //内存里的任务计划表,
	}
}
