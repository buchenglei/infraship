package skeleton

import (
	"strings"

	"github.com/rs/zerolog"
)

type Env uint8

const (
	Debug Env = iota
	Test
	Prod
	Profiling
)

func ParseEnvString(env string) (e Env) {
	switch strings.ToLower(env) {
	case "debug":
		e = Debug
	case "test":
		e = Test
	case "prod":
		e = Prod
	case "profiling":
		e = Profiling
	default:
		e = Debug
	}

	return
}

type LogLevel = zerolog.Level

// 定义常用的日志等级
var (
	DebugLevel     = zerolog.DebugLevel
	TestLevel      = zerolog.DebugLevel
	ProdLevel      = zerolog.InfoLevel
	ProfilingLevel = zerolog.NoLevel
)
