package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/lekht/go-news-aggregator/config"
	"github.com/lekht/go-news-aggregator/pkg/storage/postgres"
)

type storage interface {
	News(n int) ([]*postgres.Post, error)
}

type API struct {
	server          *http.Server
	db              storage
	shutdownTimeout time.Duration

	notify chan error
}

// Регистрация методов в маршрутизаторе
// func (api *API) endpoints() {
// 	api.server.Handler.Name("get_some_last_news").Path("/news/{n}").Methods(http.MethodGet).HandlerFunc(api.itemsHandler)
// 	api.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))

// 	// api.r.Name("get_all_news").Path("/news").Methods(http.MethodGet).HandlerFunc(a.AllPostsHandler)
// 	// api.r.Name("get_news_by_id").Path("/news/full/{id}").Methods(http.MethodGet).HandlerFunc(a.PostHandler)

// 	// api.r.HandleFunc("/news/{n}", api.itemsHandler).Methods(http.MethodGet, http.MethodOptions)
// 	// api.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
// }

// Конструктор API
func New(cfg *config.Server, db storage) *API {
	a := API{
		db:              db,
		notify:          make(chan error, 1),
		shutdownTimeout: 3 * time.Second,
	}
	handler := mux.NewRouter()
	handler.Use(a.accessMiddleware)
	handler.Name("get_last_news").Path("/news/{n}").Methods(http.MethodGet).HandlerFunc(a.itemsHandler)
	handler.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
	// api.endpoints()
	a.server = &http.Server{
		Addr:         cfg.Port,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
		Handler:      handler,
	}
	return &a
}

func (a *API) Start() {
	go func() {
		a.notify <- a.server.ListenAndServe()
		// if err != nil {
		// 	log.Panic("сервер не хочет работать(((")
		// }
		close(a.notify)
		defer fmt.Println("сервер упал")
	}()
	// select {}
}

func (a *API) Notify() <-chan error {
	return a.notify
}

func (a *API) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), a.shutdownTimeout)
	defer cancel()

	return a.server.Shutdown(ctx)
}

func (a *API) writeResponseError(w http.ResponseWriter, err error, code int) {
	w.Header().Add("Code", strconv.Itoa(code))
	// log.WithError(err).Error("api error")
	w.WriteHeader(code)
	_, _ = w.Write([]byte(err.Error()))
}
