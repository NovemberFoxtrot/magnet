package magnet

import (
	"crypto/sha1"
	"encoding/base64"
	r "github.com/christopherhesse/rethinkgo"
	s "github.com/gorilla/sessions"
	m "github.com/hoisie/mustache"
	"github.com/justinas/nosurf"
	h "net/http"
	"regexp"
	"time"
)

type User struct {
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

type Session struct {
	UserId  string `json:UserId`
	Expires int64  `json:Expires`
}

func GetUserData(cs *s.CookieStore, req *h.Request) (string, string) {
	session, _ := cs.Get(req, "magnet_session")
	return session.Values["username"].(string), session.Values["user_id"].(string)
}

func cryptPassword(password, salt string) string {
	hash := sha1.New()
	hash.Write([]byte(password + salt))
	return string(base64.URLEncoding.EncodeToString(hash.Sum(nil)))
}

func LoginHandler(r *h.Request, w h.ResponseWriter) {
	context := map[string]interface{}{
		"title":      "Access magnet",
		"csrf_token": nosurf.Token(r),
	}
	w.Write([]byte(m.RenderFileInLayout("templates/login.mustache", "templates/base.mustache", context)))
}

func LoginPostHandler(req *h.Request, w h.ResponseWriter, cs *s.CookieStore, cfg *Config, dbSession *r.Session) {
	username := req.PostFormValue("username")
	password := cryptPassword(req.PostFormValue("password"), cfg.SecretKey)
	var response []interface{}
	err := r.Db("magnet").
		Table("users").
		Filter(r.Row.Attr("Username").
		Eq(username).
		And(r.Row.Attr("Password").
		Eq(password))).
		Run(dbSession).
		All(&response)

	if err != nil || len(response) == 0 {
		WriteJsonResponse(200, true, "Invalid username or password.", req, w)
	} else {
		// Store session
		userId := response[0].(map[string]interface{})["id"].(string)
		session := Session{UserId: userId,
			Expires: time.Now().Unix() + int64(cfg.SessionExpires)}

		var response r.WriteResponse
		err = r.Db("magnet").
			Table("sessions").
			Insert(session).
			Run(dbSession).
			One(&response)

		if err != nil || response.Inserted < 1 {
			WriteJsonResponse(200, true, "Error creating the user session.", req, w)
		} else {
			session, _ := cs.Get(req, "magnet_session")
			session.Values["session_id"] = response.GeneratedKeys[0]
			session.Values["username"] = username
			session.Values["user_id"] = userId
			session.Save(req, w)
			WriteJsonResponse(200, false, "User correctly logged in.", req, w)
		}
	}
}

func LogoutHandler(cs *s.CookieStore, req *h.Request, dbSession *r.Session, w h.ResponseWriter) {
	session, _ := cs.Get(req, "magnet_session")
	var response r.WriteResponse

	r.Db("magnet").
		Table("sessions").
		Get(session.Values["session_id"]).
		Delete().
		Run(dbSession).
		One(&response)

	session.Values["user_id"] = ""
	session.Values["session_id"] = ""
	session.Values["username"] = ""
	session.Save(req, w)

	h.Redirect(w, req, "/", 301)
}

func SignUpHandler(req *h.Request, w h.ResponseWriter, dbSession *r.Session, cs *s.CookieStore, cfg *Config) {
	user := new(User)
	req.ParseForm()
	user.Username = req.PostFormValue("username")
	user.Email = req.PostFormValue("email")
	user.Password = cryptPassword(req.PostFormValue("password"), cfg.SecretKey)
	errors := ""

	if len(user.Username) == 0 || len(user.Email) == 0 {
		errors += "Empty fields. "
	}

	exp, _ := regexp.Compile(`[a-zA-Z0-9._%+-]+@([a-zA-Z0-9-]+\.)+[A-Za-z]{2,6}`)

	if !exp.MatchString(user.Email) {
		errors += "Invalid email address. "
	}

	var response []interface{}
	err := r.Db("magnet").
		Table("users").
		Filter(r.Row.Attr("Username").
		Eq(user.Username).
		Or(r.Row.Attr("Email").
		Eq(user.Email))).
		Run(dbSession).
		All(&response)

	if err != nil || len(response) != 0 {
		errors += "Username or email taken."
	} else {
		var response r.WriteResponse
		err = r.Db("magnet").
			Table("users").
			Insert(user).
			Run(dbSession).
			One(&response)

		if err != nil {
			errors += "There was an error creating the user."
		} else {
			WriteJsonResponse(201, false, "New user created.", req, w)
		}
	}

	if errors != "" {
		WriteJsonResponse(200, true, errors, req, w)
	}
}

func RequestNewToken(r *h.Request, w h.ResponseWriter) {
	WriteJsonResponse(200, false, nosurf.Token(r), r, w)
}
