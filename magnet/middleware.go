package magnet

import (
	s "github.com/gorilla/sessions"
    h "net/http"
    r "github.com/christopherhesse/rethinkgo"
)

func Authentication(cs *s.CookieStore, req *h.Request, w h.ResponseWriter, dbSession *r.Session) {
    session, _ := cs.Get(req, "magnet_session")
    // TODO:
    // - implement db
    // - add db to params
    var response map[string]interface{}
    var userId int
    err := r.Table("sessions").Filter(r.Row.Attr("id").Eq(session.Values["session_id"])).Run(dbSession).One(&response)
    if err == nil {
        userId = response["UserId"].(int)
    } else {
        userId = 0
    }
    
    if userId < 1 {
        LoginHandler(req, w)
    }
}