package web

// https://medium.com/@matryer/the-http-handlerfunc-wrapper-technique-in-golang-c60bf76e6124#.ipqttmx8e

import (
	"net/http"
	"strings"
	"time"
	"github.com/gorilla/sessions"
	"github.com/dawn-network/glogchain/db"
	"encoding/hex"
	"bytes"
	"encoding/binary"
	"github.com/tendermint/go-crypto"
	"log"
)

func GetSession(r *http.Request) (*sessions.Session) {
	// Get a session. We're ignoring the error resulted from decoding an
	// existing session: Get() always returns a session, even if empty.
	session, err := store.Get(r, "session-name")
	if (err != nil) {
		//http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Println("GetSession", err)
		panic("Unexpected error!")
	}

	return session
}

func AuthWrapper(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session := GetSession(r)

		_, ok := session.Values["user"].(db.User)
		if !ok {
			LoginHandler(w, r)
			return
		}

		fn(w, r)
	}
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// If method is GET serve an html login page
	if r.Method != "POST" {
		render(w, "login", nil)
		return
	}

	// Grab the username/password from the submitted post form
	str_PriKey := r.FormValue("PriKey")

	if (len(str_PriKey) != 128) {
		render(w, "login", ActionResult{Status: "error", Message: "PriKey must be 64 bytes in Hex String ( 128 characters)", Data: map[string]interface{}{ "PriKey": str_PriKey}})
		return
	}

	str_PriKey = strings.ToUpper(str_PriKey)

	byte_arr, err := hex.DecodeString(str_PriKey)
	if (err != nil) {
		render(w, "login", ActionResult{Status: "error", Message: err.Error(), Data: map[string]interface{}{ "PriKey": str_PriKey}})
		return
	}

	buf := &bytes.Buffer{}
	err = binary.Write(buf, binary.BigEndian, byte_arr)
	if err != nil {
		render(w, "login", ActionResult{Status: "error", Message: err.Error(), Data: map[string]interface{}{ "PriKey": str_PriKey}})
		return
	}

	var pri_key crypto.PrivKeyEd25519
	binary.Read(buf, binary.BigEndian, &pri_key);

	log.Println("LoginHandler", "pri_key", strings.ToUpper(hex.EncodeToString(pri_key.Bytes())))
	log.Println("LoginHandler PubKey KeyString=\t" + pri_key.PubKey().KeyString())

	var address string = strings.ToUpper(hex.EncodeToString(pri_key.PubKey().Address()))
	log.Println("LoginHandler Address=\t\t" + address)

	// find User in db
	user, err := db.GetUser(address)

	// If failed, redirect to the login
	if (err != nil) {
		render(w, "login", ActionResult{Status: "error", Message: err.Error(), Data: map[string]interface{}{ "PriKey": str_PriKey}})
		return
	}

	session := GetSession(r)

	session.Values["user"] = user // Set some session values.
	session.Values["private_key"] = pri_key
	err = session.Save(r, w) // Save it before we write to the response/return from the handler.
	if (err != nil) {
		log.Println("WTF - can not save session!!!")
	}

	//// test get session var
	//_, ok := session.Values["private_key"].(crypto.PrivKeyEd25519)
	//if (!ok) {
	//	log.Println("Can not get session private_key")
	//}
	_, ok := session.Values["user"].(db.User)
	if (!ok) {
		log.Println("WTF!!!")
	}

	http.Redirect(w, r, "/", http.StatusFound)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Cache-Control", "no-cache, private, max-age=0")
	w.Header().Set("Expires", time.Unix(0, 0).Format(http.TimeFormat))
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("X-Accel-Expires", "0")

	session := GetSession(r)
	session.Values["user"] = nil
	session.Values["private_key"] = nil
	session.Save(r, w)

	http.Redirect(w, r, "/", 302)
}