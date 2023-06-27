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
	logrus.SetFormatter(&MyFormatter{})
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.Level(LogLevel))

	if t == ECLog2StdFile {
		if dir != "" {
			dir = "log"
		}
		logFile := filepath.Join(dir, id+".log")
		logrus.SetOutput(&lumberjack.Logger{
			Filename:   logFile, // 日志文件路径
			MaxSize:    1,       // 单个日志文件的最大大小（MB）
			MaxBackups: 10,      // 保留的最大日志文件个数
			Compress:   false,   // 是否压缩日志文件
			LocalTime:  true,
		})
	}

	return l
}

func (l *Loger) Trace(v ...interface{}) {
	logrus.Trace(v...)
}

func (l *Loger) Debug(v ...interface{}) {
	logrus.Debug(v...)
}

func (l *Loger) Info(v ...interface{}) {
	logrus.Info(v...)
}

func (l *Loger) Warn(v ...interface{}) {
	logrus.Warn(v...)
}

func (l *Loger) Error(v ...interface{}) {
	logrus.Error(v...)
}

func (l *Loger) Fatal(v ...interface{}) {
	logrus.Fatal(v...)
}

func (l *Loger) Tracef(format string, v ...interface{}) {
	logrus.Tracef(format, v...)
}

func (l *Loger) Debugf(format string, v ...interface{}) {
	logrus.Debugf(format, v...)
}

func (l *Loger) Infof(format string, v ...interface{}) {
	logrus.Infof(format, v...)
}

func (l *Loger) Warnf(format string, v ...interface{}) {
	logrus.Warnf(format, v...)
}

func (l *Loger) Errorf(format string, v ...interface{}) {
	logrus.Errorf(format, v...)
}

func (l *Loger) Fatalf(format string, v ...interface{}) {
	logrus.Fatalf(format, v...)
}
