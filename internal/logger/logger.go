package logger

//go:generate moq -out logger_mock.go . Logger
type Logger interface {
	Fatal(v ...interface{})
	Fatalf(format string, v ...interface{})
	Info(v ...interface{})
	Infof(format string, v ...interface{})
	EnableDebug()
	Debug(v ...interface{})
	Debugf(format string, v ...interface{})
}
