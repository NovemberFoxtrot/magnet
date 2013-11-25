package main

import (
	"fmt"
	r "github.com/christopherhesse/rethinkgo"
	m "github.com/codegangsta/martini"
	s "github.com/gorilla/sessions"
	"github.com/hoisie/mustache"
	"github.com/justinas/nosurf"
	h "net/http"
	"strconv"
	"strings"
)

type Bookmark struct {
	Title   string
	Tags    []string
	Url     string
	Created int64
	User    string
}

func GetBookmarks(page int64, dbSession *r.Session, userId string) []Bookmark {
	var bookmarks []Bookmark

	r.Db("magnet").
		Table("bookmarks").
		Filter(r.Row.Attr("User").
		Eq(userId)).
		OrderBy(r.Desc("id")).
		Skip(50 * page).
		Limit(50).
		Run(dbSession).
		All(&bookmarks)

	return bookmarks
}

func GetBookmarksHandler(params m.Params, req *h.Request, w h.ResponseWriter, cs *s.CookieStore, dbSession *r.Session) {
	userId, _ := GetUserData(cs, req)
	page, _ := strconv.ParseInt(params["page"], 10, 16)
	bookmarks := GetBookmarks(page, dbSession, userId)
	JsonDataResponse(200, bookmarks, req, w)
}

func IndexHandler(req *h.Request, w h.ResponseWriter, cs *s.CookieStore, dbSession *r.Session) {
	userId, username := GetUserData(cs, req)
	context := map[string]interface{}{
		"title":      "Magnet",
		"csrf_token": nosurf.Token(req),
		"bookmarks":  GetBookmarks(0, dbSession, userId),
		"tags":       GetTags(dbSession, userId),
		"username":   username,
	}
	w.Write([]byte(mustache.RenderFileInLayout("templates/home.mustache", "templates/base.mustache", context)))
}

func NewBookmarkHandler(req *h.Request, w h.ResponseWriter, cs *s.CookieStore, dbSession *r.Session) {
	userId, _ := GetUserData(cs, req)
	bookmark := new(Bookmark)
	bookmark.Title = req.PostFormValue("title")
	bookmark.Tags = strings.Split(req.PostFormValue("tags"), ",")
	bookmark.Url = req.PostFormValue("url")
	bookmark.User = userId

	fmt.Println(bookmark.Title)
	fmt.Println(bookmark.Tags)
	fmt.Println(bookmark.Url)
}

func EditBookmarkHandler() {
}

func DeleteBookmarkHandler() {
}
