package logging

import (
	"log"

	"go.uber.org/zap"
)
var Logger *zap.SugaredLogger

// Функция создания логгера
func CreateLogger() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("creating logger: %s", err.Error())
	}
	Logger = logger.Sugar()
}
