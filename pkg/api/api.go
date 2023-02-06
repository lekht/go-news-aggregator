package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lekht/go-news-aggregator/pkg/db"
)

type API struct {
	r  *mux.Router
	db *db.DB
}

func (api *API) Router() *mux.Router {
	return api.r
}

func (api *API) endpoints() {
	api.r.HandleFunc("/news/{count}", api.itemsHandler).Methods(http.MethodGet, http.MethodOptions)
	api.r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./webapp"))))
}

func New(db *db.DB) *API {
	api := API{}
	api.db = db
	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

func (api *API) itemsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method == http.MethodOptions {
		return
	}
	c := mux.Vars(r)["count"]
	n, err := strconv.Atoi(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	posts, err := api.db.News(n)
	if err != nil {
		log.Printf("db getting data error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(posts)
}
