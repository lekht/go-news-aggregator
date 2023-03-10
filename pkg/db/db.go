package db

import (
	"context"
	"errors"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Пост из rss потока
type Post struct {
	ID          int
	Title       string
	Description string
	PubDate     string
	Url         string
}

// БД
type DB struct {
	pool *pgxpool.Pool
}

// Конструктор БД
func New(ctx context.Context, connstr string) (*DB, error) {
	dbpool, err := pgxpool.New(ctx, connstr)
	if err != nil {
		return nil, err
	}
	db := DB{
		pool: dbpool,
	}
	return &db, nil
}

// Получение n новостей из БД
func (db *DB) News(n int) ([]Post, error) {
	rows, err := db.pool.Query(context.Background(), `SELECT * FROM items ORDER BY id DESC LIMIT $1`, n)
	if err != nil {
		return nil, err
	}
	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[Post])
	if err != nil {
		log.Printf("CollectRows error: %v", err)
		return nil, err
	}
	return items, rows.Err()

}

// Записывает новости в БД
func (db *DB) AddNews(items []Post) error {
	if len(items) < 1 {
		return errors.New("adding empty slice")
	}
	tx, err := db.pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	for _, i := range items {
		_, err := tx.Exec(context.Background(), `INSERT INTO items (title, description, PubDate, url)
			VALUES ($1,$2,$3,$4) 
			RETURNING id`,
			&i.Title,
			&i.Description,
			&i.PubDate,
			&i.Url,
		)
		if err != nil {
			return err
		}
	}
	tx.Commit(context.Background())
	return nil
}
