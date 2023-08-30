package main

import (
	"dynamic-user-segmentation/config"
	v1 "dynamic-user-segmentation/internal/controller/http/v1"
	postgresql2 "dynamic-user-segmentation/internal/repository"
	"dynamic-user-segmentation/internal/service"
	"dynamic-user-segmentation/pkg/client/db/postgresql"
	"dynamic-user-segmentation/pkg/hash"
	"dynamic-user-segmentation/pkg/httpserver"
	"dynamic-user-segmentation/pkg/logger/sloglogger"
	"dynamic-user-segmentation/pkg/util/errors"
	"github.com/labstack/echo/v4"
	"os"
	"os/signal"
	"syscall"
)

// @swagger 2.0
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
	pg, err := postgresql.NewClient(config.PgUrl(
		cfg.Storage.User,
		cfg.Storage.Password,
		cfg.Storage.Host,
		cfg.Storage.Port,
		cfg.Storage.Database), postgresql.LoadMaxPoolSize(cfg.Storage.MaxPoolSize))

	if err != nil {
		logger.Error("App run  error: %w", errors.Wrap("can't crete client", err))
		os.Exit(1)
	}

	defer pg.Close()

	logger.Info("Initializing repositories...")
	repositories := postgresql2.NewRepositories(pg)

	logger.Info("Initializing services...")
	dependencies := service.ServicesDependencies{
		Repository: repositories,
		SignKey:    cfg.JWT.SignKey,
		TokenTTL:   cfg.JWT.TokenTTL,
		Hash:       hash.NewPasswordHashSHA256(cfg.Hash.Salt),
	}

	services := service.NewServices(dependencies)

	logger.Info("Initializing handlers...")
	handler := echo.New()
	handler.Use(sloglogger.NewMiddleWareLogger(logger))
	v1.NewRouter(handler, services)

	logger.Info("Starting http server...")

	logger.Debug("Server host %s: ", cfg.HTTP.Host)
	logger.Debug("Server port %s: ", cfg.HTTP.Port)

	server := httpserver.NewServer(handler, httpserver.LoadHost(cfg.HTTP.Host, cfg.HTTP.Port))
	server.Start()

	logger.Info("Configuration load successful...")
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case sign := <-interrupt:
		logger.Info("Main signal: %s", sign.String())
	case err = <-server.Notify():
		logger.Error("Main server: %w", err)
	}

	logger.Info("Shutting down...")
	err = server.Shutdown()
	if err != nil {
		logger.Error("Main server shutdown: %w", err)
	}
}
