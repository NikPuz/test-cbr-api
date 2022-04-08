package main

import (
	"fmt"
	"test-service-a/internal/cbr_api"
	"time"

	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("panic recovery.", r)
		}
	}()

	logger := initZapLogger()

	initViperConfigger(logger)

	// Назначение данных локации
	location, err := time.LoadLocation("Europe/Moscow")
	if err != nil {
		logger.Error("Error time.LoadLocation Europe/Moscow", zap.Error(err))
	}

	// // Назначение данных ежедневных запросов
	et, err := time.ParseInLocation("15:04:05.999", viper.GetString("REQUEST_EVERYDAY"), location)
	if err != nil {
		logger.Error("Error ParseInLocation REQUEST_EVERYDAY", zap.Error(err))
	}

	// Назначение данных запросов с задержкой
	// delay, err := time.ParseDuration(viper.GetString("REQUEST_DELAY"))
	// if err != nil {
	// 	logger.Error("err ParseDuration REQUEST_DELAY", zap.Error(err))
	// }

	// Старт работы
	if viper.GetString("REQUEST_EVERYDAY") != "" {
		cbr_api.NatsGetEveryday(logger, et.In(location), viper.GetString("CBR_URL"))
	}

	// if viper.GetString("REQUEST_DELAY") != "" {
	// 	cbr_api.NatsGetWithDelay(logger, delay, viper.GetString("CBR_URL"))
	// }
	fmt.Scanf(" ")
}

////////////////////
////////////////////
////////////////////

func initViperConfigger(logger *zap.Logger) {
	viper.SetConfigName("app")
	viper.AddConfigPath("config/.")
	err := viper.ReadInConfig()
	if err != nil {
		logger.Error("failed read in config", zap.Error(err))
		return
	}
}

func initZapLogger() *zap.Logger {
	logger := zap.NewExample()
	return logger
}
