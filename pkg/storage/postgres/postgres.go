package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lekht/go-news-aggregator/config"
)

// БД
type Storage struct {
	Pool *pgxpool.Pool
}

// Конструктор БД
func New(ctx context.Context, cfg *config.PG) (*Storage, error) {
	var connstr string = fmt.Sprintf("postgres://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DB)
	dbpool, err := pgxpool.New(ctx, connstr)
	if err != nil {
		return nil, err
	}
	db := Storage{
		Pool: dbpool,
	}
	return &db, nil
}

// Получение n новостей из БД
func (db *Storage) News(n int) ([]*Post, error) {
	rows, err := db.Pool.Query(context.Background(), `SELECT * FROM items ORDER BY id DESC LIMIT $1`, n)
	if err != nil {
		return nil, err
	}
	items, err := pgx.CollectRows(rows, pgx.RowToStructByName[*Post])
	if err != nil {
		log.Printf("CollectRows error: %v", err)
		return nil, err
	}
	return items, rows.Err()

}

// Записывает новости в БД
func (db *Storage) AddNews(posts []*Post) error {
	if len(posts) < 1 {
		return errors.New("adding empty slice")
	}
	tx, err := db.Pool.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())
	for _, p := range posts {
		_, err := tx.Exec(context.Background(), `INSERT INTO items (title, content, pubTime, link)
			VALUES ($1,$2,$3,$4) 
			RETURNING id`,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return err
		}
	}
	tx.Commit(context.Background())
	return nil
}

func (s *Storage) Close() {
	if s.Pool != nil {
		s.Pool.Close()
	}
}
