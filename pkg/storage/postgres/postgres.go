package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/lekht/go-news-aggregator/config"
)

// БД
type Storage struct {
	pool *pgxpool.Pool
}

// Конструктор БД
func New(ctx context.Context, cfg *config.PG) (*Storage, error) {
	var connstr string = fmt.Sprintf("postgres://%s:%s@%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.DB)
	dbpool, err := pgxpool.New(ctx, connstr)
	if err != nil {
		return nil, err
	}
	db := Storage{
		pool: dbpool,
	}
	return &db, nil
}

// Получение n новостей из БД
func (db *Storage) News(n int) ([]*Post, error) {
	rows, err := db.pool.Query(context.Background(), `
		SELECT 
			id,
			title,
			content,
			pubTime,
			link
		FROM items
		ORDER BY pubTime DESC
		LIMIT $1;`,
		n,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var posts []*Post
	for rows.Next() {
		var post Post

		err = rows.Scan(&post.ID, &post.Title, &post.Content, &post.PubTime, &post.Link)
		if err != nil {
			return nil, err
		}

		posts = append(posts, &post)
	}

	if err = rows.Err(); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return posts, nil

}

// Записывает новости в БД
func (db *Storage) AddNews(posts []*Post) error {
	if len(posts) < 1 {
		return errors.New("adding empty slice")
	}
	tx, err := db.pool.Begin(context.Background())
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
	if s.pool != nil {
		s.pool.Close()
	}
}
