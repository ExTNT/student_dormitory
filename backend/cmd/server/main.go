package main

import (
	"fmt"
	"log"

	"dormitory/backend/config"
	"dormitory/backend/internal/db"
	"dormitory/backend/internal/handler"
	"dormitory/backend/internal/service"
	"go.uber.org/zap"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("init logger: %v", err)
	}
	defer logger.Sync()

	database, err := db.New(cfg.Database)
	if err != nil {
		logger.Fatal("connect database", zap.Error(err))
	}
	defer database.Close()

	svc := service.New(database, cfg)
	router := handler.NewRouter(svc, cfg, logger)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Info("server_starting", zap.String("addr", addr))
	if err := router.Run(addr); err != nil {
		logger.Fatal("server stopped", zap.Error(err))
	}
}
