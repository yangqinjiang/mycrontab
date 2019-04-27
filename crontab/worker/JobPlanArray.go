package worker

import "github.com/yangqinjiang/mycrontab/crontab/common"

//任务计划表 ,使用最小堆实现
type JobPlanArray struct {
	JobPlanManager //任务计划管理
	jobPlanTable      map[string]*common.JobSchedulePlan //任务调度计划表内存里的任务计划表,
}
func (j *JobPlanArray)Size() int{
	return len(j.jobPlanTable)
}
//插入一个任务
func (j *JobPlanArray)Insert(info *common.JobSchedulePlan)(error){
	j.jobPlanTable[info.Job.Name] = info
	return nil
}
// 使用key 删除一个任务
func (j *JobPlanArray)Remove(key string)(error){
	//首先检查是否已存在此任务 (内存)
	if _, jobExisted := j.jobPlanTable[key]; jobExisted {
		delete(j.jobPlanTable, key)
	}
	return nil
}
func NewJobPlanArray() *JobPlanArray {
	return &JobPlanArray{
		jobPlanTable:      make(map[string]*common.JobSchedulePlan, 1000), //内存里的任务计划表,
	}
}
