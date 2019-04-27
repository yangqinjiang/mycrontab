package worker

import (
	"github.com/astaxie/beego/logs"
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"time"
)

//TODO: 使用最小堆， 动态排序任务
//  参考资料
//  https://github.com/liuyubobobo/Play-with-Algorithms/blob/master/04-Heap/Course%20Code%20(C%2B%2B)/Optional-2-Min-Heap/MinHeap.h
//任务计划表 ,使用最小堆实现
type JobPlanMinHeap struct {
	JobPlanManager //任务计划管理
	data []int
	count int
	capacity int
}
func (j *JobPlanMinHeap) ExtractEarliest(tryStartJob func(jobPlan *common.JobSchedulePlan) (err error)) time.Duration {
	return 0
}
func (mh *JobPlanMinHeap)shiftUp(k int)  {
	for k> 1 && mh.data[k/2] > mh.data[k]  {
		mh.swap(&mh.data[k/2],&mh.data[k])
		k /= 2
	}
}
func (mh *JobPlanMinHeap)shiftDown(k int)  {
	for 2*k <= mh.count  {
		j:= 2*k//在此轮循环中,data[k]和data[j]交换位置
		if j+1 <= mh.count && mh.data[j+1] < mh.data[j] {
			//左右子节点最小的一个
			j++
		}
		// data[j] 是 data[2*k]和data[2*k+1]中的最小值
		if mh.data[k] <= mh.data[j] {
			break
		}

		mh.swap(&mh.data[k],&mh.data[j])
		k = j
	}
}
/**
交换数组或切片的两个元素
*/
func (e *JobPlanMinHeap) swap(a, b *int) {
	*a, *b = *b, *a
}

// 返回堆中的元素个数
func (mh *JobPlanMinHeap)Size() int  {
	return mh.count
}
// 返回一个布尔值, 表示堆中是否为空
func (mh *JobPlanMinHeap)IsEmpty() bool  {
	return mh.count == 0
}
// 向最小堆中插入一个新的元素 item
func (mh *JobPlanMinHeap)Insert(item int)   {
	//边界
	Assert(mh.count + 1<= mh.capacity)
	mh.data[mh.count + 1] = item
	mh.shiftUp(mh.count)
	mh.count ++
}
// 从最小堆中取出堆顶元素, 即堆中所存储的最小数据
func (e *JobPlanMinHeap)ExtractMin() int  {
	Assert(e.count > 0)
	ret := e.data[1] //读取第一个,是最大值

	//交换最后和第一个元素,使得不是最大堆
	e.swap(&e.data[1],&e.data[e.count])
	e.count --
	//进行 shiftDown
	e.shiftDown(1)

	//返回第一个
	return  ret

}
// 获取最小堆中的堆顶元素
func (e *JobPlanMinHeap)GetMin() int  {
	Assert(e.count > 0)
	return e.data[1]
}
// 构造函数, 构造一个空堆, 可容纳capacity个元素
func NewJobPlanMinHeap(capacity int) *JobPlanMinHeap  {
	logs.Debug("NewJobPlanMinHeap");
	return &JobPlanMinHeap{data:make([]int,capacity+1,capacity+1),
		count:0,
		capacity:capacity}
}
//Heapify
// 构造函数, 通过一个给定数组创建一个最小堆
// 该构造堆的过程, 时间复杂度为O(n)
func NewJobPlanMinHeapByArray(arr []int,n int) *JobPlanMinHeap  {

	// 索引从1开始
	o := &JobPlanMinHeap{data:make([]int,n+1,n+1),
		count:0,
		capacity:n}

	//更新堆元素
	for i:=0;i<n;i++{
		o.data[i+1] = arr[i]
	}
	//更新计数器
	o.count=n

	//从不是子节点开始,进行shiftDown
	for i:= o.count/2 ;i>=1;i--{
		o.shiftDown(i)
	}


	return o
}
func Assert(b bool)  {
	if(!b){
		panic("出错了")
	}
}
