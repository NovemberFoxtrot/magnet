package magnet

import (
	s "github.com/gorilla/sessions"
    h "net/http"
    r "github.com/christopherhesse/rethinkgo"
    "encoding/json"
    "fmt"
)

func WriteJsonResponse(status int, error bool, message string, r *h.Request, w h.ResponseWriter) {
    w.Header().Set("Content-Type", "application/json")
    resp := make(map[string]interface{})
    resp["status"] = status
    resp["message"] = message
    resp["error"] = error
    jsonResp, _ := json.Marshal(resp)
    w.WriteHeader(status)
    w.Write(jsonResp)
}

func Authentication(cs *s.CookieStore, req *h.Request, w h.ResponseWriter, dbSession *r.Session) {
    session, _ := cs.Get(req, "magnet_session")

    var response map[string]interface{}
    fmt.Println(session.Values["session_id"])
    userId := ""
    err := r.Db("magnet").
            Table("sessions").
            Get(session.Values["session_id"]).
            Run(dbSession).
            One(&response)

    if err == nil && len(response) > 0 {
        userId = response["UserId"].(string)
    }
    
    if userId == "" {
        LoginHandler(req, w, cs)
    }
}

func CsrfFailHandler(w h.ResponseWriter, r *h.Request) {
    WriteJsonResponse(200, true, "Token invalid.", r, w)
}