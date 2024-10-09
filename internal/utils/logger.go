package utils

import (
	"shanraq.org/config"

	"github.com/sirupsen/logrus"
)

func NewLogger(cfg *config.Config) *logrus.Logger {
	logger := logrus.New()

	// Setting the logging level
	logger.SetLevel(logrus.InfoLevel)

	// Exit form
	logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})

	return logger
}