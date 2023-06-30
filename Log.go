package golibs

import (
	"bytes"
	"fmt"
	"path"
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"github.com/sirupsen/logrus"
)

type Loger struct {
	*logrus.Logger
}

const (
	ECLog2StdOut EULogType = iota
	ECLog2StdFile
)

type EULogType int

const (
	LevelTrace EULogLevel = EULogLevel(logrus.TraceLevel)
	LevelDebug EULogLevel = EULogLevel(logrus.DebugLevel)
	LevelInfo  EULogLevel = EULogLevel(logrus.InfoLevel)
	LevelWarn  EULogLevel = EULogLevel(logrus.WarnLevel)
	LevelError EULogLevel = EULogLevel(logrus.ErrorLevel)
	LevelFatal EULogLevel = EULogLevel(logrus.FatalLevel)
)

type EULogLevel int

type MyFormatter struct {
}

func (m *MyFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}
	timestamp := entry.Time.Format("2006/01/02 15:04:05")
	var file string
	var line int
	if entry.Caller != nil {
		file = path.Base(entry.Caller.File)
		line = entry.Caller.Line
	}

	newLog := fmt.Sprintf("[%s] [%s] [%s:%d] %s\n", entry.Level, timestamp, file, line, entry.Message)
	b.WriteString(newLog)
	return b.Bytes(), nil
}

func (l *Loger) Init(t EULogType, LogLevel EULogLevel, dir, id string) *Loger {

	l.Logger = logrus.New()
	l.Logger.SetFormatter(&MyFormatter{})
	l.Logger.SetReportCaller(true)
	l.Logger.SetLevel(logrus.Level(LogLevel))

	if t == ECLog2StdFile {
		if dir != "" {
			dir = "log"
		}
		logFile := filepath.Join(dir, id+".log")

		l.Logger.SetOutput(&lumberjack.Logger{
			Filename:   logFile, // 日志文件路径
			MaxSize:    1,       // 单个日志文件的最大大小（MB）
			MaxBackups: 10,      // 保留的最大日志文件个数
			Compress:   false,   // 是否压缩日志文件
			LocalTime:  true,
		})
	}

	return l
}
