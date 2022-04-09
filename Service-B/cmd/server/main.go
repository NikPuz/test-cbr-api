package main

import (
	"encoding/json"
	"net/http"
	"sync"
	"time"

	"test-service-b/internal/cbr"
	"test-service-b/internal/entity"
	"test-service-b/internal/middleware"

	"github.com/gorilla/mux"
	"github.com/nats-io/nats.go"
	"github.com/spf13/viper"
	"go.uber.org/zap"
)

func main() {

	logger := initZapLogger()

	initViperConfigger(logger)

	middleware := middleware.NewMiddleware(logger)

	router := mux.NewRouter()

	var courseValues []entity.ValCurs

	var mutex sync.Mutex

	nc, err := nats.Connect(viper.GetString("NATS_CONNECT"))
	if err != nil {
		logger.Error("Error Connect to nats", zap.Error(err))
	}

	nc.Subscribe("ValCurs", func(msg *nats.Msg) {
		var valCurs entity.ValCurs
		if err := json.Unmarshal(msg.Data, &valCurs); err != nil {
			logger.Error("Error json.Unmarshal msg.Data", zap.Error(err))
		}
		mutex.Lock()
		courseValues = append(courseValues, valCurs)
		mutex.Unlock()
	})
	nc.Flush()

	cbrRepository := cbr.NewCbrRepository(&mutex, logger, &courseValues)
	cbrService := cbr.NewCbrService(logger, cbrRepository)
	cbr.RegisterCbrHandlers(router, logger, cbrService)

	router.Use(middleware.RequestLogger)
	router.Use(middleware.PanicRecovery)

	configServer(router).ListenAndServe()
}

func configServer(router http.Handler) *http.Server {
	return &http.Server{
		Handler:        router,
		Addr:           ":" + viper.GetString("APP_PORT"),
		WriteTimeout:   5 * time.Second,
		ReadTimeout:    5 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
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
