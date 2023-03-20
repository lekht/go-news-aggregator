package postgres

// Пост из rss потока
type Post struct {
	ID      int
	Title   string
	Content string
	PubTime int64
	Link    string
}
