package rss

import (
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	strip "github.com/grokify/html-strip-tags-go"
	"github.com/lekht/go-news-aggregator/config"
	"github.com/lekht/go-news-aggregator/pkg/storage/postgres"
)

type storage interface {
	AddNews([]*postgres.Post) error
}

type Parser struct {
	db     storage
	links  []string
	period int64

	notify chan error

	errorCh chan error
	postCh  chan []*postgres.Post
}

func New(ctx context.Context, cfg *config.RSS, db storage) *Parser {
	p := &Parser{
		db:     db,
		links:  cfg.URLs,
		period: cfg.Period,

		notify: make(chan error, 1),

		errorCh: make(chan error),
		postCh:  make(chan []*postgres.Post),
	}
	// p.start(ctx)

	return p
}

// start worker
func (p *Parser) Start(ctx context.Context) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				for _, urls := range p.links {
					go parseUrl(urls, p)
				}
				time.Sleep(time.Minute * time.Duration(p.period))
			}
		}
	}()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case posts := <-p.postCh:
				err := p.db.AddNews(posts)
				if err != nil {
					log.Println(fmt.Errorf("rss - start - storage error: %w", err))
				}
			case err := <-p.errorCh:
				log.Println(fmt.Errorf("rss - start - parser error: %w", err))
			}
		}
	}()
}

// Возвращает массив раскодированных новостей
func Parse(url string) ([]*postgres.Post, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var f RSS
	err = xml.Unmarshal(body, &f)
	if err != nil {
		return nil, err
	}
	var items []*postgres.Post
	for _, item := range f.Channel.Items {
		var post postgres.Post
		post.Title = item.Title
		post.Content = item.Content
		post.Content = strip.StripTags(post.Content)
		item.PubTime = strings.ReplaceAll(item.PubTime, ",", "")
		t, err := time.Parse("Mon 2 Jan 2006 15:04:05 -0700", item.PubTime)
		if err != nil {
			t, err = time.Parse("Mon 2 Jan 2006 15:04:05 GMT", item.PubTime)
		}
		if err == nil {
			post.PubTime = t.Unix()
		}

		post.Link = item.Link
		items = append(items, &post)
	}
	return items, nil
}

// Чтение rss-потока и отправка раскодированных постов и ошибок в каналы.
func parseUrl(url string, p *Parser) {
	log.Println("парсим сайт")
	for {
		feeds, err := Parse(url)
		if err != nil {
			p.errorCh <- err
			continue
		}
		p.postCh <- feeds
		return
	}
}
