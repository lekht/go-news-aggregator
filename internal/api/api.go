package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lekht/go-news-aggregator/config"
)

type API struct {
	r  *mux.Router
	db storage
}

// Регистрация методов в маршрутизаторе
func (a *API) endpoints() {
	a.r.Use(a.accessMiddleware, a.requestIdMiddlware, a.logRequestMiddlware)
	a.r.Name("new_by_id").Path("/news/").Methods(http.MethodGet).HandlerFunc(a.newByID)
	a.r.Name("last_news").Path("/news/{n}").Methods(http.MethodGet).HandlerFunc(a.news)
	a.r.Name("page_and_filter").Path("/news").Methods(http.MethodGet).HandlerFunc(a.filterAndPage)
	// a.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

func (a *API) Router() *mux.Router {
	return a.r
}

// Конструктор API
func New(cfg *config.Server, db storage) *API {
	a := API{
		db: db,
		r:  mux.NewRouter(),
	}
	a.endpoints()
	return &a
}
