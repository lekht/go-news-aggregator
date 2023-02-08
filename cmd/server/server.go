package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/lekht/go-news-aggregator/pkg/api"
	"github.com/lekht/go-news-aggregator/pkg/db"
	"github.com/lekht/go-news-aggregator/pkg/rss"
)

type config struct {
	URLs   []string `json:"rss"`
	Period int      `json:"request_period"`
}

const connstr string = "postgres://postgres:password@server.domain/items"

func main() {
	// Инициализация БД
	dbase, err := db.New(context.Background(), connstr)
	if err != nil {
		log.Fatal(err)
	}
	// Инициализация API
	api := api.New(dbase)

	b, err := os.ReadFile("./config.json")
	if err != nil {
		log.Fatal(err)
	}
	var config config
	err = json.Unmarshal(b, &config)
	if err != nil {
		log.Fatal(err)
	}

	// Запуск парсинга rss-потоков новостных сайтов
	chPost := make(chan []db.Post)
	chErr := make(chan error)
	go func() {
		for _, urls := range config.URLs {
			go parseUrl(urls, chPost, chErr)
		}
		time.Sleep(time.Minute * time.Duration(config.Period))
	}()

	// Записаь постов из канала в бд
	go func() {
		for posts := range chPost {
			err = dbase.AddNews(posts)
			if err != nil {
				log.Println("channel posts adding error: ", err)
			}
		}
	}()

	// Обработка ошибок
	go func() {
		for err := range chErr {
			log.Println("rss goroutine parsing error: ", err)
		}
	}()

	// Запуск сервера
	err = http.ListenAndServe("localhost:80", api.Router())
	if err != nil {
		panic(err)
	}
}

// Чтение rss-потока и отправка раскодированных постов и ошибок в каналы.
func parseUrl(url string, posts chan<- []db.Post, errs chan<- error) {
	for {
		feeds, err := rss.Parse(url)
		if err != nil {
			errs <- err
			continue
		}
		posts <- feeds
		return
	}
}
