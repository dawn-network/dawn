package web

// base on simple web server - ref https://reinbach.com/golang-webapps-1.html
// use https://github.com/gorilla/mux

import (
	//"github.com/gorilla/sessions"
	"net/http"
	"log"
	"time"
	"github.com/baabeetaa/glogchain/app"
	"github.com/gorilla/mux"
	//"encoding/gob"
	//"github.com/tendermint/go-crypto"
	//"github.com/baabeetaa/glogchain/db"
	"github.com/gorilla/sessions"
	//"github.com/gorilla/securecookie"
	"encoding/gob"
	"github.com/tendermint/go-crypto"
	"github.com/baabeetaa/glogchain/db"
)

type Context struct {
	Title  		string
	Static 		string
	//Request 	*http.Request
	SessionValues 	map[interface{}]interface{}
	Data 		interface{}
}

type ActionResult struct {
	Status 		string 		// success or error
	Message 	string
	Data 		interface{}
}

var store *sessions.CookieStore
//var store = sessions.NewCookieStore(securecookie.GenerateRandomKey(16), securecookie.GenerateRandomKey(16))

//func HomeHandler(w http.ResponseWriter, req *http.Request) {
//	context := Context{Title: "Welcome!"}
//	context.Static = "/static/"
//	render(w, "home", context)
//}

func StartWebServer() error  {
	store = sessions.NewCookieStore([]byte("something-very-secret"))
	store.Options = &sessions.Options{
		//Domain:   "localhost",
		Path:     "/",
		MaxAge:   3600 * 8, // 8 hours
		//HttpOnly: true,
	}

	r := mux.NewRouter()

	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(app.GlogchainConfigGlobal.WebRootDir + "/web/static/"))))

	//r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/", CategoryHandler)
	r.HandleFunc("/login", LoginHandler)
	r.HandleFunc("/logout", LogoutHandler)
	r.HandleFunc("/account/create", AccountCreate)
	r.HandleFunc("/account/generate_keypair", GenerateKeyPair)
	//r.HandleFunc("/category", CategoryHandler)
	r.HandleFunc("/post", ViewSinglePostHandler)
	r.HandleFunc("/post/create", AuthWrapper(PostCreateHandler))
	r.HandleFunc("/post/edit", AuthWrapper(PostEditHandler))
	r.HandleFunc("/wallet/view", AuthWrapper(WalletViewHandler))
	r.HandleFunc("/wallet/sendtoken", AuthWrapper(WalletSendTokenHandler))

	// Subrouter
	//s := r.PathPrefix("/secure").Subrouter()
	//// /secure/test
	//s.HandleFunc("/test", AuthWrapper(HomeHandler))

	srv := &http.Server{
		Handler:      r,
		Addr:         app.GlogchainConfigGlobal.GlogchainWebAddr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe()) // Bind to a port and pass our router in

	return nil
}

func init()  {
	// to store a complex datatype within a session
	gob.Register(crypto.PrivKeyEd25519 {})
	gob.Register(db.User {})
}