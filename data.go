package main

import (
	"github.com/christopherhesse/rethinkgo"
	"github.com/codegangsta/martini"
	"github.com/gorilla/sessions"
	"strconv"
)

type Connection struct {
	session *rethinkgo.Session
}

func NewBookmark(userID string, dbSession *rethinkgo.Session, bookmark map[string]interface{}) (rethinkgo.WriteResponse, error) {
	var response rethinkgo.WriteResponse

	err := rethinkgo.Db("magnet").
		Table("bookmarks").
		Insert(bookmark).
		Run(dbSession).
		One(&response)

	return response, err
}

func DeleteBookmark(userID string, params martini.Params, dbSession *rethinkgo.Session) (rethinkgo.WriteResponse, error) {
	var response rethinkgo.WriteResponse

	err := rethinkgo.Db("magnet").
		Table("bookmarks").
		Filter(rethinkgo.Row.Attr("User").
		Eq(userID).
		And(rethinkgo.Row.Attr("id").
		Eq(params["bookmark"]))).
		Delete().
		Run(dbSession).
		One(&response)

	return response, err
}

func EditBookmark(userID string, params martini.Params, dbSession *rethinkgo.Session, bookmark map[string]interface{}) (rethinkgo.WriteResponse, error) {
	var response rethinkgo.WriteResponse

	err := rethinkgo.Db("magnet").
		Table("bookmarks").
		Filter(rethinkgo.Row.Attr("User").
		Eq(userID).
		And(rethinkgo.Row.Attr("id").
		Eq(params["bookmark"]))).
		Update(bookmark).
		Run(dbSession).
		One(&response)

	return response, err
}

func Search(userID string, params martini.Params, dbSession *rethinkgo.Session, query string) ([]interface{}, error) {
	var response []interface{}
	page, _ := strconv.ParseInt(params["page"], 10, 16)

	err := rethinkgo.Db("magnet").
		Table("bookmarks").
		Filter(rethinkgo.Row.Attr("Title").
		Match("(?i)" + query).
		And(rethinkgo.Row.Attr("User").
		Eq(userID))).
		OrderBy(rethinkgo.Desc("Created")).
		Skip(50 * page).
		Limit(50).
		Run(dbSession).
		All(&response)

	return response, err
}

func GetTag(userID string, params martini.Params, dbSession *rethinkgo.Session) ([]interface{}, error) {
	var response []interface{}
	page, _ := strconv.ParseInt(params["page"], 10, 16)

	err := rethinkgo.Db("magnet").
		Table("bookmarks").
		Filter(rethinkgo.Row.Attr("User").
		Eq(userID).
		And(rethinkgo.Row.Attr("Tags").
		Contains(params["tag"]))).
		OrderBy(rethinkgo.Desc("Created")).
		Skip(50 * page).
		Limit(50).
		Run(dbSession).
		All(&response)

	return response, err
}

func LoginPost(dbSession *rethinkgo.Session, username, password string) ([]interface{}, error) {
	var response []interface{}

	err := rethinkgo.Db("magnet").
		Table("users").
		Filter(rethinkgo.Row.Attr("Username").
		Eq(username).
		And(rethinkgo.Row.Attr("Password").
		Eq(password))).
		Run(dbSession).
		All(&response)

	return response, err
}

func LoginPostInsertSession(dbSession *rethinkgo.Session, session Session) (rethinkgo.WriteResponse, error) {
	var response rethinkgo.WriteResponse

	err := rethinkgo.Db("magnet").
		Table("sessions").
		Insert(session).
		Run(dbSession).
		One(&response)

	return response, err
}

func Logout(dbSession *rethinkgo.Session, session *sessions.Session) (rethinkgo.WriteResponse, error) {
	var response rethinkgo.WriteResponse

	err := rethinkgo.Db("magnet").
		Table("sessions").
		Get(session.Values["session_id"]).
		Delete().
		Run(dbSession).
		One(&response)

	return response, err
}

func SignUp(dbSession *rethinkgo.Session, user *User) ([]interface{}, error) {
	var response []interface{}

	err := rethinkgo.Db("magnet").
		Table("users").
		Filter(rethinkgo.Row.Attr("Username").
		Eq(user.Username).
		Or(rethinkgo.Row.Attr("Email").
		Eq(user.Email))).
		Run(dbSession).
		All(&response)

	return response, err
}
