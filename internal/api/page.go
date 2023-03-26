package api

import "github.com/lekht/go-news-aggregator/pkg/storage/postgres"

const newsPerPage int = 40

type Response struct {
	Page  Page
	Posts []*postgres.Post
}

type Page struct {
	TotalPages  int
	CurrentPage int
	NewsPerPage int
}
