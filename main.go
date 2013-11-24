package main

import (
	"./magnet"
	"encoding/json"
	"fmt"
	r "github.com/christopherhesse/rethinkgo"
	"github.com/codegangsta/martini"
	s "github.com/gorilla/sessions"
	"github.com/justinas/nosurf"
	h "net/http"
	"os"
)

func initDatabase(connectionString string) *r.Session {
	// TODO: Erase all expired sessions

	session, err := r.Connect(connectionString, "magnet")
	if err != nil {
		fmt.Println("Error connecting:", err)
		return nil
	}

	r.DbCreate("magnet").Run(session).Exec()
	r.TableCreate("users").Run(session).Exec()
	r.TableCreate("bookmarks").Run(session).Exec()
	r.TableCreate("sessions").Run(session).Exec()
	r.TableCreate("tags").Run(session).Exec()

	return session
}

func main() {
	// TODO
	// - Config init and map

	m := martini.Classic()

	// Read config
	reader, _ := os.Open("config.json")
	decoder := json.NewDecoder(reader)
	config := &magnet.Config{}
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
	// It will be available to all handlers as *magnet.Config
	m.Map(config)
	// public folder will serve the static content
	m.Use(martini.Static("public"))

	// Tag-related routes
	//m.Get("/tag/:tag", magnet.AuthRequired, magnet.GetTagHandler)
	//m.Post("/tag", magnet.AuthRequired, magnet.NewTagHandler)
	//m.Put("/tag/:tag", magnet.AuthRequired, magnet.EditTagHandler)
	//m.Delete("/tag/:tag", magnet.AuthRequired, magnet.EditTagHandler)

	// Bookmark-related routes
	m.Get("/bookmarks/:page", magnet.AuthRequired, magnet.GetBookmarksHandler)
	m.Post("/bookmark", magnet.AuthRequired, magnet.NewBookmarkHandler)
	m.Put("/bookmark/:bookmark", magnet.AuthRequired, magnet.EditBookmarkHandler)
	m.Delete("/bookmark/:bookmark", magnet.AuthRequired, magnet.DeleteBookmarkHandler)

	// Search
	//m.Post("/search", magnet.AuthRequired, magnet.SearchHandler)

	// User-related routes
	m.Post("/login", magnet.LoginPostHandler)
	m.Get("/logout", magnet.AuthRequired, magnet.LogoutHandler)
	m.Post("/signup", magnet.SignUpHandler)
	m.Post("/new_token", magnet.AuthRequired, magnet.RequestNewToken)

	// Home
	m.Get("/", func(cs *s.CookieStore, req *h.Request, w h.ResponseWriter, dbSession *r.Session) {
		if magnet.GetUserId(cs, req, dbSession) == "" {
			magnet.LoginHandler(req, w)
		}
	}, magnet.IndexHandler)

	csrfHandler := nosurf.New(m)
	csrfHandler.SetFailureHandler(h.HandlerFunc(magnet.CsrfFailHandler))

	h.ListenAndServe(config.Port, csrfHandler)
}
