package main

import (
	"os"

	"gomall-lite-api/config"
	"gomall-lite-api/internal/logger"
	"gomall-lite-api/internal/model"
	"gomall-lite-api/internal/router"
)

func main() {
	logger.Init()
	cfg := config.Load()

	logger.Default().Info("starting gomall-lite api", "port", cfg.Port)
	if err := model.InitDB(cfg); err != nil {
		logger.Default().Error("init database failed", "error", err)
		os.Exit(1)
	}

	r := router.SetupRouter(cfg)
	logger.Default().Info("gomall-lite api running", "port", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		logger.Default().Error("server failed", "error", err)
		os.Exit(1)
	}
}
