package env

import (
	"github.com/sirupsen/logrus"
)

var Log logrus.FieldLogger

func init() {
	setupLogging()
}

func setupLogging() {
	logger := logrus.StandardLogger()
	logger.Formatter = &logrus.TextFormatter{DisableTimestamp: true}
	Log = logger
}
