package main

import (
	"encoding/json"
	"fmt"
	r "github.com/christopherhesse/rethinkgo"
	"github.com/codegangsta/martini"
	s "github.com/gorilla/sessions"
	"github.com/justinas/nosurf"
	h "net/http"
	"os"
	"time"
)

func initDatabase(connectionString string) *r.Session {
	session, err := r.Connect(connectionString, "magnet")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return nil
	}

	r.DbCreate("magnet").Run(session).Exec()
	r.TableCreate("users").Run(session).Exec()
	r.TableCreate("bookmarks").Run(session).Exec()
	r.TableCreate("sessions").Run(session).Exec()

	// Delete all expired sessions
	var rsp r.WriteResponse
	r.Db("magnet").
		Table("sessions").
		Filter(r.Row.Attr("Expires").
		Lt(time.Now().Unix())).
		Delete().
		Run(session).
		One(&rsp)

	return session
}

func main() {
	m := martini.Classic()

	// Read config
	reader, _ := os.Open("config.json")
	decoder := json.NewDecoder(reader)
	config := &Config{}
	decoder.Decode(&config)

	// Init database
	dbSession := initDatabase(config.ConnectionString)
	if dbSession == nil {
		os.Exit(2)
	}

	// Create a new cookie store
	store := s.NewCookieStore([]byte(config.SecretKey))

	// It will be available to all handlers as *sessions.CookieStore
	m.Map(store)

	// It will be available to all handlers as *r.Session
	m.Map(dbSession)

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
	m.Get("/", func(cs *s.CookieStore, req *h.Request, w h.ResponseWriter, dbSession *r.Session) {
		if GetUserID(cs, req, dbSession) == "" {
			LoginHandler(req, w)
		}
	}, IndexHandler)

	csrfHandler := nosurf.New(m)
	csrfHandler.SetFailureHandler(h.HandlerFunc(CsrfFailHandler))

	h.ListenAndServe(config.Port, csrfHandler)
}
