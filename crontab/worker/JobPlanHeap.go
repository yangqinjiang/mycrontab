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
	JobPlanManager //任务计划管理
	data           []common.JobSchedulePlan
	keyIndex       map[string]int
	count          int
	capacity       int
}

func (j *JobPlanMinHeap) PrintList() {
	for i := 1; i <= j.Size(); i++ {
		logs.Debug(j.data[i].Job, j.data[i].NextTime, j.data[i].Job.CronExpr)
	}

}
func (j *JobPlanMinHeap) ExtractEarliest(tryStartJob func(jobPlan *common.JobSchedulePlan) (err error)) (t time.Duration, err error) {
	var (
		mini_plan         *common.JobSchedulePlan
		mini_plan_extract *common.JobSchedulePlan
	)
	//计算时间
	now := time.Now()
	before_len := j.count

	//获取堆顶元素
	mini_plan = j.GetMin()
	if nil == mini_plan {
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
		logs.Debug("取出最小堆顶元素 item_1=", mini_plan_extract.Job.Name, "已过期, 准备执行任务...")

		elapsed = time.Since(now) //更新遍历时间
		//尝试执行任务
		tryStartJob(mini_plan_extract)
		mini_plan_extract.NextTime = mini_plan_extract.Expr.Next(now) //执行后,更新下次执行时间的值

		if err := j.Insert(mini_plan_extract); err != nil {
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
	for k > 1 && mh.data[k/2].NextTime.Unix() > mh.data[k].NextTime.Unix() {
		mh.swap(&mh.data[k/2], &mh.data[k])
		k /= 2
	}
}
func (mh *JobPlanMinHeap) shiftDown(k int) {
	for 2*k <= mh.count {
		j := 2 * k //在此轮循环中,data[k]和data[j]交换位置
		if j+1 <= mh.count && mh.data[j+1].NextTime.Unix() < mh.data[j].NextTime.Unix() {
			//左右子节点最小的一个
			j++
		}
		// data[j] 是 data[2*k]和data[2*k+1]中的最小值
		if mh.data[k].NextTime.Unix() <= mh.data[j].NextTime.Unix() {
			break
		}

		mh.swap(&mh.data[k], &mh.data[j])
		k = j
	}
}

/**
交换数组或切片的两个元素
*/
func (e *JobPlanMinHeap) swap(a, b *common.JobSchedulePlan) {
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
	logs.Debug("插入新的JobPlan值 Job.Name=", item.Job.Name)
	//mh.PrintList()
	//边界
	myIndex := mh.count + 1

	if myIndex > mh.capacity{
		return errors.New("任务已超出额定数量")
	}
	//如果已存在这个元素,则跳过
	if _, exist := mh.keyIndex[item.Job.Name]; exist {
		logs.Error("再次插入JobPlan 失败")
		return errors.New("已存在这个数据")
	}

	mh.data[myIndex] = *item
	mh.keyIndex[item.Job.Name] = myIndex
	mh.shiftUp(mh.count)
	mh.count++
	logs.Debug("再次插入mini_plan的值item.Job.Name=", item.Job.Name, ",mh.count=", mh.count)
	//mh.PrintList()
	return nil
}

// 使用key 删除一个任务
func (mh *JobPlanMinHeap) Remove(key string) error {
	//存在,则删除 索引为myIndex的数据
	if myIndex, exist := mh.keyIndex[key]; exist {


		logs.Info("exist ,remove",key," myIndex=",myIndex," data.Size:", len(mh.data))
		//操作切片
		if myIndex >= len(mh.data){
			mh.data = mh.data[:myIndex]
		}else{ // myIndex < len(mh.data)
			mh.data = append(mh.data[:myIndex], mh.data[myIndex+1:]...)
		}

		//删除keyIndex中的数据
		delete(mh.keyIndex, key)
		mh.count--
		//交换最后和第一个元素,使它不是最小堆
		mh.swap(&mh.data[1], &mh.data[mh.count])
		//进行 shiftDown
		mh.shiftDown(1)
		return nil
	}

	return errors.New("删除失败,不存在这个数据key= " + key)

}

// 从最小堆中取出堆顶元素, 即堆中所存储的最小数据
func (e *JobPlanMinHeap) ExtractMin() *common.JobSchedulePlan {
	Assert(e.count > 0)

	ret := e.data[1] //读取第一个,是最小值

	logs.Debug("ExtractMin: Before Swap ,Job.Name= ",ret.Job.Name )

	//交换最后和第一个元素,使它不是最小堆
	e.swap(&e.data[1], &e.data[e.count])
	e.count--
	//进行 shiftDown
	e.shiftDown(1)
	delete(e.keyIndex, ret.Job.Name) //删除key_value
	//返回第一个
	logs.Debug("ExtractMin: After Swap ,Job.Name=",ret.Job.Name )
	return &ret

}

// 将最小堆中索引为i的元素修改为newItem
func (e *JobPlanMinHeap) change(key string, newItem common.JobSchedulePlan) {
	//存在,则修改
	if myIndex, exist := e.keyIndex[key]; exist {
		i := myIndex + 1
		e.data[i] = newItem
		e.shiftDown(i)
	}
}

// 获取最小堆中的堆顶元素
func (e *JobPlanMinHeap) GetMin() *common.JobSchedulePlan {
	if e.count <= 0 {
		return nil
	}
	return &e.data[1]
}

// 判断arr数组是否有序
func (e *JobPlanMinHeap) IsSorted() bool {

	data_size := e.Size() ;
	logs.Info("Before ExtractMin,data_size=",data_size)
	data := make([]*common.JobSchedulePlan,data_size)
	//从堆中依次取出元素
	for i := 0; i < data_size; i++ {
		logs.Info("Before ExtractMin,Size=",i," ,Size=",e.Size())
		data[i] = e.ExtractMin()
	}
	logs.Info("After ExtractMin,Size=",e.Size(),data)
	for i := 0; i <  data_size - 1; i++ {
		logs.Info(data[i].Job.Name,data[i].NextTime.Unix())
		if data[i].NextTime.Unix() > data[i+1].NextTime.Unix() {

			return false
		}
	}
	return true
}

// 构造函数, 构造一个空堆, 可容纳capacity个元素
func NewJobPlanMinHeap(capacity int) *JobPlanMinHeap {
	logs.Debug("NewJobPlanMinHeap")
	return &JobPlanMinHeap{data: make([]common.JobSchedulePlan, capacity+1, capacity+1),
		keyIndex: make(map[string]int, capacity+1),
		count:    0,
		capacity: capacity}
}

func Assert(b bool) {
	if !b {
		panic("出错了")
	}
}
