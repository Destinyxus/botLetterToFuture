package logger

import (
	commander "github.com/Destinyxus/botLetterToFuture/internal/bot_commander"
	"github.com/sirupsen/logrus"
)

type logger struct {
	logrus *logrus.Logger
}

func New(l *logrus.Logger) commander.Logger {
	return logger{logrus: l}
}

func (l logger) Debug(args ...interface{}) {
	l.logrus.Debug(args...)
}

func (l logger) Debugf(format string, args ...interface{}) {
	l.logrus.Debugf(format, args...)
}

func (l logger) Infof(format string, args ...interface{}) {
	l.logrus.Infof(format, args...)
}
