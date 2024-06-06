package logger

import "testing"

func TestLoggerBWarnfase(t *testing.T) {
	defer Logger.CatchPanic()
	Logger.Infof("this is logger %s this is green %s", Green, Reset)
	Logger.Warnf("this is logger")
	Logger.Error("this is logger")
	panic(123)
}
