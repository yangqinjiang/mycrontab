package job_plan

import (
	"fmt"
	"github.com/pkg/errors"
	logs "github.com/sirupsen/logrus"
	"github.com/yangqinjiang/mycrontab/worker/common"
	"time"
)

//使用最小堆， 动态排序任务
//  参考资料
//  https://github.com/liuyubobobo/Play-with-Algorithms/blob/master/04-Heap/Course%20Code%20(C%2B%2B)/Optional-2-Min-Heap/MinHeap.h
//任务计划表 ,使用最小堆实现
type JobPlanMinHeap struct {
	JobPlanManager //任务计划管理
	indexes []int64    // index=>时间
	//keyIndex       map[string]int
	count          int
	capacity       int

	jobPlanTable   []common.JobSchedulePlan //使用map,保存任务调度计划元素
	jobPlanMap   map[string]int //使用map,保存任务调度计划元素
}

func (j *JobPlanMinHeap)SizeTrue() bool  {
	logs.Info("indexes len -1 =",len(j.indexes)-1," ,jobPlanTable len -1 =",len(j.jobPlanTable)-1, " ,jobPlanMap len=",len(j.jobPlanMap))
	//
	return len(j.indexes) -1  == len(j.jobPlanTable) - 1 && len(j.jobPlanTable) - 1 == len(j.jobPlanMap)
}
func (j *JobPlanMinHeap) PrintList() {
	for i := 1; i <= j.Size(); i++ {
		//logs.Debug(j.indexes[i].Job, j.indexes[i].NextTime, j.indexes[i].Job.CronExpr)
	}

}
func (j *JobPlanMinHeap) ExtractEarliest(tryStartJob func(jobPlan *common.JobSchedulePlan) (err error)) (t time.Duration, err error) {
	var (
		mini_plan         common.JobSchedulePlan
		mini_plan_extract common.JobSchedulePlan
	)
	//计算时间
	now := time.Now()

	//获取堆顶元素
	mini_plan = j.GetMin()
	if nil == mini_plan.Job {
		return 0, nil
	}
	logs.Info("ExtractEarliest min item=",mini_plan.Job.Name)
	//判断是否快过期
	isExpire := mini_plan.NextTime.Before(now) || mini_plan.NextTime.Equal(now)
	elapsed := time.Since(now)
	//这里执行任务
	if isExpire && nil != tryStartJob {

		now = time.Now()
		//从最小堆中取出堆顶元素
		mini_plan_extract = j.ExtractMin()
		elapsed = time.Since(now) //更新遍历时间
		logs.Warn("取出最小堆顶元素,任务名称= [ ",mini_plan_extract.Job.Name," ] Command=[",mini_plan_extract.Job.Command," ] ,ShellName=[",mini_plan_extract.Job.ShellName," ] , CronExpr = [",mini_plan_extract.Job.CronExpr, " ],执行时间= [ ", mini_plan_extract.NextTime," ],并重新进行 shiftDown,耗时: ", elapsed)
		if mini_plan_extract.Del{
			logs.Error("已标识为DEL,跳过,不执行任务")
			return 0,nil
		}

		if mini_plan.Job != mini_plan_extract.Job {
			panic(fmt.Sprintf("从GetMin得到的[ %s ]Job 与 ExtractMin得到的[ %s ] Job 不是同一个",mini_plan.Job.Name,mini_plan_extract.Job.Name))
		}

		logs.Warn("已过期, 准备执行任务...")
		//尝试执行任务
		tryStartJob(&mini_plan_extract)
		//执行后,更新下次执行时间的值
		mini_plan_extract.NextTime = mini_plan_extract.Expr.Next(now)
		logs.Warn("执行完成,更新下次执行时间的值 ,下次执行时间=", mini_plan_extract.NextTime)

		//再次插入数据
		if err := j.InsertAgain(&mini_plan_extract); err != nil {
			logs.Error(err)
			return 0, err
		}

	} else {
		logs.Info("未过期或者没设置tryStartJob参数")
	}

	return mini_plan.NextTime.Sub(now), nil //返回最小的时间,用于睡眠或定时
}
func (mh *JobPlanMinHeap) shiftUp(k int) {
	for k > 1 && mh.jobPlanTable[mh.indexes[k/2]].NextTime.Unix() > mh.jobPlanTable[mh.indexes[k]].NextTime.Unix() {
		mh.swap(&mh.indexes[k/2], &mh.indexes[k])
		k /= 2
	}
}
func (mh *JobPlanMinHeap) shiftDown(k int) {
	for 2*k <= mh.count {
		j := 2 * k //在此轮循环中,indexes[k]和data[j]交换位置
		if j+1 <= mh.count && mh.jobPlanTable[mh.indexes[j+1]].NextTime.Unix()< mh.jobPlanTable[mh.indexes[j]].NextTime.Unix() {
			//左右子节点最小的一个
			j++
		}
		// indexes[j] 是 indexes[2*k]和data[2*k+1]中的最小值
		if mh.jobPlanTable[mh.indexes[k]].NextTime.Unix() <= mh.jobPlanTable[mh.indexes[j]].NextTime.Unix() {
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
//再次插入原有数据
func (mh *JobPlanMinHeap)InsertAgain(item *common.JobSchedulePlan) error  {
	index,exist := mh.jobPlanMap[item.Job.Name] //
	if !exist{
		return errors.New("不存在此Job,Name="+item.Job.Name)
	}
	
	//index+=1
	mh.jobPlanTable[index]= *item
	// 找到indexes[j] = i, j表示data[i]在堆中的位置
	// 之后shiftUp(j), 再shiftDown(j)
	count := mh.count + 1
	for j:=1;j<=count;j++{
		if mh.indexes[j] == int64(index){
			logs.Debug("存在mh.indexes[j] == int64(index)",j,index," ,mh.count=",mh.count)
			mh.count ++
			mh.shiftUp(j)
			mh.shiftDown(j)
			return nil
		}else{
			logs.Debug("indexes不存在 mh.indexes=",mh.indexes," ,index=",index," ,j=",j," , count=",count)
		}
	}
	return nil
}
// 向最小堆中插入一个新的元素 item
func (mh *JobPlanMinHeap) Insert(item *common.JobSchedulePlan) error {

	if index,exist:=mh.jobPlanMap[item.Job.Name];exist{
		mh.jobPlanTable[index] = *item
		return errors.New("已存在 "+item.Job.Name+" 任务,并更新它")
	}

	//边界
	i := mh.count + 1
	logs.Debug(" Before Insert mh.count=",mh.count," ,i=",i)
	if !(i <= mh.capacity)  {
		return errors.New("任务数量 已超出额定边界")
	}
	if !(i >= 1 && i <= mh.capacity ){
		return errors.New("任务数量 已超出额定边界")
	}

	logs.Debug("插入新的JobPlan值 Job.Name=", item.Job.Name)

	//保存到索引数组
	mh.indexes[i] = int64(i)
	logs.Debug("Insert func ,Indexes=",mh.indexes)
	// 保存到map
	mh.jobPlanTable[i] = *item  //
	mh.jobPlanMap[item.Job.Name] = i //保存JobName为索引

	mh.count++
	mh.shiftUp(mh.count)
	logs.Debug("再次插入mini_plan的值item.Job.Name=", item.Job.Name,item.NextTime,item.Job.CronExpr, ",mh.count=", mh.count)

	return nil
}

// 使用key 删除一个任务
// 将最小堆中索引为i的元素修改为newItem
func (mh *JobPlanMinHeap) Remove(key string,newItem *common.Job) error {

	//存在,则修改
	var index int;
	index , exist := mh.jobPlanMap[key]
	if !exist{
		err_str := "Remove,不存在key="+key+"的Job"
		logs.Error(err_str)
		return errors.New(err_str)
	}
	jobSchedulePlan, err := common.BuildDeleteJobSchedulePlan(newItem)
	if err != nil {
		return err
	}

	i := int64(index)
	logs.Warn("Remove & change",key,"的Job"," ,index=",i)
	mh.jobPlanTable[i]= *jobSchedulePlan
	delete(mh.jobPlanMap,key)
	return nil
}

// 从最小堆中取出堆顶元素, 即堆中所存储的最小数据
func (e *JobPlanMinHeap) ExtractMin() common.JobSchedulePlan {
	Assert(e.count > 0)

	 //读取第一个,是最小值
	plan := e.jobPlanTable[e.indexes[1]]

	//logs.Debug("ExtractMin: Before Swap ,Job.Name= ", job_plan.Job.Name)

	//交换最后和第一个元素,使它不是最小堆
	e.swap(&e.indexes[1], &e.indexes[e.count])
	e.count--
	//进行 shiftDown
	e.shiftDown(1)
	//delete(e.jobPlanMap, job_plan.Job.Name) //删除key_value
	//返回第一个
	logs.Debug("ExtractMin: After Swap ,e.indexes=", e.indexes," ,e.count=",e.count)
	return plan

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
	logs.Debug("Before ExtractMin,data_size=", data_size)
	data := make([]common.JobSchedulePlan, data_size)
	//从堆中依次取出元素
	for i := 0; i < data_size; i++ {
		logs.Debug("Before ExtractMin,Size=", i, " ,Size=", e.Size())
		data[i] = e.ExtractMin()
	}
	logs.Debug("After ExtractMin,Size=", e.Size())
	for i := 0; i < data_size-1; i++ {
		logs.Debug(data[i].Job.Name, data[i].NextTime.Unix())
		if data[i].NextTime.Unix() > data[i+1].NextTime.Unix() {
			logs.Error("ERROR: ", data[i+1].Job.Name, data[i+1].NextTime.Unix())
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
