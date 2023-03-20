package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

// Метод получения записей из БД
func (api *API) itemsHandler(w http.ResponseWriter, r *http.Request) {
	// w.Header().Set("Content-Type", "application/json")
	// w.Header().Set("Access-Control-Allow-Origin", "*")
	nn := mux.Vars(r)["n"]
	n, err := strconv.Atoi(nn)
	if err != nil {
		api.writeResponseError(w, err, http.StatusBadRequest)
		return
	}

	news, err := api.db.News(n)
	if err != nil {
		api.writeResponseError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Add("Code", strconv.Itoa(http.StatusOK))
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	_ = json.NewEncoder(w).Encode(news)
}
