package app

import (
	"os"
	"strconv"

	"github.com/buchenglei/infraship/skeleton"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Builder interface {
	Build() zerolog.Logger
	ReplaceGlobalLogger()
}

type Logger struct {
	l zerolog.Logger
}

func NewLogger(name string, env skeleton.Env) *Logger {
	logger := &Logger{}

	var level skeleton.LogLevel
	switch env {
	case skeleton.Debug:
		level = skeleton.DebugLevel
		logger.l = buildLoggerContextDevMode(name)
	case skeleton.Test:
		level = skeleton.TestLevel
		logger.l = buildLoggerContextTestMode(name)
	default:
		level = skeleton.ProdLevel
		logger.l = buildLoggerContextProdMode(name)
	}

	zerolog.SetGlobalLevel(level)

	// 通用设置
	zerolog.TimestampFieldName = "loc_time"

	return logger
}

func (l *Logger) Build() zerolog.Logger {
	return l.l
}

func (l *Logger) ReplaceGlobalLogger() {
	log.Logger = l.l
}

func buildLoggerContextDevMode(appname string) zerolog.Logger {
	zerolog.CallerMarshalFunc = cutCallerFileName

	return zerolog.New(os.Stdout).With().
		Str("app_name", appname).Timestamp().
		Caller().Stack().
		Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func buildLoggerContextTestMode(appname string) zerolog.Logger {
	zerolog.CallerMarshalFunc = cutCallerFileName

	return zerolog.New(os.Stdout).With().
		Str("app_name", appname).Timestamp().Stack().
		CallerWithSkipFrameCount(2).Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func buildLoggerContextProdMode(appname string) zerolog.Logger {
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"

	return zerolog.New(os.Stdout).With().
		Str("app_name", appname).Timestamp().
		Logger()
}

func buildLoggerContextProfilingMode(appname string) zerolog.Logger {
	return zerolog.New(os.Stdout)
}

func cutCallerFileName(pc uintptr, file string, line int) string {
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short
	return file + ":" + strconv.Itoa(line)
}
