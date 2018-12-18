package worker

import (
	"github.com/yangqinjiang/mycrontab/crontab/common"
	"os/exec"
	"time"
)

//任务执行器
type Executor struct {
}

var (
	G_executor *Executor
)

//执行一个任务
func (executor *Executor) ExecuteJob(info *common.JobExecuteInfo) {
	//启动协程
	go func() {
		var (
			cmd     *exec.Cmd
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
		err = jobLock.TryLock()
		//释放锁
		defer jobLock.Unlock()
		if err != nil {
			result.Err = err
			result.EndTime = time.Now()
		} else {
			//上锁成功后,重置任务启动时间
			result.StartTime = time.Now()
			//执行shell命令
			cmd = exec.CommandContext(info.CancelCtx, "/bin/bash", "-c", info.Job.Command)
			//执行并捕获输出
			output, err = cmd.CombinedOutput()
			//记录结束时间
			result.EndTime = time.Now()
			result.Output = output
			result.Err = err
		}
		//无论是否抢到锁,都返回
		//任务执行完成后,把执行的结果返回给scheduler,它会从executingTable删除执行记录
		G_scheduler.PushJobResult(result)
	}()
}

//初始化执行器
func InitExecutor() (err error) {

	G_executor = &Executor{}
	return
}
