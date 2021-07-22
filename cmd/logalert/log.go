package main

import (
	logs "github.com/sirupsen/logrus"
	"io"
	"os"
)

func convertLogLevel(level string) logs.Level {
	switch level {
	case "debug":
		return logs.DebugLevel
	case "info":
		return logs.InfoLevel
	case "warn":
		return logs.WarnLevel
	case "error":
		return logs.ErrorLevel
	}
	return logs.InfoLevel
}

func initLog() (err error) {
	logPath := appConfig["logs"]["path"]
	logAdapter := appConfig["logs"]["adapter"]
	logLevel := appConfig["logs"]["level"]

	if len(logPath) == 0 {
		logAdapter = "console"
	}

	logs.SetFormatter(&logs.JSONFormatter{
		FieldMap: logs.FieldMap{
			logs.FieldKeyMsg: "message",
		},
	})
	if logAdapter == "file" {
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			logs.Error("open log file err:", err)
		}
		writer := io.Writer(file)
		logs.SetOutput(writer)
	} else {
		writer := os.Stdout
		logs.SetOutput(writer)
	}

	logs.SetLevel(convertLogLevel(logLevel))
	logs.Info("init logger success")
	return
}
