package main

import (
	"encoding/json"
	"github.com/gorilla/sessions"
	"net/http"
	"net/url"
	"time"
)

// JSONDataResponse writes JSON data to ResponseWriter
func JSONDataResponse(status int, err bool, data interface{}, r *http.Request, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]interface{})
	resp["status"] = status
	resp["data"] = data
	resp["error"] = err
	jsonResp, _ := json.Marshal(resp)
	w.WriteHeader(status)
	w.Write(jsonResp)
}

// WriteJSONResponse writes JSON to the ResponseWriter
func WriteJSONResponse(status int, error bool, message string, r *http.Request, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	resp := make(map[string]interface{})
	resp["status"] = status
	resp["message"] = message
	resp["error"] = error
	jsonResp, _ := json.Marshal(resp)
	w.WriteHeader(status)
	w.Write(jsonResp)
}

// GetUserID fetches userID from rethinkdb
func GetUserID(cs *sessions.CookieStore, req *http.Request, connection *Connection) string {
	session, _ := cs.Get(req, "magnet_session")
	var response map[string]interface{}

	userID := ""

	response, err := connection.GetUnexpiredSession(session)

	if err == nil && len(response) > 0 {
		if int64(response["Expires"].(float64)) > time.Now().Unix() {
			userID = response["UserId"].(string)
		}
	}

	return userID
}

// AuthRequired checks user session
func AuthRequired(cs *sessions.CookieStore, req *http.Request, w http.ResponseWriter, connection *Connection) {
	if GetUserID(cs, req, connection) == "" {
		WriteJSONResponse(401, true, "User is not logged in.", req, w)
	}
}

// IsValidURL checks if URL can be parsed
func IsValidURL(urlStr string) bool {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	return parsedURL.IsAbs()
}
