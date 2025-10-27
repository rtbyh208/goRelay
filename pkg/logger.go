package pkg

import (
	"fmt"
	"runtime"
	"strings"
	"time"
)

type Logger struct {
	Level int
}

var Log Logger

const (
	NoneLevel = iota
	PanicLevel
	FatalLevel
	ErrorLevel
	WarnLevel
	InfoLevel
	DebugLevel
	TraceLevel
)

func NewLogger() *Logger {
	return &Log
}

func (l Logger) loggerHandle() (string, int) {
	_, file, line, _ := runtime.Caller(2)

	var shortFile string
	shortFileList := strings.Split(file, "/")
	shortFile = shortFileList[len(shortFileList)-1]

	return shortFile, line
}

func newLine(file string, line int, args ...interface{}) string {
	nowTime := time.Now().Format("2006-01-02T15:04:05")
	return fmt.Sprintf("%s %s:%d %s", nowTime, file, line, fmt.Sprint(args...))
}

func (l *Logger) SetLogger(level int) {
	l.Level = level
}

func (l Logger) Panic(args ...interface{}) {
	file, line := l.loggerHandle()
	if l.Level >= PanicLevel {
		fmt.Println("[Panic]", newLine(file, line, args...))
	}
}

func (l Logger) Fatal(args ...interface{}) {
	file, line := l.loggerHandle()
	if l.Level >= FatalLevel {
		fmt.Println("[Fatal] ", newLine(file, line, args...))
	}
}

func (l Logger) Error(args ...interface{}) {
	file, line := l.loggerHandle()
	if l.Level >= ErrorLevel {
		fmt.Println("[Error] ", newLine(file, line, args...))
	}
}

func (l Logger) Warn(args ...interface{}) {
	file, line := l.loggerHandle()
	if l.Level >= WarnLevel {
		fmt.Println("[Warn] ", newLine(file, line, args...))
	}
}

func (l Logger) Info(args ...interface{}) {
	file, line := l.loggerHandle()
	if l.Level >= InfoLevel {
		fmt.Println("[Info] ", newLine(file, line, args...))
	}
}

func (l Logger) Debug(args ...interface{}) {
	file, line := l.loggerHandle()
	if l.Level >= DebugLevel {
		fmt.Println("[Debug] ", newLine(file, line, args...))
	}
}

func (l Logger) Trace(args ...interface{}) {
	file, line := l.loggerHandle()
	if l.Level >= TraceLevel {
		fmt.Println("[Trace] ", newLine(file, line, args...))
	}
}
