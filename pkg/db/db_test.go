package db

import (
	"context"
	"log"
	"testing"
)

func newDB() *DB {
	db, err := New(context.Background(), "postgres://postgres:password@server.domain/items")
	if err != nil {
		log.Fatalf("db creating error: %v", err)
	}
	return db
}

func TestDB_AddNews(t *testing.T) {
	db := newDB()

	type fields struct {
		db *DB
	}
	type args struct {
		items []Post
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "empty slice",
			fields:  fields{db: db},
			args:    args{[]Post{}},
			wantErr: true,
		},
		{
			name:   "single post",
			fields: fields{db: db},
			args: args{[]Post{
				{
					Title:       "Test1",
					Description: "something interesting",
					PubDate:     "12.12.2012",
					Url:         "test.com",
				},
			}},
			wantErr: false,
		},
		{
			name:   "several posts",
			fields: fields{db: db},
			args: args{[]Post{
				{
					Title:       "Test1",
					Description: "something interesting",
					PubDate:     "12.12.2012",
					Url:         "test.com/1",
				},
				{
					Title:       "Test2",
					Description: "something interesting",
					PubDate:     "12.12.2012",
					Url:         "test.com/2",
				},
				{
					Title:       "Test3",
					Description: "something interesting",
					PubDate:     "12.12.2012",
					Url:         "test.com/3",
				},
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := db.AddNews(tt.args.items); (err != nil) != tt.wantErr {
				t.Errorf("DB.AddNews() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDB_News(t *testing.T) {
	db := newDB()
	posts, err := db.News(5)
	if err != nil {
		t.Errorf("db getting data error: %v", err)
	}
	if len(posts) != 5 {
		t.Error("posts' count error")
	}
	t.Log(posts)
}
