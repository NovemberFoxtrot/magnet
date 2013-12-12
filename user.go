package main

import (
	"crypto/sha1"
	"encoding/base64"
	"github.com/gorilla/sessions"
	"github.com/justinas/nosurf"
	"net/http"
)

// User for JSON schema
type User struct {
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

// Session for JSON schema
type Session struct {
	UserID  string `json:"UserId"`
	Expires int64  `json:"Expires"`
}

// GetUserData fetches user session data
func GetUserData(cs *sessions.CookieStore, req *http.Request) (string, string) {
	session, _ := cs.Get(req, "magnet_session")
	return session.Values["username"].(string), session.Values["user_id"].(string)
}

func cryptPassword(password, salt string) string {
	hash := sha1.New()
	hash.Write([]byte(password + salt))
	return string(base64.URLEncoding.EncodeToString(hash.Sum(nil)))
}

// RequestNewToken writes out nosurf token
func RequestNewToken(r *http.Request, w http.ResponseWriter) {
	WriteJSONResponse(200, false, nosurf.Token(r), r, w)
}
