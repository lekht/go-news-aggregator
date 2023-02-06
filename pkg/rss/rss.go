package rss

import (
	"encoding/xml"
	"io"
	"net/http"

	"github.com/lekht/go-news-aggregator/pkg/db"
)

type Feeds struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Title       string `xml:"title"`
	Description string `xml:"description"`
	Link        string `xml:"link"`
	Items       []Item `xml:"item"`
}

type Item struct {
	Title           string `xml:"title"`
	Description     string `xml:"description"`
	PublicationDate string `xml:"pubDate"`
	Link            string `xml:"link"`
}

func Parse(url string) ([]db.Post, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var f Feeds
	err = xml.Unmarshal(body, &f)
	if err != nil {
		return nil, err
	}
	var items []db.Post
	for _, item := range f.Channel.Items {
		var p db.Post
		p.Title = item.Title
		p.Description = item.Description
		p.PubDate = item.PublicationDate
		p.Url = item.Link
		items = append(items, p)
	}
	return items, nil
}
