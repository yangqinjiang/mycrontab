package log
import(
	"github.com/yangqinjiang/mycrontab/worker/common"
)
//日志接口类
type JobLoger interface {
	Write(jobLog *common.LogBatch) (n int, err error)
}