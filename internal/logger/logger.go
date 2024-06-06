package logger

import (
	"bytes"
	"fmt"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"runtime"
	"strings"

	"github.com/alecthomas/chroma/quick"
	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
)

const (
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Reset   = "\033[0m"
	BgGreen = "\033[42m" // Background green color
	Emoji   = "🔍"        // Emoji symbol
)

// mLogger 封装了logrus的mLogger
type mLogger struct {
	logger  *log.Logger
	isDebug bool
}

// NewLogger 创建并初始化新的Logger
func NewLogger() *mLogger {
	logger := log.New()
	logger.SetOutput(os.Stdout)

	logger.SetOutput(colorable.NewColorableStdout())
	logger.SetFormatter(&log.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
	})

	logger.SetLevel(log.DebugLevel)

	return &mLogger{logger: logger}
}

func (l *mLogger) SetIsDebug(debug bool) {
	l.isDebug = debug
}

func (l *mLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *mLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *mLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
	l.printCallerCode(log.ErrorLevel)
}

func (l *mLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
	l.printCallerCode(log.ErrorLevel)
}

func (l *mLogger) Debug(args ...interface{}) {
	if !l.isDebug {
		return
	}
	l.logger.Debug(args...)
}

func (l *mLogger) Debugf(format string, args ...interface{}) {
	if !l.isDebug {
		return
	}
	l.logger.Debugf(format, args...)
}

func (l *mLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
	l.printCallerCode(log.WarnLevel)
}

func (l *mLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
	l.printCallerCode(log.WarnLevel)
}

func (l *mLogger) Panic(args ...interface{}) {
	defer l.printCallerCode(log.PanicLevel)
	l.logger.Panic(args...)
}

func (l *mLogger) Panicf(format string, args ...interface{}) {
	defer l.printCallerCode(log.PanicLevel)
	l.logger.Panicf(format, args...)
}

func (l *mLogger) Fatal(args ...interface{}) {
	defer l.printCallerCode(log.FatalLevel)
	l.logger.Fatal(args...)
}

func (l *mLogger) Fatalf(format string, args ...interface{}) {
	defer l.printCallerCode(log.FatalLevel)
	l.logger.Fatalf(format, args...)
}

func getCallerCode(file string, line int) (string, error) {
	src, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, file, src, parser.AllErrors)
	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if err := printer.Fprint(&buf, fset, node); err != nil {
		return "", err
	}

	lines := strings.Split(buf.String(), "\n")
	startLine := max(line-3, 0)
	endLine := min(line+3, len(lines)-1)

	var codeSnippet bytes.Buffer
	for i := startLine; i <= endLine; i++ {
		codeLine := lines[i]
		if i == line-1 {
			// Apply emoji to the specific line
			codeLine = Emoji + codeLine
		}
		codeSnippet.WriteString(codeLine + "\n")
	}

	return codeSnippet.String(), nil
}

func (l *mLogger) printCallerCode(level log.Level) {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return
	}

	code, err := getCallerCode(file, line)
	if err != nil {
		l.logger.Warnf("Unable to retrieve code snippet: %v", err)
		return
	}

	var highlightedCode bytes.Buffer
	if err := quick.Highlight(&highlightedCode, code, "go", "terminal256", "monokai"); err != nil {
		l.logger.Warnf("Unable to highlight code: %v", err)
		return
	}

	l.logger.Log(level, "Code Context:\n\n"+highlightedCode.String())
}

func (l *mLogger) CatchPanic() {
	if r := recover(); r != nil {
		stack := make([]uintptr, 15)
		length := runtime.Callers(3, stack[:])
		stack = stack[:length]
		buf := &bytes.Buffer{}

		fmt.Fprintf(buf, "%sPanic recovered: %v%s\n", Red, r, Reset)
		fmt.Fprintln(buf, "\nStack trace:")
		fmt.Fprintln(buf, "------------------------------------------")
		fmt.Fprintln(buf, "File:Line                | Function:")
		fmt.Fprintln(buf, "------------------------------------------")

		frames := runtime.CallersFrames(stack)
		for {
			frame, more := frames.Next()
			formattedFile := fileWithLine(frame.File, frame.Line)
			formattedFunc := fmt.Sprintf("%.30s", frame.Function)
			fmt.Fprintf(buf, "%-24s | %s\n", formattedFile, formattedFunc)
			if !more {
				break
			}
		}
		fmt.Fprintln(buf, "------------------------------------------")
		l.logger.Error(buf.String())
		l.printCallerCode(log.FatalLevel)
	}
}

func fileWithLine(file string, line int) string {
	shortFile := file
	maxLength := 30
	if len(file) > maxLength {
		shortFile = "..." + file[len(file)-maxLength+3:]
	}
	return fmt.Sprintf("%s:%d", shortFile, line)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
