package main

import (
	"github.com/christopherhesse/rethinkgo"
	"github.com/codegangsta/martini"
	"github.com/gorilla/sessions"
	"github.com/hoisie/mustache"
	"github.com/justinas/nosurf"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// CsrfFailHandler writes invalid token response
func CsrfFailHandler(w http.ResponseWriter, r *http.Request) {
	WriteJSONResponse(200, true, "Provided token is not valid.", r, w)
}

// GetBookmarksHandler writes bookmarks to JSON data
func GetBookmarksHandler(params martini.Params, req *http.Request, w http.ResponseWriter, cs *sessions.CookieStore, dbSession *rethinkgo.Session) {
	_, userID := GetUserData(cs, req)
	page, _ := strconv.ParseInt(params["page"], 10, 16)
	bookmarks := GetBookmarks(page, dbSession, userID)
	JSONDataResponse(200, false, bookmarks, req, w)
}

// IndexHandler writes out templates
func IndexHandler(req *http.Request, w http.ResponseWriter, cs *sessions.CookieStore, dbSession *rethinkgo.Session) {
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
func NewBookmarkHandler(req *http.Request, w http.ResponseWriter, cs *sessions.CookieStore, dbSession *rethinkgo.Session) {
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

		response, _ := NewBookmark(userID, dbSession, bookmark)

		if response.Inserted > 0 {
			WriteJSONResponse(200, false, response.GeneratedKeys[0], req, w)
		} else {
			WriteJSONResponse(200, true, "Error inserting bookmark.", req, w)
		}
	}
}

// EditBookmarkHandler writes out response to editing a URL
func EditBookmarkHandler(req *http.Request, w http.ResponseWriter, cs *sessions.CookieStore, dbSession *rethinkgo.Session, params martini.Params) {
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

		response, err := EditBookmark(userID, params, dbSession, bookmark)

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
func DeleteBookmarkHandler(params martini.Params, req *http.Request, w http.ResponseWriter, cs *sessions.CookieStore, dbSession *rethinkgo.Session) {
	_, userID := GetUserData(cs, req)

	response, err := DeleteBookmark(userID, params, dbSession)

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
func SearchHandler(params martini.Params, req *http.Request, w http.ResponseWriter, cs *sessions.CookieStore, dbSession *rethinkgo.Session) {
	_, userID := GetUserData(cs, req)
	query := req.PostFormValue("query")

	response, err := Search(userID, params, dbSession, query)

	if err != nil {
		WriteJSONResponse(200, true, "Error retrieving bookmarks", req, w)
	} else {
		JSONDataResponse(200, false, response, req, w)
	}
}

// GetTagHandler fetches books for a given tag
func GetTagHandler(params martini.Params, req *http.Request, w http.ResponseWriter, cs *sessions.CookieStore, dbSession *rethinkgo.Session) {
	_, userID := GetUserData(cs, req)

	response, err := GetTag(userID, params, dbSession)

	if err != nil {
		WriteJSONResponse(200, true, "Error getting bookmarks for tag "+params["tag"], req, w)
	} else {
		JSONDataResponse(200, false, response, req, w)
	}
}

// LoginHandler writes out login template
func LoginHandler(r *http.Request, w http.ResponseWriter) {
	context := map[string]interface{}{
		"title":      "Access magnet",
		"csrf_token": nosurf.Token(r),
	}
	w.Write([]byte(mustache.RenderFileInLayout("templates/login.mustache", "templates/base.mustache", context)))
}

// LoginPostHandler writes out login response
func LoginPostHandler(req *http.Request, w http.ResponseWriter, cs *sessions.CookieStore, cfg *Config, dbSession *rethinkgo.Session) {
	username := req.PostFormValue("username")
	password := cryptPassword(req.PostFormValue("password"), cfg.SecretKey)

	var response []interface{}

	response, err := LoginPost(dbSession, username, password)

	if err != nil || len(response) == 0 {
		WriteJSONResponse(200, true, "Invalid username or password.", req, w)
	} else {
		// Store session
		userID := response[0].(map[string]interface{})["id"].(string)
		session := Session{UserID: userID,
			Expires: time.Now().Unix() + int64(cfg.SessionExpires)}

		response, err := LoginPostInsertSession(dbSession, session)

		if err != nil || response.Inserted < 1 {
			WriteJSONResponse(200, true, "Error creating the user session.", req, w)
		} else {
			session, _ := cs.Get(req, "magnet_session")
			session.Values["session_id"] = response.GeneratedKeys[0]
			session.Values["username"] = username
			session.Values["user_id"] = userID
			session.Save(req, w)
			WriteJSONResponse(200, false, "User correctly logged in.", req, w)
		}
	}
}

// LogoutHandler writes out logout response
func LogoutHandler(cs *sessions.CookieStore, req *http.Request, dbSession *rethinkgo.Session, w http.ResponseWriter) {
	session, _ := cs.Get(req, "magnet_session")

	_, _ = Logout(dbSession, session)

	session.Values["user_id"] = ""
	session.Values["session_id"] = ""
	session.Values["username"] = ""
	session.Save(req, w)

	http.Redirect(w, req, "/", 301)
}

// SignUpHandler writes out response to singing up
func SignUpHandler(req *http.Request, w http.ResponseWriter, dbSession *rethinkgo.Session, cs *sessions.CookieStore, cfg *Config) {
	user := new(User)
	req.ParseForm()
	user.Username = req.PostFormValue("username")
	user.Email = req.PostFormValue("email")
	user.Password = cryptPassword(req.PostFormValue("password"), cfg.SecretKey)
	errors := ""

	if len(user.Username) == 0 || len(user.Email) == 0 {
		errors += "Empty fields. "
	}

	exp, _ := regexp.Compile(`[a-zA-Z0-9._%+-]+@([a-zA-Z0-9-]+\.)+[A-Za-z]{2,6}`)

	if !exp.MatchString(user.Email) {
		errors += "Invalid email address. "
	}

	response, err := SignUp(dbSession, user)

	if err != nil || len(response) != 0 {
		errors += "Username or email taken."
	} else {
		_, err = SignUpInsert(dbSession, user)

		if err != nil {
			errors += "There was an error creating the user."
		} else {
			WriteJSONResponse(201, false, "New user created.", req, w)
		}
	}

	if errors != "" {
		WriteJSONResponse(200, true, errors, req, w)
	}
}
