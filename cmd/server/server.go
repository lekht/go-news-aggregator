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
	dbase, err := db.New(context.Background(), connstr)
	if err != nil {
		log.Fatal(err)
	}

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
	chPost := make(chan []db.Post)
	chErr := make(chan error)
	go func() {
		for _, urls := range config.URLs {
			go parseUrl(urls, chPost, chErr)
		}
		time.Sleep(time.Minute * time.Duration(config.Period))
	}()

	go func() {
		for posts := range chPost {
			err = dbase.AddNews(posts)
			if err != nil {
				log.Println("channel posts adding error: ", err)
			}
		}
	}()

	go func() {
		for err := range chErr {
			log.Println("rss goroutine parsing error: ", err)
		}
	}()

	err = http.ListenAndServe(":8080", api.Router())
	if err != nil {
		panic(err)
	}
}

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
