package logger

import (
	"github.com/sirupsen/logrus"
)

type Logger struct {
	logger *logrus.Logger
}

func (l *Logger) EnableDebug() {
	l.logger.SetLevel(logrus.DebugLevel)
}

func (l *Logger) Debug(v ...interface{}) {
	l.logger.Debug(v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.logger.Debugf(format, v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.logger.Warn(v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.logger.Warnf(format, v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.logger.Info(v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.logger.Infof(format, v...)
}

func New() *Logger {
	return &Logger{logger: logrus.New()}
}
