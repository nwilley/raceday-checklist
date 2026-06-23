package main

import (
	"context"
	"log"
	"time"

	"raceday-checklist/api/internal/checklist"
	"raceday-checklist/api/internal/config"
	"raceday-checklist/api/internal/database"
	"raceday-checklist/api/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	checklistRepository := checklist.NewMySQLRepository(db)
	checklistService := checklist.NewService(checklistRepository)

	router := server.NewRouter(checklistService)
	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
