package main

import (
	r "github.com/christopherhesse/rethinkgo"
	m "github.com/codegangsta/martini"
	s "github.com/gorilla/sessions"
	"github.com/hoisie/mustache"
	"github.com/justinas/nosurf"
	h "net/http"
	"strconv"
	"strings"
	"time"
)

type Bookmark struct {
	Title   string
	Tags    []string
	Url     string
	Created float64
	User    string
	Date    string
}

func GetBookmarks(page int64, dbSession *r.Session, userId string) []Bookmark {
	var bookmarks []Bookmark

	err := r.Db("magnet").
		Table("bookmarks").
		Filter(r.Row.Attr("User").
		Eq(userId)).
		OrderBy(r.Desc("Created")).
		Skip(50 * page).
		Limit(50).
		Run(dbSession).
		All(&bookmarks)

	if err == nil {
		for i, _ := range bookmarks {
			if len(bookmarks[i].Tags) < 1 {
				bookmarks[i].Tags = []string{"No tags"}
			}
		}
	}

	return bookmarks
}

func GetBookmarksHandler(params m.Params, req *h.Request, w h.ResponseWriter, cs *s.CookieStore, dbSession *r.Session) {
	_, userId := GetUserData(cs, req)
	page, _ := strconv.ParseInt(params["page"], 10, 16)
	bookmarks := GetBookmarks(page, dbSession, userId)
	JsonDataResponse(200, bookmarks, req, w)
}

func IndexHandler(req *h.Request, w h.ResponseWriter, cs *s.CookieStore, dbSession *r.Session) {
	username, userId := GetUserData(cs, req)
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
	_, userId := GetUserData(cs, req)
	bookmark := new(Bookmark)
	bookmark.Title = req.PostFormValue("title")
	if req.PostFormValue("tags") != "" {
		bookmark.Tags = strings.Split(req.PostFormValue("tags"), ",")
		for i, v := range bookmark.Tags {
			bookmark.Tags[i] = strings.ToLower(strings.TrimSpace(v))
		}
	}
	bookmark.Created = float64(time.Now().Unix())
	bookmark.Date = time.Unix(int64(bookmark.Created), 0).Format("Jan 2, 2006 at 3:04pm")
	bookmark.Url = req.PostFormValue("url")
	bookmark.User = userId

	var response r.WriteResponse
	r.Db("magnet").
		Table("bookmarks").
		Insert(bookmark).
		Run(dbSession).
		One(&response)

	if response.Inserted > 0 {
		WriteJsonResponse(200, false, response.GeneratedKeys[0].(string), req, w)
	} else {
		WriteJsonResponse(200, true, "Error inserting bookmark.", req, w)
	}
}

func EditBookmarkHandler() {
}

func DeleteBookmarkHandler() {
}
