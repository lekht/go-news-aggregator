package api

import (
	"context"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/lekht/go-news-aggregator/internal/logger"
)

type ContextKey string

const ContextRequestKey ContextKey = "request_id"

func (a *API) accessMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}

func (a *API) requestIdMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Query().Get("request_id")

		if id == "" {
			id = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), ContextRequestKey, id)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (a *API) logRequestMiddlware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			reqInf := &logger.HTTPReqInfo{
				Time:         time.Now().Format("2006-01-02 15:04:05.000000 MST"),
				Method:       r.Method,
				ResponseCode: w.Header().Get("Code"),
				IPadress:     logger.RequestGetRemoteAddress(r),
				RequestID:    r.Context().Value(ContextRequestKey),
			}
			reqInf.Info()
		}()

		next.ServeHTTP(w, r)
	})
}
