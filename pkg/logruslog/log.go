package logruslog

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func New(logLevel string) (*logrus.Logger, error) {
	lvl, err := logrus.ParseLevel(logLevel)
	if err != nil {
		return nil, fmt.Errorf("error parsing log level: %w", err)
	}

	l := logrus.New()

	l.SetFormatter(&logrus.TextFormatter{DisableColors: false})
	l.SetLevel(lvl)

	return l, nil
}
