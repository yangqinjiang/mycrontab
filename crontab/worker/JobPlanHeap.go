package worker

import (
	"github.com/astaxie/beego/logs"
	"github.com/pkg/errors"
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"time"
)

//使用最小堆， 动态排序任务
//  参考资料
//  https://github.com/liuyubobobo/Play-with-Algorithms/blob/master/04-Heap/Course%20Code%20(C%2B%2B)/Optional-2-Min-Heap/MinHeap.h
//任务计划表 ,使用最小堆实现
type JobPlanMinHeap struct {
	JobPlanManager  //任务计划管理
	indexes []int64 // index=>时间
	//keyIndex       map[string]int
	count          int
	capacity       int

	jobPlanTable   []common.JobSchedulePlan //使用map,保存任务调度计划元素
	jobPlanMap   map[string]int //使用map,保存任务调度计划元素
}

func (j *JobPlanMinHeap) PrintList() {
	for i := 1; i <= j.Size(); i++ {
		//logs.Debug(j.indexes[i].Job, j.indexes[i].NextTime, j.indexes[i].Job.CronExpr)
	}

}
func (j *JobPlanMinHeap) ExtractEarliest(tryStartJob func(jobPlan common.JobSchedulePlan) (err error)) (t time.Duration, err error) {
	var (
		mini_plan         common.JobSchedulePlan
		mini_plan_extract common.JobSchedulePlan
	)
	//计算时间
	now := time.Now()
	before_len := j.count

	//获取堆顶元素
	mini_plan = j.GetMin()
	if nil == mini_plan.Job {
		return 0, nil
	}
	logs.Debug("GetMin item=", mini_plan.Job.Name)
	//判断是否快过期
	isExpire := mini_plan.NextTime.Before(now) || mini_plan.NextTime.Equal(now)
	elapsed := time.Since(now)
	//这里执行任务
	if isExpire && nil != tryStartJob {
		//从最小堆中取出堆顶元素
		mini_plan_extract = j.ExtractMin()
		logs.Debug("取出最小堆顶元素 item_1=", mini_plan_extract.Job.Name, " ,执行时间=", mini_plan.NextTime, "已过期, 准备执行任务...")
		if mini_plan.Job != mini_plan_extract.Job {
			panic("从GetMin得到的Job 不等于 ExtractMin得到的Jog")
		}
		//判断执行时间与当前时间是否差太多
		//if time.Now().Sub(mini_plan_extract.NextTime).Seconds() > 1.{
		//	panic("执行时间与当前时间 的差, > 1s ")
		//}
		elapsed = time.Since(now) //更新遍历时间
		//尝试执行任务
		tryStartJob(mini_plan_extract)
		//执行后,更新下次执行时间的值
		mini_plan_extract.NextTime = mini_plan_extract.Expr.Next(now)
		logs.Debug("执行后,更新下次执行时间的值 item_1=", mini_plan_extract.Job.Name, " ,下次执行时间=", mini_plan_extract.NextTime)
		if err := j.Insert(&mini_plan_extract); err != nil {
			logs.Error(err)
			return 0, err
		}

	} else {
		logs.Debug("最小堆顶元素 item=", mini_plan.Job.Name, " 未过期")
	}

	after_len := j.count
	logs.Debug("最小堆顶元素 item=", mini_plan.Job.Name, " ,NextTime=", mini_plan.NextTime, "遍历耗时: ", elapsed, " 元素个数:(before=", before_len, "/after=", after_len, ")")
	return mini_plan.NextTime.Sub(now), nil //返回最小的时间,用于睡眠或定时
}
func (mh *JobPlanMinHeap) shiftUp(k int) {
	for k > 1 && mh.indexes[k/2] > mh.indexes[k] {
		mh.swap(&mh.indexes[k/2], &mh.indexes[k])
		k /= 2
	}
}
func (mh *JobPlanMinHeap) shiftDown(k int) {
	for 2*k <= mh.count {
		j := 2 * k //在此轮循环中,indexes[k]和data[j]交换位置
		if j+1 <= mh.count && mh.indexes[j+1] < mh.indexes[j] {
			//左右子节点最小的一个
			j++
		}
		// indexes[j] 是 indexes[2*k]和data[2*k+1]中的最小值
		if mh.indexes[k] <= mh.indexes[j] {
			break
		}

		mh.swap(&mh.indexes[k], &mh.indexes[j])
		k = j
	}
}

/**
交换数组或切片的两个元素
*/
func (e *JobPlanMinHeap) swap(a, b *int64) {
	*a, *b = *b, *a
}

// 返回堆中的元素个数
func (mh *JobPlanMinHeap) Size() int {
	return mh.count
}

// 返回一个布尔值, 表示堆中是否为空
func (mh *JobPlanMinHeap) IsEmpty() bool {
	return mh.count == 0
}

