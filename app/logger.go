package app

import (
	"os"
	"path"
	"strconv"

	"github.com/buchenglei/infraship/kit"
	"github.com/buchenglei/infraship/skeleton"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LogBuilder func(l zerolog.Logger) zerolog.Logger

type Logger struct {
	l zerolog.Logger
}

func NewLogger(name string, level skeleton.LogLevel, customBuilers ...LogBuilder) *Logger {
	logger := &Logger{}

	zerolog.SetGlobalLevel(level)

	switch level {
	case skeleton.ProdLevel:
		logger.l = buildLoggerContextProdMode(name)
	default:
		logger.l = buildLoggerContextDevMode(name)
	}

	// 通用设置
	zerolog.TimestampFieldName = "loc_time"
	for _, cstBuiler := range customBuilers {
		logger.l = cstBuiler(logger.l)
	}

	return logger
}

func (l *Logger) Logger() *zerolog.Logger {
	return &l.l
}

func (l *Logger) AsGlobalLooger() {
	log.Logger = l.l
}

func buildLoggerContextDevMode(appname string) zerolog.Logger {
	zerolog.CallerMarshalFunc = cutCallerFileName

	return zerolog.New(os.Stdout).With().
		Str("app_name", appname).Timestamp().
		Caller().Stack().
		Logger().Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func buildLoggerContextProdMode(appname string) zerolog.Logger {
	zerolog.TimestampFieldName = "t"
	zerolog.LevelFieldName = "l"
	zerolog.MessageFieldName = "m"

	return zerolog.New(os.Stdout).With().
		Str("app_name", appname).Timestamp().
		Logger()
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

func LoggerWithFileOutput(filename string, maxSize, maxAge int, compress bool) LogBuilder {
	if filename == "" {
		panic("必须要指定日志输出的文件")
	}
	if maxSize <= 0 {
		maxSize = 20 // MB
	}
	if maxAge <= 0 {
		maxAge = 3 // days
	}

	dir := path.Dir(filename)
	if !kit.PathExists(dir) {
		os.MkdirAll(dir, 0766)
	}
	return func(l zerolog.Logger) zerolog.Logger {
		return l.Output(&lumberjack.Logger{
			Filename:  filename,
			MaxSize:   maxSize,
			MaxAge:    maxAge,
			LocalTime: true,
			Compress:  compress,
		}).With().Logger()
	}
}
