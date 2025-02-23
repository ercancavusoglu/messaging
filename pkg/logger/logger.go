package logger

import (
	"log"
	"os"
)

var (
	infoLogger  *log.Logger
	errorLogger *log.Logger
)

func init() {
	infoLogger = log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}

func Info(msg string, keyvals ...interface{}) {
	infoLogger.Printf("%s %v", msg, keyvals)
}

func Error(msg string, keyvals ...interface{}) {
	errorLogger.Printf("%s %v", msg, keyvals)
}
