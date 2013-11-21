package main

import (
	"github.com/codegangsta/martini"
	s "github.com/gorilla/sessions"
    "./magnet"
    "github.com/justinas/nosurf"
    r "github.com/christopherhesse/rethinkgo"
    "encoding/json"
    "net/http"
    "fmt"
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
    //m.Get("/tag/:tag", magnet.ResponseAuthentication, magnet.GetTagHandler)
    //m.Post("/tag", magnet.ResponseAuthentication, magnet.NewTagHandler)
    //m.Put("/tag/:tag", magnet.ResponseAuthentication, magnet.EditTagHandler)
    //m.Delete("/tag/:tag", magnet.ResponseAuthentication, magnet.EditTagHandler)
    
    // Bookmark-related routes
    //m.Get("/bookmarks/:page", magnet.ResponseAuthentication, magnet.GetBookmarksHandler)
    //m.Post("/bookmark", magnet.ResponseAuthentication, magnet.NewBookmarkHandler)
    //m.Put("/bookmark/:bookmark", magnet.ResponseAuthentication, magnet.EditBookmarkHandler)
    //m.Delete("/bookmark/:bookmark", magnet.ResponseAuthentication, magnet.EditBookmarkHandler)
	
    // Search
    //m.Post("/search", magnet.ResponseAuthentication, magnet.SearchHandler)
    
    // User-related routes
    m.Post("/login", magnet.LoginPostHandler)
    m.Post("/logout", magnet.LogoutHandler)
    m.Post("/signup", magnet.SignUpHandler)

    // Home
	m.Get("/", magnet.Authentication, magnet.IndexHandler)
    csrfHandler := nosurf.New(m)
    csrfHandler.SetFailureHandler(http.HandlerFunc(magnet.CsrfFailHandler))

	http.ListenAndServe(config.Port, csrfHandler)
}