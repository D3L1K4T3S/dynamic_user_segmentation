package main

import (
	"dynamic-user-segmentation/config"
	"dynamic-user-segmentation/pkg/client/db"
	"dynamic-user-segmentation/pkg/logger/sloglogger"
	"dynamic-user-segmentation/pkg/util/erros"
	"os"
)

// @title Dynamic User Segmentation Swagger API
// @version 0.0.1
// @description This is service for user analytics and slugs
// @termOfService http://swagger.io/terms/

// @contact.name Egor Zhelagin
// @contact.email zhelagin.egor@yandex.ru

// @license.name MIT

// @host localhost:8080
// @BasePath /

const configPath = "config/config.yaml"

func main() {

	// Load config
	cfg := config.MustLoadConfig(configPath)

	// Setup logger
	logger := sloglogger.NewLogger(sloglogger.SetLevel(cfg.Log.Level))

	logger.Info("Initializing client postgreSQl")
	_, err := db.NewClient(config.PgUrl(
		cfg.Storage.User,
		cfg.Storage.Password,
		cfg.Storage.Host,
		cfg.Storage.Port,
		cfg.Storage.Database), db.LoadMaxPoolSize(cfg.Storage.MaxPoolSize))

	if err != nil {
		logger.Error("App run  error: %w", erros.Wrap("can't crete client", err))
		os.Exit(1)
	}

	logger.Info("Initializing repositories")

}
