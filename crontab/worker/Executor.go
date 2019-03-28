package worker

import (
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"math/rand"
	"sync"
	"time"
)

//任务执行器
type Executor struct {
}

var (
	G_executor *Executor
	onceexec       sync.Once
)

//执行一个任务
func (executor *Executor) Exec(info *common.JobExecuteInfo,
	callFunc func(info *common.JobExecuteInfo) ([]byte, error)) (err error) {
	//启动协程
	go func() {
		var (

			output  []byte
			err     error
			result  *common.JobExecuteResult
			jobLock *JobLock
		)
		//任务执行的结果
		result = &common.JobExecuteResult{
			ExecuteInfo: info,
			Output:      make([]byte, 0),
			StartTime:   time.Now(),
		}

		//初始化分布式锁
		jobLock = G_jobMgr.CreateJobLock(info.Job.Name)
		//随机睡眠(0~1s),解决抢锁频繁问题
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

		//尝试抢锁
		err = jobLock.TryLock()
		//释放锁
		defer jobLock.Unlock()
		if err != nil {
			//抢锁失败,则返回
			result.Err = err
			result.EndTime = time.Now()
		} else {
			//上锁成功后,重置任务启动时间
			result.StartTime = time.Now()
			//使用闭包函数,抽离主业务代码
			output, err =  callFunc(info)

			//记录结束时间
			result.EndTime = time.Now()
			result.Output = output
			result.Err = err
		}
		//TODO:使用接口优化下面的代码, 涉及到 scheduler引用executer,反过来又引用它, 相互引用  ,囧....
		//无论是否抢到锁,都返回
		//任务执行完成后,把执行的结果返回给scheduler,它会从executingTable删除执行记录
		G_scheduler.PushJobResult(result)
	}()
	return
}

//初始化执行器
func InitExecutor() (err error) {
	onceexec.Do(func() {
		G_executor = &Executor{}
	})
	return
}
