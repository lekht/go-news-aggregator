package rss

import "encoding/xml"

type RSS struct {
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
	Title   string `xml:"title"`
	Content string `xml:"description"`
	PubDate string `xml:"pubDate"`
	Link    string `xml:"link"`
}
