package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lekht/go-news-aggregator/config"
	"github.com/lekht/go-news-aggregator/pkg/storage/postgres"
)

type storage interface {
	News(n int) ([]*postgres.Post, error)
}

// type API struct {
// 	server          *http.Server
// 	db              storage
// 	shutdownTimeout time.Duration

//		notify chan error
//	}
type API struct {
	r  *mux.Router
	db storage
}

// Регистрация методов в маршрутизаторе
func (a *API) endpoints() {
	// a.r.Use(a.accessMiddleware)
	// a.r.Name("get_last_news").Path("/news/{n}").Methods(http.MethodGet).HandlerFunc(a.itemsHandler)
	// a.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))

	a.r.Use(a.accessMiddleware)
	a.r.HandleFunc("/news/{n}", a.itemsHandler).Methods(http.MethodGet, http.MethodOptions)
	a.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

func (a *API) Router() *mux.Router {
	return a.r
}

// Конструктор API
func New(cfg *config.Server, db storage) *API {
	// a := API{
	// 	db:              db,
	// 	notify:          make(chan error, 1),
	// 	shutdownTimeout: 3 * time.Second,
	// }
	a := API{
		db: db,
		r:  mux.NewRouter(),
	}
	a.endpoints()
	return &a
}

// func (a *API) Start() {
// 	go func() {
// 		a.notify <- a.server.ListenAndServe()
// 		// if err != nil {
// 		// 	log.Panic("сервер не хочет работать(((")
// 		// }
// 		close(a.notify)
// 		defer fmt.Println("сервер упал")
// 	}()
// 	// select {}
// }

// func (a *API) Notify() <-chan error {
// 	return a.notify
// }

// func (a *API) Shutdown() error {
// 	ctx, cancel := context.WithTimeout(context.Background(), a.shutdownTimeout)
// 	defer cancel()

// 	return a.server.Shutdown(ctx)
// }

// func (a *API) writeResponseError(w http.ResponseWriter, err error, code int) {
// 	w.Header().Add("Code", strconv.Itoa(code))
// 	// log.WithError(err).Error("api error")
// 	w.WriteHeader(code)
// 	_, _ = w.Write([]byte(err.Error()))
// }
