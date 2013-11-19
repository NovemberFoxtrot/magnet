package main

import (
	"github.com/codegangsta/martini"
	"github.com/gorilla/sessions"
    "./magnet"
)

func main() {
    // TODO
    // - Database init and map
    // - Config init and map
    
	m := martini.Classic()
    
    // Create a new cookie store
    store := sessions.NewCookieStore([]byte("my fancy secret code"))
    // It will be available to all handlers as *sessions.CookieStore
    m.Map(store)
    // public folder will serve the static content
	m.Use(martini.Static("public"))

    // Tag-related routes
    m.Get("/tag/:tag", magnet.ResponseAuthentication, magnet.GetTagHandler)
    m.Post("/tag", magnet.ResponseAuthentication, magnet.NewTagHandler)
    m.Put("/tag/:tag", magnet.ResponseAuthentication, magnet.EditTagHandler)
    m.Delete("/tag/:tag", magnet.ResponseAuthentication, magnet.EditTagHandler)
    
    // Bookmark-related routes
    m.Get("/bookmarks/:page", magnet.ResponseAuthentication, magnet.GetBookmarksHandler)
    m.Post("/bookmark", magnet.ResponseAuthentication, magnet.NewBookmarkHandler)
    m.Put("/bookmark/:bookmark", magnet.ResponseAuthentication, magnet.EditBookmarkHandler)
    m.Delete("/bookmark/:bookmark", magnet.ResponseAuthentication, magnet.EditBookmarkHandler)
	
    // Search
    m.Post("/search", magnet.ResponseAuthentication, magnet.SearchHandler)
    
    // User-related routes
    m.Post("/login", magnet.LoginHandler)
    m.Post("/logout", magnet.LogoutHandler)
    m.Post("/signup", magnet.SignUpHandler)

    // Home
	m.Get("/", magnet.Authentication, magnet.IndexHandler)

	m.Run()
}