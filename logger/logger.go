package logger

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"os"
)

var Log = logrus.New()

const rootPath = "logs/"

func InitLogger() {
	appLog := &lumberjack.Logger{
		Filename:   rootPath + "app.log",
		MaxSize:    10, // MB
		MaxBackups: 5,
		MaxAge:     28, // days
		Compress:   true,
	}

	Log.SetOutput(io.MultiWriter(os.Stdout, appLog))
	Log.SetFormatter(defaultFormatter())
	Log.SetLevel(logrus.InfoLevel)

	log.Println("✅  Logger initialized successfully.")
	log.Println()
}

func defaultFormatter() *logrus.TextFormatter {
	return &logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	}
}
