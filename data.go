package main

import (
	"github.com/christopherhesse/rethinkgo"
	"github.com/codegangsta/martini"
	"github.com/gorilla/sessions"
	"strconv"
	"time"
)

type Connection struct {
	session *rethinkgo.Session
}

func (c *Connection) NewBookmark(userID string, bookmark map[string]interface{}) (rethinkgo.WriteResponse, error) {
	var response rethinkgo.WriteResponse

	err := rethinkgo.Db("magnet").
		Table("bookmarks").
		Insert(bookmark).
		Run(c.session).
		One(&response)

	return response, err
}

func (c *Connection) DeleteBookmark(userID string, params martini.Params) (rethinkgo.WriteResponse, error) {
	var response rethinkgo.WriteResponse

	err := rethinkgo.Db("magnet").
		Table("bookmarks").
		Filter(rethinkgo.Row.Attr("User").
		Eq(userID).
		And(rethinkgo.Row.Attr("id").
		Eq(params["bookmark"]))).
		Delete().
		Run(c.session).
		One(&response)

	return response, err
}

func (c *Connection) EditBookmark(userID string, params martini.Params, bookmark map[string]interface{}) (rethinkgo.WriteResponse, error) {
	var response rethinkgo.WriteResponse

	err := rethinkgo.Db("magnet").
		Table("bookmarks").
		Filter(rethinkgo.Row.Attr("User").
		Eq(userID).
		And(rethinkgo.Row.Attr("id").
		Eq(params["bookmark"]))).
		Update(bookmark).
		Run(c.session).
		One(&response)

	return response, err
}

func (c *Connection) Search(userID string, params martini.Params, query string) ([]interface{}, error) {
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
		Run(c.session).
		All(&response)

	return response, err
}

func (c *Connection) GetTag(userID string, params martini.Params) ([]interface{}, error) {
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
		Run(c.session).
		All(&response)

	return response, err
}

func (c *Connection) LoginPost(username, password string) ([]interface{}, error) {
	var response []interface{}

	err := rethinkgo.Db("magnet").
		Table("users").
		Filter(rethinkgo.Row.Attr("Username").
		Eq(username).
		And(rethinkgo.Row.Attr("Password").
		Eq(password))).
		Run(c.session).
		All(&response)

	return response, err
}

func (c *Connection) LoginPostInsertSession(session Session) (rethinkgo.WriteResponse, error) {
	var response rethinkgo.WriteResponse

	err := rethinkgo.Db("magnet").
		Table("sessions").
		Insert(session).
		Run(c.session).
		One(&response)

	return response, err
}

func (c *Connection) Logout(session *sessions.Session) (rethinkgo.WriteResponse, error) {
	var response rethinkgo.WriteResponse

	err := rethinkgo.Db("magnet").
		Table("sessions").
		Get(session.Values["session_id"]).
		Delete().
		Run(c.session).
		One(&response)

	return response, err
}

func (c *Connection) SignUp(user *User) ([]interface{}, error) {
	var response []interface{}

	err := rethinkgo.Db("magnet").
		Table("users").
		Filter(rethinkgo.Row.Attr("Username").
		Eq(user.Username).
		Or(rethinkgo.Row.Attr("Email").
		Eq(user.Email))).
		Run(c.session).
		All(&response)

	return response, err
}

func (c *Connection) SignUpInsert(user *User) (rethinkgo.WriteResponse, error) {
	var response rethinkgo.WriteResponse

	err := rethinkgo.Db("magnet").
		Table("users").
		Insert(user).
		Run(c.session).
		One(&response)

	return response, err
}

func InitDatabase(session *rethinkgo.Session) {
	rethinkgo.DbCreate("magnet").Run(session).Exec()
	rethinkgo.TableCreate("users").Run(session).Exec()
	rethinkgo.TableCreate("bookmarks").Run(session).Exec()
	rethinkgo.TableCreate("sessions").Run(session).Exec()
}

func WipeExpiredSessions(session *rethinkgo.Session) (rethinkgo.WriteResponse, error) {
	var response rethinkgo.WriteResponse

	err := rethinkgo.Db("magnet").
		Table("sessions").
		Filter(rethinkgo.Row.Attr("Expires").
		Lt(time.Now().Unix())).
		Delete().
		Run(session).
		One(&response)

	return response, err
}
