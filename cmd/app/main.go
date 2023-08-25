package main

import (
	"github.com/coven-discord-bot/config"
	"github.com/coven-discord-bot/internal/app"
	"log"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)
}
