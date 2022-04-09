package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
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

	go startNatsGetEveryday(logger)

	go startNatsGetWithDelay(logger)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
}

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

func startNatsGetEveryday(logger *zap.Logger) {
	if viper.GetString("REQUEST_EVERYDAY") != "" {

		location, _ := time.LoadLocation("Europe/Moscow")

		et, err := time.ParseInLocation("15:04:05.999", viper.GetString("REQUEST_EVERYDAY"), location)
		if err != nil {
			logger.Error("Error ParseInLocation REQUEST_EVERYDAY", zap.Error(err))
		}

		// Начало работы
		cbr_api.NatsGetEveryday(logger, et.In(location), viper.GetString("CBR_URL"))
	} else {
		logger.Info("NatsGetEveryday is off")
	}
}

func startNatsGetWithDelay(logger *zap.Logger) {
	if viper.GetString("REQUEST_DELAY") != "" {

		delay, err := time.ParseDuration(viper.GetString("REQUEST_DELAY"))
		if err != nil {
			logger.Error("err ParseDuration REQUEST_DELAY", zap.Error(err))
		}

		// Начало работы
		cbr_api.NatsGetWithDelay(logger, delay, viper.GetString("CBR_URL"))
	} else {
		logger.Info("NatsGetWithDelay is off")
	}
}
