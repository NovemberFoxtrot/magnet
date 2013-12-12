package main

import (
	"encoding/json"
	"github.com/codegangsta/martini"
	"github.com/gorilla/sessions"
	"github.com/justinas/nosurf"
	"net/http"
	"os"
)

func main() {
	m := martini.Classic()

	DB := &Connection{}

	// Read config
	reader, _ := os.Open("config.json")
	decoder := json.NewDecoder(reader)
	config := &Config{}
	decoder.Decode(&config)

	// Init database
	DB.initDatabase(config.ConnectionString)

	// Create a new cookie store
	store := sessions.NewCookieStore([]byte(config.SecretKey))

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
