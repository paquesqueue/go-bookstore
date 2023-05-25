package common

import (
	"io"
	"os"

	log "github.com/sirupsen/logrus"
)

type Log interface {
	Info(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})
}

func InitLog() *log.Logger {
	file, err := os.OpenFile("logs.text", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Error("Error Create Log File Failed")
	}

	logs := &log.Logger{
		Out:       io.MultiWriter(file, os.Stdout),
		Formatter: new(log.JSONFormatter),
		Hooks:     make(log.LevelHooks),
		Level:     log.InfoLevel,
		ExitFunc:  os.Exit,
	}
	return logs
}

func InitRequestLog() *log.Logger {
	file, err := os.OpenFile("request-logs.text", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Error("Error Create Request Log File Failed")
	}

	logs := &log.Logger{
		Out:       io.MultiWriter(file, os.Stdout),
		Formatter: new(log.JSONFormatter),
		Hooks:     make(log.LevelHooks),
		Level:     log.InfoLevel,
		ExitFunc:  os.Exit,
	}
	return logs
}
