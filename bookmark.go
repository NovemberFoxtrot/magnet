package main

import (
	r "github.com/christopherhesse/rethinkgo"
)

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
func GetBookmarks(page int64, dbSession *r.Session, userID string) []Bookmark {
	var bookmarks []Bookmark

	err := r.Db("magnet").
		Table("bookmarks").
		Filter(r.Row.Attr("User").
		Eq(userID)).
		OrderBy(r.Desc("Created")).
		Skip(50 * page).
		Limit(50).
		Run(dbSession).
		All(&bookmarks)

	if err == nil {
		for i := range bookmarks {
			if len(bookmarks[i].Tags) < 1 {
				bookmarks[i].Tags = []string{"No tags"}
			}
		}
	}

	return bookmarks
}
