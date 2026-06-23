package main

import (
	"context"
	"log"
	"time"

	"raceday-checklist/api/internal/config"
	"raceday-checklist/api/internal/database"
	"raceday-checklist/api/internal/seed"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("load config: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, cfg.Database)
	if err != nil {
		log.Fatalf("connect database: %v", err)
	}
	defer db.Close()

	if err := seed.Run(ctx, db); err != nil {
		log.Fatalf("seed database: %v", err)
	}

	log.Print("database seed complete")
}
