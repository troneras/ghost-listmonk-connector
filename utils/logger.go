package utils

import (
	"log"
	"os"
)

type Logger struct {
	*log.Logger
}

var (
	DebugLogger *Logger
	InfoLogger  *Logger
	ErrorLogger *Logger
)

func init() {
	DebugLogger = &Logger{log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)}
	InfoLogger = &Logger{log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)}
	ErrorLogger = &Logger{log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)}
}



func (l *Logger) Infof(format string, v ...interface{}) {
	l.Printf(format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.Printf(format, v...)
}
