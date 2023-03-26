package main

import (
	"log"

	"github.com/lekht/go-news-aggregator/config"
	"github.com/lekht/go-news-aggregator/internal/app"
)

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	app.Run(cfg)
}
