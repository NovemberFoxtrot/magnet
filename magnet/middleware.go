package magnet

import (
	s "github.com/gorilla/sessions"
    h "net/http"
    r "github.com/christopherhesse/rethinkgo"
)

func Authentication(cs *s.CookieStore, req *h.Request, w h.ResponseWriter, dbSession *r.Session) {
    session, _ := cs.Get(req, "magnet_session")

    var response map[string]interface{}
    var userId string
    err := r.Db("magnet").
            Table("sessions").
            Filter(r.Row.Attr("id").
            Eq(session.Values["session_id"])).
            Run(dbSession).
            One(&response)
    if err == nil {
        userId = response["UserId"].(string)
    } else {
        userId = ""
    }
    
    if userId == "" {
        LoginHandler(req, w, cs)
    }
}

func CsrfFailHandler(w h.ResponseWriter, r *h.Request) {
    w.Header().Set("Content-Type", "text/plain")
    w.Write([]byte("Bad access token.\n"))
}