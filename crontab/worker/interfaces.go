package worker

import (
	"io"
)

//日志接口类
type Log interface {
	io.Writer
}
