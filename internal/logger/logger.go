package logger

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

type Logger interface {
}

type logger struct {
	*logrus.Logger
}

func New(logLevel string) (Logger, error) {
	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, fmt.Errorf("error parsing log level: %w", err)
	}

	l := logrus.New()

	l.SetFormatter(&logrus.TextFormatter{DisableColors: false})
	l.SetLevel(lvl)

	return logger{l}, nil
}
