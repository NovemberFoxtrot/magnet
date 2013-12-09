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

// GetBookmarksHandler writes bookmarks to JSON data
func GetBookmarksHandler(params m.Params, req *h.Request, w h.ResponseWriter, cs *s.CookieStore, dbSession *r.Session) {
	_, userID := GetUserData(cs, req)
	page, _ := strconv.ParseInt(params["page"], 10, 16)
	bookmarks := GetBookmarks(page, dbSession, userID)
	JSONDataResponse(200, false, bookmarks, req, w)
}

// IndexHandler writes out templates
func IndexHandler(req *h.Request, w h.ResponseWriter, cs *s.CookieStore, dbSession *r.Session) {
	username, userID := GetUserData(cs, req)
	context := map[string]interface{}{
		"title":      "Magnet",
		"csrf_token": nosurf.Token(req),
		"bookmarks":  GetBookmarks(0, dbSession, userID),
		"tags":       GetTags(dbSession, userID),
		"username":   username,
	}
	context["load_more"] = len(context["bookmarks"].([]Bookmark)) == 2
	w.Write([]byte(mustache.RenderFileInLayout("templates/home.mustache", "templates/base.mustache", context)))
}

// NewBookmarkHandler writes out new bookmark JSON response
func NewBookmarkHandler(req *h.Request, w h.ResponseWriter, cs *s.CookieStore, dbSession *r.Session) {
	// We use a map instead of Bookmark because id would be ""
	bookmark := make(map[string]interface{})
	bookmark["Title"] = req.PostFormValue("title")
	bookmark["Url"] = req.PostFormValue("url")

	if !IsValidURL(bookmark["Url"].(string)) || len(bookmark["Title"].(string)) < 1 {
		WriteJSONResponse(200, true, "The url is not valid or the title is empty.", req, w)
	} else {
		_, userID := GetUserData(cs, req)
		if req.PostFormValue("tags") != "" {
			bookmark["Tags"] = strings.Split(req.PostFormValue("tags"), ",")
			for i, v := range bookmark["Tags"].([]string) {
				bookmark["Tags"].([]string)[i] = strings.ToLower(strings.TrimSpace(v))
			}
		}
		bookmark["Created"] = float64(time.Now().Unix())
		bookmark["Date"] = time.Unix(int64(bookmark["Created"].(float64)), 0).Format("Jan 2, 2006 at 3:04pm")
		bookmark["User"] = userID

		var response r.WriteResponse
		r.Db("magnet").
			Table("bookmarks").
			Insert(bookmark).
			Run(dbSession).
			One(&response)

		if response.Inserted > 0 {
			WriteJSONResponse(200, false, response.GeneratedKeys[0], req, w)
		} else {
			WriteJSONResponse(200, true, "Error inserting bookmark.", req, w)
		}
	}
}

// EditBookmarkHandler writes out response to editing a URL
func EditBookmarkHandler(req *h.Request, w h.ResponseWriter, cs *s.CookieStore, dbSession *r.Session, params m.Params) {
	// We use a map instead of Bookmark because id would be ""
	bookmark := make(map[string]interface{})
	bookmark["Title"] = req.PostFormValue("title")
	bookmark["Url"] = req.PostFormValue("url")

	if !IsValidURL(bookmark["Url"].(string)) || len(bookmark["Title"].(string)) < 1 {
		WriteJSONResponse(200, true, "The url is not valid or the title is empty.", req, w)
	} else {
		_, userID := GetUserData(cs, req)
		if req.PostFormValue("tags") != "" {
			bookmark["Tags"] = strings.Split(req.PostFormValue("tags"), ",")
			for i, v := range bookmark["Tags"].([]string) {
				bookmark["Tags"].([]string)[i] = strings.ToLower(strings.TrimSpace(v))
			}
		}

		var response r.WriteResponse
		err := r.Db("magnet").
			Table("bookmarks").
			Filter(r.Row.Attr("User").
			Eq(userID).
			And(r.Row.Attr("id").
			Eq(params["bookmark"]))).
			Update(bookmark).
			Run(dbSession).
			One(&response)

		if err != nil {
			WriteJSONResponse(200, true, "Error deleting bookmark.", req, w)
		} else {
			if response.Updated > 0 || response.Unchanged > 0 || response.Replaced > 0 {
				WriteJSONResponse(200, false, "Bookmark updated successfully.", req, w)
			} else {
				WriteJSONResponse(200, true, "Error updating bookmark.", req, w)
			}
		}
	}
}

// DeleteBookmarkHandler writes out response to deleting a bookmark
func DeleteBookmarkHandler(params m.Params, req *h.Request, w h.ResponseWriter, cs *s.CookieStore, dbSession *r.Session) {
	_, userID := GetUserData(cs, req)
	var response r.WriteResponse

	err := r.Db("magnet").
		Table("bookmarks").
		Filter(r.Row.Attr("User").
		Eq(userID).
		And(r.Row.Attr("id").
		Eq(params["bookmark"]))).
		Delete().
		Run(dbSession).
		One(&response)

	if err != nil {
		WriteJSONResponse(200, true, "Error deleting bookmark.", req, w)
	} else {
		if response.Deleted > 0 {
			WriteJSONResponse(200, false, "Bookmark deleted successfully.", req, w)
		} else {
			WriteJSONResponse(200, true, "Error deleting bookmark.", req, w)
		}
	}
}

// SearchHandler writes out response when searching for a URL
func SearchHandler(params m.Params, req *h.Request, w h.ResponseWriter, cs *s.CookieStore, dbSession *r.Session) {
	_, userID := GetUserData(cs, req)
	var response []interface{}
	page, _ := strconv.ParseInt(params["page"], 10, 16)
	query := req.PostFormValue("query")

	err := r.Db("magnet").
		Table("bookmarks").
		Filter(r.Row.Attr("Title").
		Match("(?i)" + query).
		And(r.Row.Attr("User").
		Eq(userID))).
		OrderBy(r.Desc("Created")).
		Skip(50 * page).
		Limit(50).
		Run(dbSession).
		All(&response)

	if err != nil {
		WriteJSONResponse(200, true, "Error retrieving bookmarks", req, w)
	} else {
		JSONDataResponse(200, false, response, req, w)
	}
}
