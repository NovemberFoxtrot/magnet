package magnet

import (
    "net/http"
    "github.com/hoisie/mustache"
)

func LoginHandler(r *http.Request, w http.ResponseWriter) {
    // TODO:
    // - csrf token
    // - flash messages
    w.Write([]byte(mustache.RenderFileInLayout("base.html", "login.html", nil)))
}

// Not implemented
func LogoutHandler() {
}

// Not implemented
func SignUpHandler() {
}