package logger

var Logger *mLogger

func init() {
	Logger = NewLogger()
}
