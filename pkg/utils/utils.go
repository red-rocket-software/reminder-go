package utils

import (
	log "github.com/sirupsen/logrus"
)

// ConfigureLogger create logrus configuration for whole application
func ConfigureLogger() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		ForceColors:     true,
		DisableColors:   false,
	})

	log.SetLevel(log.DebugLevel)
}
