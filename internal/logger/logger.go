package logger

import "github.com/sirupsen/logrus"

// SetupLogger configures logrus to include timestamp and use text formatter.
func SetupLogger() {
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}
