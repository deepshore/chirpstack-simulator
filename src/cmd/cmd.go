package cmd

import (
	"os"

	log "github.com/sirupsen/logrus"
)

var defaultLogLevel = "info"

// setLogLevel configures the standard logger to include more context and
// supports log level filtering via the LOG_LEVEL environment variable.
func setLogLevel() {
	log.SetOutput(os.Stdout)
	log.SetReportCaller(true)

	if level, err := log.ParseLevel(getEnv("LOG_LEVEL", defaultLogLevel)); err != nil {
		log.Warnf("Invalid log level: %v, using default level: %s", err, defaultLogLevel)
		level = log.InfoLevel
	} else {
		log.SetLevel(level)
	}

	log.SetFormatter(&log.JSONFormatter{
		TimestampFormat: "2006-01-02T15:04:05Z07:00",
		PrettyPrint:     false,
	})
}

// getEnv returns the value of the environment variable or a default value.
func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}
