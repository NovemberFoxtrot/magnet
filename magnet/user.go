package magnet

import (
    h "net/http"
    m "github.com/hoisie/mustache"
    r "github.com/christopherhesse/rethinkgo"
    s "github.com/gorilla/sessions"
    "github.com/justinas/nosurf"
    "crypto/sha1"
    "encoding/base64"
    "fmt"
)

type User struct {
    Username string `json:"Username"`
    Email    string `json:"Email"`
    Password string `json:"Password"`
}

type Session struct {
    UserId    string
    Expires    int32
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

func LoginPostHandler(req *h.Request, w h.ResponseWriter, cs *s.CookieStore, cfg *Config) string {
    username := r.PostFormValue("username")
    password := cryptPassword(req.PostFormValue("password"), cfg.SecretKey)
    var response []interface{}
    err := r.Db("magnet").
            Table("users").
            Filter(r.Row.Attr("Username").
                Eq(user.Username).
                And(r.Row.Attr("Password").
                    Eq(user.Password))).
            Run(dbSession).
            All(&response)
            
    if err != nil || len(response) != 0 {
        // Login failed
    } else {
        // Login correct
        // Store session
    }
}

// Not implemented
func LogoutHandler() {
}

func SignUpHandler(req *h.Request, dbSession *r.Session, cs *s.CookieStore, cfg *Config) {
    user := new(User)
    user.Username = req.PostFormValue("username")
    user.Email = req.PostFormValue("email")
    user.Password = cryptPassword(req.PostFormValue("password"), cfg.SecretKey)
    
    if len(user.Username) == 0 || len(user.Email) == 0 {
        // Throw error
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
        // Username or email taken
    } else {
        // Can insert
        var response r.WriteResponse
        err = r.Db("magnet").
                Table("users").
                Insert(user).
                Run(dbSession).
                One(&response)
                
        if err != nil {
            // Error happened
        }
    }
    
    // TODO: Response
}