package lib

import (
	"Worker/common"
	"math/rand"
	"sync"
	"time"
	"errors"
)

//Goroutine任务执行器,也是命令的调用者
type GoroutineExecutor struct {
	JobExecuter								// 实现任务执行的接口
	jobLock JobLocker						//任务锁锁对象
	jobResultReceiver JobResultReceiver   //任务执行结果的接收器
	command Command
}

var (
	G_GoroutineExecutor *GoroutineExecutor
	onceexec            sync.Once
)

//设置命令对象
func (t *GoroutineExecutor)SetCommand(c Command)  {
	t.command = c
}
//设置任务执行结果的接收器
func (t *GoroutineExecutor)SetJobResultReceiver(jobResultReceiver JobResultReceiver)  {
	t.jobResultReceiver = jobResultReceiver
}
//执行一个任务
func (executor *GoroutineExecutor) Exec(info *common.JobExecuteInfo) (err error) {
	//启动协程
	go func() {
		var (

			output  []byte
			err     error
			result  *common.JobExecuteResult
		)
		//任务执行的结果
		result = &common.JobExecuteResult{
			ExecuteInfo: info,
			Output:      make([]byte, 0),
			StartTime:   time.Now(),
		}

		//初始化分布式锁
		executor.jobLock = G_EtcdJobMgr.CreateJobLock(info.Job.Name)
		//随机睡眠(0~1s),解决抢锁频繁问题
		time.Sleep(time.Duration(rand.Intn(1000)) * time.Millisecond)

		//尝试抢锁
		err = executor.jobLock.TryLock()
		//释放锁
		defer executor.jobLock.Unlock()
		if err != nil {
			//抢锁失败,则返回
			result.Err = err
			result.EndTime = time.Now()
		} else {
			//上锁成功后,重置任务启动时间
			result.StartTime = time.Now()
			info.PrintStatus()
			//生成一个命令对象
			executor.SetCommand(CommandFactory(info.Job.ShellName))
			//执行命令对象
			if( nil != executor.command){
				output, err =  executor.command.Execute(info)
			}else{
				output, err = []byte{}, errors.New("不存在此类命令")
			}
			//记录结束时间
			result.EndTime = time.Now()
			result.Output = output
			result.Err = err
		}
		//无论是否抢到锁,都返回
		//任务执行完成后,把执行的结果返回给scheduler,它会从executingTable删除执行记录
		if nil != executor.jobResultReceiver{
			executor.jobResultReceiver.PushResult(result)
		}
	}()
	return
}

//初始化异步任务执行器
func InitGoroutineExecutor() (err error) {
	onceexec.Do(func() {
		G_GoroutineExecutor = &GoroutineExecutor{}
	})
	return
}
