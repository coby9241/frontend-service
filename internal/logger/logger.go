package logger

import (
	"os"
	"sync"

	log "github.com/sirupsen/logrus"
)

// Logger is an alias for logrus logger
type Logger = log.Logger

var instance *Logger
var once sync.Once

// GetInstance returns a Logger pointer to log
func GetInstance() *Logger {
	once.Do(func() {
		instance = log.New()
		// Log as text with options
		instance.SetFormatter(&log.TextFormatter{
			TimestampFormat: "2006-01-02 15:04:05",
			FullTimestamp:   true,
		})
		// Output to stdout instead of the default stderr
		// Can be any io.Writer, see below for File example
		instance.SetOutput(os.Stdout)
		// Set log level to info by default
		instance.SetLevel(log.InfoLevel)
	})

	return instance
}