// 向最小堆中插入一个新的元素 item
func (mh *JobPlanMinHeap) Insert(item *common.JobSchedulePlan) error {

	if index,exist:=mh.jobPlanMap[item.Job.Name];exist{
		mh.jobPlanTable[index] = *item
		return errors.New("已存在 "+item.Job.Name+" 任务,并更新它")
	}
	//边界
	if !(mh.count + 1 <= mh.capacity)  {
		return errors.New("任务数量 已超出额定边界")
	}
	if !(mh.count+1 >= 1 && mh.count + 1 <= mh.capacity ){
		return errors.New("任务数量 已超出额定边界")
	}
	i := mh.count + 1
	logs.Debug("插入新的JobPlan值 Job.Name=", item.Job.Name)

	//保存到索引数组
	mh.indexes[i] = int64(i)
	// 保存到map
	mh.jobPlanTable[i] = *item  //
	mh.jobPlanMap[item.Job.Name] = i //保存JobName为索引

	mh.count++
	mh.shiftUp(mh.count)
	logs.Debug("再次插入mini_plan的值item.Job.Name=", item.Job.Name, ",mh.count=", mh.count)

	return nil
}

// 使用key 删除一个任务
func (mh *JobPlanMinHeap) Remove(key string) error {
	//存在,则更新字段值
	//if myIndex, exist := mh.keyIndex[key]; exist {
	//
	//	logs.Debug(myIndex)
	//
	//	return nil
	//}

	return errors.New("删除失败,不存在这个数据key= " + key)

}

// 从最小堆中取出堆顶元素, 即堆中所存储的最小数据
func (e *JobPlanMinHeap) ExtractMin() common.JobSchedulePlan {
	Assert(e.count > 0)

	planIndex := e.indexes[1] //读取第一个,是最小值
	plan := e.jobPlanTable[planIndex];

	logs.Debug("ExtractMin: Before Swap ,Job.Name= ", plan.Job.Name)

	//交换最后和第一个元素,使它不是最小堆
	e.swap(&e.indexes[1], &e.indexes[e.count])
	e.count--
	//进行 shiftDown
	e.shiftDown(1)
	//delete(e.keyIndex, ret.Job.Name) //删除key_value
	//返回第一个
	logs.Debug("ExtractMin: After Swap ,Job.Name=", plan.Job.Name)
	return plan

}

// 将最小堆中索引为i的元素修改为newItem
func (e *JobPlanMinHeap) change(key string, newItem common.JobSchedulePlan) {
	//存在,则修改
	var i int64;
	i = 0
	i+=1
	e.jobPlanTable[i]= newItem
	// 找到indexes[j] = i, j表示data[i]在堆中的位置
	// 之后shiftUp(j), 再shiftDown(j)
	for j:=1;j<=e.count;j++{
		if e.indexes[j] == i{
			e.shiftUp(j)
			e.shiftDown(j)
			return
		}
	}
}

// 获取最小堆中的堆顶元素
func (e *JobPlanMinHeap) GetMin() common.JobSchedulePlan {
	if e.count <= 0 {
		return common.JobSchedulePlan{}
	}
	//读取第一个,是最小值
	plan := e.jobPlanTable[e.indexes[1]];
	return plan
}

// 判断arr数组是否有序
func (e *JobPlanMinHeap) IsSorted() bool {

	data_size := e.Size()
	logs.Info("Before ExtractMin,data_size=", data_size)
	data := make([]common.JobSchedulePlan, data_size)
	//从堆中依次取出元素
	for i := 0; i < data_size; i++ {
		logs.Info("Before ExtractMin,Size=", i, " ,Size=", e.Size())
		data[i] = e.ExtractMin()
	}
	logs.Info("After ExtractMin,Size=", e.Size())
	for i := 0; i < data_size-1; i++ {
		logs.Info(data[i].Job.Name, data[i].NextTime.Unix())
		if data[i].NextTime.Unix() > data[i+1].NextTime.Unix() {
			logs.Info("ERROR: ", data[i+1].Job.Name, data[i+1].NextTime.Unix())
			return false
		}
	}
	return true
}

// 构造函数, 构造一个空堆, 可容纳capacity个元素
func NewJobPlanMinHeap(capacity int) *JobPlanMinHeap {
	logs.Debug("NewJobPlanMinHeap")
	return &JobPlanMinHeap{indexes: make([]int64, capacity+1, capacity+1),
		jobPlanTable: make([]common.JobSchedulePlan,capacity+1, capacity+1),
		jobPlanMap: make(map[string]int,capacity+1),
		count:    0,
		capacity: capacity}
}

func Assert(b bool) {
	if !b {
		panic("出错了")
	}
}
