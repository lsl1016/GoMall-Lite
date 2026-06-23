package main

import (
	"log"

	"gomall-lite-api/config"
	"gomall-lite-api/internal/model"
	"gomall-lite-api/internal/router"
)

func main() {
	cfg := config.Load()

	if err := model.InitDB(cfg); err != nil {
		log.Fatalf("init database failed: %v", err)
	}

	r := router.SetupRouter(cfg)
	log.Printf("GoMall Lite API running on :%s", cfg.Port)
	if err := r.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
