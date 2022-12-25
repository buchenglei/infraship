package skeleton

import "context"

// Pinger 用于检测各个外部服务或功能组件是否可以正常联通
// 同样可以用来描述服务内部异步运行的模块是否正常
type Pinger interface {
	Ping(context.Context) error
}
