package magnet

import (
    h "net/http"
    m "github.com/hoisie/mustache"
    r "github.com/christopherhesse/rethinkgo"
)

type User struct {
    Username string
    Email    string
    Password string
    Salt     string
}

type Session struct {
    UserId    string
    Expires    int32
}

func LoginHandler(r *h.Request, w hp.ResponseWriter) {
    // TODO:
    // - csrf token
    // - flash messages
    w.Write([]byte(m.RenderFileInLayout("base.mustache", "login.mustache", nil)))
}

// Not implemented
func LogoutHandler() {
}

// Not implemented
func SignUpHandler(req *h.Request, dbSession *r.Session) {
    csrfToken := req.PostFormValue("csrf_token")
    user := new(User)
    user.Username = req.PostFormValue("username")
    user.Email = req.PostFormValue("email")
    // TODO: hash and salt
    user.Password = req.PostFormValue("password")
    // TODO: insert
    // TODO: response
}