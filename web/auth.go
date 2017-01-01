package web

import (
	"net/http"
	"strings"
	"time"
	"github.com/gorilla/sessions"
)

type User struct {
	Username string
	Id string
	Secret string
}

func GetSession(r *http.Request) (*sessions.Session) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, err := store.Get(r, "session-name")
	if err != nil {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		panic("Unexpected error!")
	}

	return session
}

// https://medium.com/@matryer/the-http-handlerfunc-wrapper-technique-in-golang-c60bf76e6124#.ipqttmx8e
func AuthWrapper(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := GetSession(r)

		user, ok := session.Values["user"].(*User)
		if !ok {
			LoginHandler(w, r)
			return
		}
		//log.Println("AuthWrapper user.Username:", user.Username)
		_ = user // avoid declared and not used

		fn(w, r)
	}
}


func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// If method is GET serve an html login page
	if r.Method != "POST" {
		http.ServeFile(w, r, "web/templates/login.html")
		return
	}

	// Grab the username/password from the submitted post form
	username := r.FormValue("username")
	password := r.FormValue("password")

	// If wrong password redirect to the login
	if !((strings.Compare(username, "admin") == 0) && (strings.Compare(password, "123456") == 0)) {
		http.Redirect(w, r, "/login", 301)
		return
	}

	var user User
	user.Username = username
	user.Secret = password;

	session := GetSession(r)

	// Set some session values.
	session.Values["user"] = user

	// Save it before we write to the response/return from the handler.
	session.Save(r, w)

	// If the login succeeded
	//res.Write([]byte("Hello " + databaseUsername))
	http.Redirect(w, r, "/", http.StatusFound)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Accel-Expires", "0")

	session := GetSession(r)
	session.Values["user"] = nil

	session.Save(r, w)

	http.Redirect(w, r, "/", 302)
}