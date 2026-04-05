package logger

import (
    "log"
    "os"
)

type Logger struct {
    info  *log.Logger
    error *log.Logger
    debug *log.Logger
    level string
}

var Global *Logger

func InitLogger(level string) {
    Global = &Logger{
        info:  log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
        error: log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
        debug: log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile),
        level: level,
    }
}

func (l *Logger) Info(v ...interface{}) {
    l.info.Println(v...)
}

func (l *Logger) Error(v ...interface{}) {
    l.error.Println(v...)
}

func (l *Logger) Debug(v ...interface{}) {
    if l.level == "debug" {
        l.debug.Println(v...)
    }
}

func (l *Logger) Infof(format string, v ...interface{}) {
    l.info.Printf(format, v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
    l.error.Printf(format, v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
    if l.level == "debug" {
        l.debug.Printf(format, v...)
    }
}

func (l *Logger) Fatal(v ...interface{}) {
    l.error.Println(v...)
    os.Exit(1)
}
