package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

func init() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	logrus.SetLevel(logrus.InfoLevel)
	if l, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL")); err == nil {
		logrus.SetLevel(l)
	}
}

func New(scope string) *logrus.Entry {
	return logrus.WithField("scope", scope)
}
