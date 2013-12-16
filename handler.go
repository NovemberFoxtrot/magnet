package main

import (
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

func Start(DB *Connection, config *Config) {
	// Create a new cookie store
	store := sessions.NewCookieStore([]byte(config.SecretKey))

	m := martini.Classic()

	// It will be available to all handlers as *sessions.CookieStore
	m.Map(store)

	// It will be available to all handlers as *connection *Connection
	m.Map(DB)

	// It will be available to all handlers as *Config
	m.Map(config)

	// public folder will serve the static content
	m.Use(martini.Static("public"))

	// Tag-related routes
	m.Get("/tag/:tag/:page", AuthRequired, GetTagHandler)

	// Bookmark-related routes
	m.Get("/bookmarks/:page", AuthRequired, GetBookmarksHandler)
	m.Post("/bookmark/new", AuthRequired, NewBookmarkHandler)
	m.Post("/bookmark/update/:bookmark", AuthRequired, EditBookmarkHandler)
	m.Delete("/bookmark/delete/:bookmark", AuthRequired, DeleteBookmarkHandler)

	// Search
	m.Post("/search/:page", AuthRequired, SearchHandler)

	// User-related routes
	m.Post("/login", LoginPostHandler)
	m.Get("/logout", AuthRequired, LogoutHandler)
	m.Post("/signup", SignUpHandler)
	m.Post("/new_token", AuthRequired, RequestNewToken)

	// Test
	m.Get("/test", TestHandler)

	// Home
	m.Get("/", func(cs *sessions.CookieStore, req *http.Request, w http.ResponseWriter, connection *Connection) {
		if GetUserID(cs, req, connection) == "" {
			LoginHandler(req, w)
		}
	}, IndexHandler)

	csrfHandler := nosurf.New(m)
	csrfHandler.SetFailureHandler(http.HandlerFunc(CsrfFailHandler))

	http.ListenAndServe(config.Port, csrfHandler)
}

// CsrfFailHandler writes invalid token response
func CsrfFailHandler(w http.ResponseWriter, r *http.Request) {
	WriteJSONResponse(200, true, "Provided token is not valid.", r, w)
}

// GetBookmarksHandler writes bookmarks to JSON data
func GetBookmarksHandler(params martini.Params, req *http.Request, w http.ResponseWriter, cs *sessions.CookieStore, connection *Connection) {
	_, userID := GetUserData(cs, req)
	page, _ := strconv.ParseInt(params["page"], 10, 16)
	bookmarks := GetBookmarks(page, connection, userID)
	JSONDataResponse(200, false, bookmarks, req, w)
}

// IndexHandler writes out templates
func IndexHandler(req *http.Request, w http.ResponseWriter, cs *sessions.CookieStore, connection *Connection) {
	username, userID := GetUserData(cs, req)
	context := map[string]interface{}{
		"title":      "Magnet",
		"csrf_token": nosurf.Token(req),
		"bookmarks":  GetBookmarks(0, connection, userID),
		"tags":       GetTags(connection, userID),
		"username":   username,
	}

	context["load_more"] = len(context["bookmarks"].([]Bookmark)) == 50

	w.Write([]byte(mustache.RenderFileInLayout("templates/home.mustache", "templates/base.mustache", context)))
}

// TestHandler runs the tests
func TestHandler(req *http.Request, w http.ResponseWriter, cs *sessions.CookieStore, connection *Connection) {
	w.Write([]byte(mustache.RenderFileInLayout("templates/test.mustache", "templates/test.base.mustache", nil)))
}

// NewBookmarkHandler writes out new bookmark JSON response
func NewBookmarkHandler(req *http.Request, w http.ResponseWriter, cs *sessions.CookieStore, connection *Connection) {
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

		response, _ := connection.NewBookmark(userID, bookmark)

		if response.Inserted > 0 {
			WriteJSONResponse(200, false, response.GeneratedKeys[0], req, w)
		} else {
			WriteJSONResponse(200, true, "Error inserting bookmark.", req, w)
		}
	}
}

// EditBookmarkHandler writes out response to editing a URL
func EditBookmarkHandler(req *http.Request, w http.ResponseWriter, cs *sessions.CookieStore, connection *Connection, params martini.Params) {
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

		response, err := connection.EditBookmark(userID, params, bookmark)

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
func DeleteBookmarkHandler(params martini.Params, req *http.Request, w http.ResponseWriter, cs *sessions.CookieStore, connection *Connection) {
	_, userID := GetUserData(cs, req)

	response, err := connection.DeleteBookmark(userID, params)

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
func SearchHandler(params martini.Params, req *http.Request, w http.ResponseWriter, cs *sessions.CookieStore, connection *Connection) {
	_, userID := GetUserData(cs, req)
	query := req.PostFormValue("query")

	response, err := connection.Search(userID, params, query)

	if err != nil {
		WriteJSONResponse(200, true, "Error retrieving bookmarks", req, w)
	} else {
		JSONDataResponse(200, false, response, req, w)
	}
}

// GetTagHandler fetches books for a given tag
func GetTagHandler(params martini.Params, req *http.Request, w http.ResponseWriter, cs *sessions.CookieStore, connection *Connection) {
	_, userID := GetUserData(cs, req)

	response, err := connection.GetTag(userID, params)

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
func LoginPostHandler(req *http.Request, w http.ResponseWriter, cs *sessions.CookieStore, cfg *Config, connection *Connection) {
	username := req.PostFormValue("username")
	password := cryptPassword(req.PostFormValue("password"), cfg.SecretKey)

	var response []interface{}

	response, err := connection.LoginPost(username, password)

	if err != nil || len(response) == 0 {
		WriteJSONResponse(200, true, "Invalid username or password.", req, w)
	} else {
		// Store session
		userID := response[0].(map[string]interface{})["id"].(string)
		session := Session{UserID: userID,
			Expires: time.Now().Unix() + int64(cfg.SessionExpires)}

		response, err := connection.LoginPostInsertSession(session)

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
func LogoutHandler(cs *sessions.CookieStore, req *http.Request, connection *Connection, w http.ResponseWriter) {
	session, _ := cs.Get(req, "magnet_session")

	_, _ = connection.Logout(session)

	session.Values["user_id"] = ""
	session.Values["session_id"] = ""
	session.Values["username"] = ""
	session.Save(req, w)

	http.Redirect(w, req, "/", 301)
}

// SignUpHandler writes out response to singing up
func SignUpHandler(req *http.Request, w http.ResponseWriter, connection *Connection, cs *sessions.CookieStore, cfg *Config) {
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

	response, err := connection.SignUp(user)

	if err != nil || len(response) != 0 {
		errors += "Username or email taken."
	} else {
		_, err = connection.SignUpInsert(user)

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
