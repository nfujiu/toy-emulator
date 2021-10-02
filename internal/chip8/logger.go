package chip8

import (
		"log"
		"os"
)

type Logger interface {
		Debug(string, ...interface{})
		Info(string, ...interface{})
}

type CustomLogger struct {
		Logger *log.Logger
}

var logger Logger = &CustomLogger {
		Logger: log.New(os.Stderr, "", log.LstdFlags),
}

func (l *CustomLogger) Debug(format string, v ...interface{}) {
		l.Logger.Printf("[DEBUG]  " + format, v...)
}

func (l *CustomLogger) Info(format string, v ...interface{}) {
		l.Logger.Printf("[INFO]  " + format, v...)
}