package job_plan

import ("time"
	"github.com/yangqinjiang/mycrontab/worker/common")

//任务计划 接口
type JobPlanManager interface {
	Size() int
	//插入一个任务
	Insert(info *common.JobSchedulePlan)(error)
	// 使用key 删除一个任务
	Remove(key string,newItem  *common.Job)(error)
	//找出最早
	ExtractEarliest(func (jobPlan *common.JobSchedulePlan)(err error)) (time.Duration,error)
}