package main

import (
	"log"

	"github.com/lekht/go-news-aggregator/config"
	"github.com/lekht/go-news-aggregator/internal/app"
)

// type config struct {
// 	URLs   []string `json:"rss"`
// 	Period int      `json:"request_period"`
// }

// const connstr string = "postgres://postgres:password@server.domain/items"

func main() {
	// Configuration
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}

	app.Run(cfg)
}

// 	// Инициализация БД
// 	dbase, err := db.New(context.Background(), connstr)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Инициализация API
// 	api := api.New(dbase)

// 	b, err := os.ReadFile("./config/config.json")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	var config config
// 	err = json.Unmarshal(b, &config)
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	// Запуск парсинга rss-потоков новостных сайтов
// 	chPost := make(chan []db.Post)
// 	chErr := make(chan error)
// 	go func() {

// 	}()

// 	// Запуск сервера
// 	err = http.ListenAndServe(":80", api.Router())
// 	if err != nil {
// 		panic(err)
// 	}
// }
