package main

// Bookmark for JSON schema
type Bookmark struct {
	ID      string `json:"id"`
	Title   string
	Tags    []string
	URL     string
	Created float64
	User    string
	Date    string
}

// GetBookmarks fetches bookmarks from rethinkdb
func GetBookmarks(page int64, connection *Connection, userID string) []Bookmark {
	bookmarks, err := connection.GetBookmarks(userID, page)

	if err == nil {
		for i := range bookmarks {
			if len(bookmarks[i].Tags) < 1 {
				bookmarks[i].Tags = []string{"No tags"}
			}
		}
	}

	return bookmarks
}
