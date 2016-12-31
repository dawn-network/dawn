package web

import (
	"net/http"
	"log"
	"strings"
	"time"
)

type User struct {
	Username string
	Id string
	Secret string
}

// https://medium.com/@matryer/the-http-handlerfunc-wrapper-technique-in-golang-c60bf76e6124#.ipqttmx8e
func AuthWrapper(fn http.HandlerFunc) http.HandlerFunc {
	// called once per wrapping
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("MyWrapper")

		// Get a session. We're ignoring the error resulted from decoding an
		// existing session: Get() always returns a session, even if empty.
		session, err := store.Get(r, "session-name")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		user, ok := session.Values["user"].(*User)
		if !ok {
			LoginHandler(w, r)
			return
		}
		log.Println("AuthWrapper user.Username:", user.Username)

		fn(w, r)
	}
}


func LoginHandler(w http.ResponseWriter, req *http.Request) {
	// If method is GET serve an html login page
	if req.Method != "POST" {
		http.ServeFile(w, req, "web/templates/login.html")
		return
	}

	// Grab the username/password from the submitted post form
	username := req.FormValue("username")
	password := req.FormValue("password")

	// If wrong password redirect to the login
	if !((strings.Compare(username, "admin") == 0) && (strings.Compare(password, "123456") == 0)) {
		http.Redirect(w, req, "/login", 301)
		return
	}

	var user User
	user.Username = username
	user.Secret = password;

	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, err := store.Get(req, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Set some session values.
	session.Values["user"] = user

	// Save it before we write to the response/return from the handler.
	session.Save(req, w)


	// If the login succeeded
	//res.Write([]byte("Hello " + databaseUsername))
	http.Redirect(w, req, "/", http.StatusFound)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Accel-Expires", "0")

	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, err := store.Get(r, "session-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	session.Values["user"] = nil

	session.Save(r, w)

	//cookie := &http.Cookie{
	//	 Name:   "session-name",
	//	 Value:  "",
	//	 Path:   "/",
	//	 MaxAge: -1,
	//     }
	//http.SetCookie(w, cookie)


	http.Redirect(w, r, "/", 302)
}