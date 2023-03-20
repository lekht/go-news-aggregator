package app

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/lekht/go-news-aggregator/config"
	"github.com/lekht/go-news-aggregator/internal/api"
	"github.com/lekht/go-news-aggregator/pkg/storage/postgres"
)

func Run(cfg *config.Config) {
	ctx := context.Background()
	pg, err := postgres.New(ctx, &cfg.PG)
	if err != nil {
		log.Fatalf("app - Run - postgres.New: %v", err)
	}
	defer pg.Close()
	// rss.New(ctx, &cfg.RSS, pg)
	api := api.New(&cfg.Server, pg)
	api.Start()
	// router := api.Router()
	// httpServer := server.New(router, server.Port(cfg.Server.Port))

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-interrupt:
		log.Println(fmt.Errorf("app - Run - signal: " + s.String()))
	case err = <-api.Notify():
		log.Println(fmt.Errorf("app - Run - server.Notify: %w", err))
	}

	err = api.Shutdown()
	if err != nil {
		log.Println(fmt.Errorf("app - Run - server.Shutdown: %w", err))
	}
}