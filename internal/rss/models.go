package rss

import "encoding/xml"

type RSS struct {
	XMLName xml.Name `xml:"rss"`
	Channel Channel  `xml:"channel"`
}

type Channel struct {
	Items []Item `xml:"item"`
}

type Item struct {
	Title   string `xml:"title"`
	Content string `xml:"description"`
	PubTime string `xml:"pubDate"`
	Link    string `xml:"link"`
}
