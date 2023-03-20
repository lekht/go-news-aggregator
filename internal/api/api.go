package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lekht/go-news-aggregator/config"
	"github.com/lekht/go-news-aggregator/pkg/storage/postgres"
)

type API struct {
	r  *mux.Router
	db *postgres.Storage
}

// Регистрация методов в маршрутизаторе
func (a *API) endpoints() {
	a.r.Use(a.accessMiddleware)
	a.r.Name("last_news").Path("/news/{n}").Methods(http.MethodGet).HandlerFunc(a.itemsHandler)
	a.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

func (a *API) Router() *mux.Router {
	return a.r
}

// Конструктор API
func New(cfg *config.Server, db *postgres.Storage) *API {
	a := API{
		db: db,
		r:  mux.NewRouter(),
	}
	a.endpoints()
	return &a
}
