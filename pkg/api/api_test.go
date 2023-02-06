package api

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/lekht/go-news-aggregator/pkg/db"
)

var connstr string = "postgres://postgres:password@server.domain/items"

func TestAPI_itemsHandler(t *testing.T) {
	dbase, _ := db.New(context.Background(), connstr)
	api := New(dbase)
	req := httptest.NewRequest(http.MethodGet, "/news/5", nil)
	rr := httptest.NewRecorder()
	api.r.ServeHTTP(rr, req)
	if !(rr.Code == http.StatusOK) {
		t.Errorf("код неверен, получили %d, а хотели %d", rr.Code, http.StatusOK)
	}
	b, err := io.ReadAll(rr.Body)
	if err != nil {
		t.Fatalf("не удалось прочитать ответ сервера: %v", err)
	}
	var data []db.Post
	err = json.Unmarshal(b, &data)
	if err != nil {
		t.Fatalf("не удалось раскодировать сообщение сервера: %v", err)
	}
	const wantLen = 5
	if len(data) != wantLen {
		t.Fatalf("получено %d записей, когда ожидалось %d", len(data), wantLen)
	}
	t.Log(string(b))
}
