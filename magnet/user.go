package magnet

import (
    h "net/http"
    m "github.com/hoisie/mustache"
    r "github.com/christopherhesse/rethinkgo"
    s "github.com/gorilla/sessions"
    "github.com/justinas/nosurf"
    "crypto/sha1"
    "encoding/base64"
    "time"
)

type User struct {
    Username string `json:"Username"`
    Email    string `json:"Email"`
    Password string `json:"Password"`
}

type Session struct {
    UserId    string `json:userid`
    Expires    int64 `json:expires`
}

func cryptPassword(password, salt string) string {
    hash := sha1.New()
    hash.Write([]byte(password + salt))
    return string(base64.URLEncoding.EncodeToString(hash.Sum(nil)))
}

func LoginHandler(r *h.Request, w h.ResponseWriter, cs *s.CookieStore) {
    context := map[string]interface{} {
        "title" : "Access magnet",
        "csrf_token" : nosurf.Token(r),
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
            session.Save(req, w)
            WriteJsonResponse(200, false, "User correctly logged in.", req, w)
        }
    }
}

// Not implemented yet
func LogoutHandler() {
}

func SignUpHandler(req *h.Request, w h.ResponseWriter, dbSession *r.Session, cs *s.CookieStore, cfg *Config) {
    user := new(User)
    user.Username = req.PostFormValue("username")
    user.Email = req.PostFormValue("email")
    user.Password = cryptPassword(req.PostFormValue("password"), cfg.SecretKey)
    
    if len(user.Username) == 0 || len(user.Email) == 0 {
        WriteJsonResponse(200, true, "Empty fields.", req, w)
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
        WriteJsonResponse(200, true, "Username or email taken.", req, w)
    } else {
        // Can insert
        var response r.WriteResponse
        err = r.Db("magnet").
                Table("users").
                Insert(user).
                Run(dbSession).
                One(&response)
                
        if err != nil {
            WriteJsonResponse(200, true, "There was an error creating the user.", req, w)
        } else {
            WriteJsonResponse(201, false, "New user created.", req, w)
        }
    }
}