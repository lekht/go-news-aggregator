package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lekht/go-news-aggregator/pkg/storage/postgres"
)

type storage interface {
	News(n int) ([]*postgres.Post, error)
	NewsPageFilter(page int, filter string, newsPerPage int) ([]*postgres.Post, error)
	NewByID(id int) (*postgres.Post, error)
	// NewsAllInAll() (int, error)
}

// Метод получения записей из БД
func (a *API) news(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}
	c := mux.Vars(r)["n"]
	n, err := strconv.Atoi(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	posts, err := a.db.News(n)
	if err != nil {
		log.Printf("api - news - db getting data error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(posts)
	if err != nil {
		log.Println("api - news - response's data encoding error: ", err)
	}
}

func (a *API) filterAndPage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodOptions {
		return
	}

	var page int

	q := r.URL.Query()
	pageValue := q.Get("page")
	filter := q.Get("filter")
	if pageValue == "" {
		page = 1
	} else {
		i, err := strconv.Atoi(pageValue)
		if err != nil {
			log.Println("api - filterPage - params handler err: ", err)
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		page = i
	}
	log.Printf("api - filterPage - gotten params: %d, %s", page, filter)

	posts, err := a.db.NewsPageFilter(page, filter, newsPerPage)
	if err != nil {
		log.Printf("api - news - db getting data error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	postsNum := len(posts)

	response := &Response{
		Page: Page{
			TotalPages:  postsNum / newsPerPage,
			CurrentPage: page,
			NewsPerPage: newsPerPage,
		},
		Posts: posts,
	}

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("api - filterPage - response's data encoding error: ", err)
	}
}

func (a *API) newByID(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	idValue := q.Get("id")
	id, err := strconv.Atoi(idValue)
	if err != nil {
		log.Printf("api - newByID - convertig id error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	p, err := a.db.NewByID(id)
	if err != nil {
		log.Printf("api - news - db data getting error: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(p)
	if err != nil {
		log.Println("api - news - response's data encoding error: ", err)
	}
}
