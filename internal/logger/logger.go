package logger

import (
	"bytes"
	"fmt"
	"os"
	"runtime"

	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
)

const (
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Reset  = "\033[0m"
)

// mLogger 封装了logrus的mLogger
type mLogger struct {
	logger *log.Logger
}

// NewLogger 创建并初始化新的Logger
func NewLogger() *mLogger {
	logger := log.New()
	logger.SetOutput(os.Stdout)

	// 设置输出为彩色
	logger.SetOutput(colorable.NewColorableStdout())

	// 设置日志格式为文本格式并强制使用颜色
	logger.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05", // 自定义时间格式
	})

	// 设置日志级别
	logger.SetLevel(log.DebugLevel)

	return &mLogger{logger: logger}
}

// Info 记录普通信息日志
func (l *mLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

// Infof 记录普通信息日志 (格式化)
func (l *mLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

// Error 记录错误日志
func (l *mLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

// Errorf 记录错误日志 (格式化)
func (l *mLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

// Debug 记录调试信息日志
func (l *mLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

// Debugf 记录调试信息日志 (格式化)
func (l *mLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

// Warn 记录警告日志
func (l *mLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

// Warnf 记录警告日志 (格式化)
func (l *mLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

// Panic 记录严重错误日志并引发 panic
func (l *mLogger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

// Panicf 记录严重错误日志并引发 panic (格式化)
func (l *mLogger) Panicf(format string, args ...interface{}) {
	l.logger.Panicf(format, args...)
}

// Fatal 记录致命错误日志并终止程序
func (l *mLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

// Fatalf 记录致命错误日志并终止程序 (格式化)
func (l *mLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

// CatchPanic 捕获 panic 并使用 Logger 友好输出
func (l *mLogger) CatchPanic() {
	if r := recover(); r != nil {
		stack := make([]uintptr, 15)
		length := runtime.Callers(3, stack[:])
		stack = stack[:length]
		buf := &bytes.Buffer{}

		// Print the panic message in red color
		fmt.Fprintf(buf, "%sPanic recovered: %v%s\n", Red, r, Reset)
		fmt.Fprintln(buf, "\nStack trace:")
		fmt.Fprintln(buf, "------------------------------------------")
		fmt.Fprintln(buf, "File:Line                | Function:")
		fmt.Fprintln(buf, "------------------------------------------")

		frames := runtime.CallersFrames(stack)
		for {
			frame, more := frames.Next()
			formattedFile := fileWithLine(frame.File, frame.Line)
			formattedFunc := fmt.Sprintf("%.30s", frame.Function) // Adjust to desired length
			fmt.Fprintf(buf, "%-24s | %s\n", formattedFile, formattedFunc)
			if !more {
				break
			}
		}
		fmt.Fprintln(buf, "------------------------------------------")
		l.logger.Error(buf.String())
	}
}

// fileWithLine formats the file path with line number ensuring proper width
func fileWithLine(file string, line int) string {
	shortFile := file
	// Use a max length for file paths for consistency
	maxLength := 30
	if len(file) > maxLength {
		shortFile = "..." + file[len(file)-maxLength+3:]
	}
	return fmt.Sprintf("%s:%d", shortFile, line)
}
