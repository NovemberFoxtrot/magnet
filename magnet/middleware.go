package magnet

import (
	"github.com/gorilla/sessions"
    "net/http"
)

func Authentication(cs *sessions.CookieStore, r *http.Request, w http.ResponseWriter) {
    session, _ := cs.Get(r, "magnet_session")
    // TODO:
    // - implement db
    // - add db to params
    userId := db.UserIdForSession(session.Values["session_id"].(string))
    
    if userId < 1 {
        LoginHandler(r, w)
    }
}