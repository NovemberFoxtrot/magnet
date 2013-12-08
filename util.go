package main

import (
	"encoding/json"
	r "github.com/christopherhesse/rethinkgo"
	s "github.com/gorilla/sessions"
	h "net/http"
	"net/url"
	"time"
)

func JsonDataResponse(status int, err bool, data interface{}, r *h.Request, w h.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]interface{})
	resp["status"] = status
	resp["data"] = data
	resp["error"] = err
	jsonResp, _ := json.Marshal(resp)
	w.WriteHeader(status)
	w.Write(jsonResp)
}

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

func GetUserId(cs *s.CookieStore, req *h.Request, dbSession *r.Session) string {
	session, _ := cs.Get(req, "magnet_session")
	var response map[string]interface{}
	userId := ""
	// Get user session if it hasn't expired yet
	err := r.Db("magnet").
		Table("sessions").
		Get(session.Values["session_id"]).
		Run(dbSession).
		One(&response)

	if err == nil && len(response) > 0 {
		if int64(response["Expires"].(float64)) > time.Now().Unix() {
			userId = response["UserId"].(string)
		}
	}

	return userId
}

func AuthRequired(cs *s.CookieStore, req *h.Request, w h.ResponseWriter, dbSession *r.Session) {
	if GetUserId(cs, req, dbSession) == "" {
		WriteJsonResponse(401, true, "User is not logged in.", req, w)
	}
}

func CsrfFailHandler(w h.ResponseWriter, r *h.Request) {
	WriteJsonResponse(200, true, "Provided token is not valid.", r, w)
}

func IsValidUrl(urlStr string) bool {
	parsedUrl, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	return parsedUrl.IsAbs()
}
