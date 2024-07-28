package main

import (
	"fmt"
	"messaggio/internal/repository"
	"messaggio/internal/router"
	"messaggio/internal/service/kafka"
	"messaggio/internal/service/messages"
	"path"
	"runtime"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	setLogger()
	err := godotenv.Load()
	if err != nil {
		logrus.Fatalf(".env file not found.")
	}
	repository := repository.New()
	messageService := messages.New(repository)
	kafka := kafka.New(messageService)
	messageService.SetKafka(kafka)
	router := router.New(messageService)
	router.StartServer()
}

func setLogger() {
	logrus.SetReportCaller(true)
	logrus.SetLevel(logrus.DebugLevel)
	formatter := &logrus.TextFormatter{
		CallerPrettyfier: func(f *runtime.Frame) (function string, file string) {
			filename := path.Base(f.File)
			return fmt.Sprintf("%s()", f.Function), fmt.Sprintf("%s:%d", filename, f.Line)
		},
		DisableColors: false,
		FullTimestamp: true,
	}
	logrus.SetFormatter(formatter)
}
