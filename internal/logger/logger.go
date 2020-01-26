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
		instance = &Logger{}
		// Log as JSON instead of the default ASCII formatter.
		instance.SetFormatter(&log.JSONFormatter{})
		// Output to stdout instead of the default stderr
		// Can be any io.Writer, see below for File example
		instance.SetOutput(os.Stdout)
		// Only log the warning severity or above.
		instance.SetLevel(log.InfoLevel)
	})

	return instance
}
