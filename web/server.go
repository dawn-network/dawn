package web

// base on simple web server - ref https://reinbach.com/golang-webapps-1.html
// use https://github.com/gorilla/mux

import (
	"github.com/gorilla/sessions"
	"net/http"
	"fmt"
	"encoding/hex"
	"log"
	"time"
	"github.com/baabeetaa/glogchain/config"
	"io/ioutil"
	"github.com/gorilla/mux"
	"encoding/gob"
	"github.com/baabeetaa/glogchain/db"
	"strings"
	"github.com/tendermint/go-crypto"
	"encoding/binary"
	"bytes"
	"github.com/baabeetaa/glogchain/protocol"
	"encoding/json"
	"golang.org/x/crypto/ripemd160"
	"github.com/baabeetaa/glogchain/service"
)

type Context struct {
	Title  		string
	Static 		string
	Data 		interface{}
}

type ActionResult struct {
	Status 		string 		// success or error
	Message 	string
	Data 		interface{}
}

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func HomeHandler(w http.ResponseWriter, req *http.Request) {
	context := Context{Title: "Welcome!"}
	context.Static = "/static/"
	render(w, "home", context)
}

func CategoryHandler(w http.ResponseWriter, req *http.Request) {
	//context := Context{Title: "Welcome!"}
	//context.Static = "/static/"

	cat := req.FormValue("cat") // category id
	//categoryId, err := strconv.ParseInt(cat, 10, 64)
	//if err != nil {
	//	panic(err)
	//}

	posts, err := db.GetPostsByCategory(cat, 0, 20)
	if err != nil {
		log.Println("CategoryHandler", err)
		panic(err)
	}

	render(w, "category", posts)
}

func ViewSinglePostHandler(w http.ResponseWriter, req *http.Request) {
	//context := Context{Title: "Welcome!"}
	//context.Static = "/static/"

	p := req.FormValue("p")

	//postId, err := strconv.ParseInt(p, 10, 64)
	//if err != nil {
	//	panic(err)
	//}

	post, err := db.GetPost(p)
	if err != nil {
		panic(err)
	}

	//context.Data = postId
	//log.Println("ViewSinglePostHandler: " + (context.Data.(db.Post)).PostTitle )

	render(w, "single_post", post)
}

func AccountCreateView(w http.ResponseWriter, req *http.Request) {
	context := Context{Title: "Welcome!"}
	context.Static = "/static/"
	context.Data = map[string]interface{}{ "username": "", "pubkey": ""}
	render(w, "account_create", context)
}

func AccountCreateHandler(w http.ResponseWriter, req *http.Request) {
	username := req.FormValue("username")
	pubkey := req.FormValue("pubkey")

	log.Println("AccountCreateHandler", "username", username)
	log.Println("AccountCreateHandler", "pubkey", pubkey)

	if (len(username) < 6) {
		render(w, "account_create", ActionResult{Status: "error", Message: "username must be at least 6 characters", Data: map[string]interface{}{ "username": username, "pubkey": pubkey}})
		return
	}

	if (len(pubkey) != 64) {
		render(w, "account_create", ActionResult{Status: "error", Message: "PubKey must be 32 bytes in Hex String ( 64 characters)", Data: map[string]interface{}{ "username": username, "pubkey": pubkey}})
		return
	}


	pubkey = strings.ToUpper(pubkey)
	byte_arr, err := hex.DecodeString(pubkey)
	if (err != nil) {
		render(w, "account_create", ActionResult{Status: "error", Message: err.Error(), Data: map[string]interface{}{ "username": username, "pubkey": pubkey}})
		return
	}

	buf := &bytes.Buffer{}
	err = binary.Write(buf, binary.BigEndian, byte_arr)
	if err != nil {
		render(w, "account_create", ActionResult{Status: "error", Message: err.Error(), Data: map[string]interface{}{ "username": username, "pubkey": pubkey}})
		return
	}

	var key crypto.PubKeyEd25519
	binary.Read(buf, binary.BigEndian, &key);

	log.Println("AccountCreateHandler", "key", key.KeyString())
	var address string = strings.ToUpper(hex.EncodeToString(key.Address()))
	log.Println("AccountCreateHandler Address=\t\t" + address)


	//////////////////////////////////////
	var user db.User
	user.ID = address
	user.Username = username
	user.Pubkey = pubkey
	user.UserRegistered = "2017-01-06 09:00:28"
	user.DisplayName = username

	tx := protocol.OperationEnvelope{ Type: "AccountCreateOperation", Operation: protocol.AccountCreateOperation{ Fee: 0, User: user }}
	byte_arr, err = json.Marshal(tx)
	if err != nil {
		log.Fatal("AccountCreateHandler Cannot encode to JSON ", err)
	}

	tx_json := string(byte_arr[:])
	log.Println("AccountCreateHandler tx_json=", tx_json)

	tx_json_hex := make([]byte, len(tx_json) * 2)
	hex.Encode(tx_json_hex, []byte(tx_json))
	log.Println("AccountCreateHandler tx_json_hex", string(tx_json_hex[:]))


	var url_request string = config.GlogchainConfigGlobal.TmRpcLaddr + "/broadcast_tx_commit?tx=%22" + string(tx_json_hex[:]) + "%22"
	log.Println("AccountCreateHandler url_request: %#v\n", url_request)
	resp, err := http.Get(url_request)
	if err != nil {
		log.Println("AccountCreateHandler http.Get error: %#v\n", err)
		return;
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("AccountCreateHandler ioutil.ReadAll error: %#v\n", err)
		return;
	}
	json_response_string := string(body[:])
	log.Println("AccountCreateHandler json_response_string: %#v\n", json_response_string)


	render(w, "account_create", ActionResult{Status: "success", Message: "ok", Data: map[string]interface{}{ "username": username, "pubkey": pubkey}})
}

func AboutHandler(w http.ResponseWriter, req *http.Request) {
	context := Context{Title: "About"}
	render(w, "about", context)
}

func PostCreateHandler(w http.ResponseWriter, req *http.Request) {
	log.Println("---PostCreateHandler---------------------------------------------------------------------")

	var post db.Post

	// If method is GET serve an html login page
	if req.Method != "POST" {
		log.Println("PostCreateHandler GET")
		context := Context{Title: "PostCreate", Data: post}
		render(w, "PostCreate", context)
		return
	}

	log.Println("PostCreateHandler POST")

	session := GetSession(req)
	user, ok := session.Values["user"].(*db.User)
	if !ok {
		panic("PostCreateHandler: wtf?? session user is nil")
		return
	}

	timeNow := time.Now()

	/////////
	post.ID = "" // hex string generate by ripemd160 hash the post, will be assinged latter
	post.PostAuthor = user.ID 				// take from session
	post.PostDate = timeNow.String() 			// take from current datetime
	post.PostContent = req.FormValue("PostContent")
	post.PostTitle = req.FormValue("Title")
	post.PostModified = timeNow.String() 		// take from current datetime
	post.Thumb = req.FormValue("Thumb")
	post.Cat = req.FormValue("Categories") 	// Categories in json string array

	// very basic from validating
	if ((len(post.PostTitle) < 6) || (len(post.PostContent) < 6) || (len(post.Thumb) < 6)) {
		render(w, "PostCreate",
			ActionResult{
				Status: "error",
				Message: "field must be at least 6 characters",
				Data: post,
			})
		return
	}

	// validating Categories
	post.Cat = strings.ToLower(post.Cat)
	if (len(post.Cat) < 6) {
		post.Cat = `[]`
	}

	cats_string := []string{}
	err := json.Unmarshal([]byte(post.Cat), &cats_string)
	if (err != nil) {
		render(w, "PostCreate",
			ActionResult{
				Status: "error",
				Message: "Categories json array string is invalid",
				Data: post,
			})
		return
	}

	// generate id of new post
	hasher := ripemd160.New()
	str_mix := fmt.Sprint("%s | %s | %s | %s", post.PostAuthor, post.PostTitle, post.PostDate, service.RandSeq(8))
	hasher.Write([]byte(str_mix))
	buf_id := hasher.Sum(nil)
	post.ID = strings.ToUpper(hex.EncodeToString(buf_id))

	log.Println("PostCreateHandler", "id=", post.ID, "PostAuthor=", post.PostAuthor, "PostDate=",
		post.PostDate, "Title=", post.PostTitle, "PostModified=", post.PostModified, "Thumb=",
		post.Thumb, "Categories=", post.Cat, "PostContent=", post.PostContent)


	tx := protocol.OperationEnvelope {
		Type: "PostCreateOperation",
		Operation: protocol.PostCreateOperation {
			Fee: 0,
			Post: post },
	}

	byte_arr, err := json.Marshal(tx)
	if err != nil {
		log.Fatal("PostCreateHandler Cannot encode to JSON ", err)
		return
	}

	tx_json := string(byte_arr[:])
	log.Println("PostCreateHandler tx_json=", tx_json)

	tx_json_hex := make([]byte, len(tx_json) * 2)
	hex.Encode(tx_json_hex, []byte(tx_json))
	log.Println("PostCreateHandler tx_json_hex", string(tx_json_hex[:]))


	var url_request string = config.GlogchainConfigGlobal.TmRpcLaddr + "/broadcast_tx_commit?tx=%22" + string(tx_json_hex[:]) + "%22"
	log.Println("PostCreateHandler url_request: %#v\n", url_request)
	resp, err := http.Get(url_request)
	if err != nil {
		log.Println("PostCreateHandler http.Get error: %#v\n", err)
		return;
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("PostCreateHandler ioutil.ReadAll error: %#v\n", err)
		return;
	}
	json_response_string := string(body[:])
	log.Println("PostCreateHandler json_response_string: %#v\n", json_response_string)

	// delay sometime then Redirect to the new post
	time.Sleep(1000 * time.Millisecond) // 1s
	http.Redirect(w, req, "/post?p=" + post.ID  , http.StatusFound)
}

func StartWebServer() error  {
	gob.Register(&db.User{})
	gob.Register(&crypto.PrivKeyEd25519{})


	r := mux.NewRouter()

	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	//r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/", CategoryHandler)
	r.HandleFunc("/login", LoginHandler)
	r.HandleFunc("/logout", LogoutHandler)
	r.HandleFunc("/account/create", AccountCreateView)
	r.HandleFunc("/account/create_handler", AccountCreateHandler)
	//r.HandleFunc("/category", CategoryHandler)
	r.HandleFunc("/post", ViewSinglePostHandler)
	r.HandleFunc("/about/", AuthWrapper(AboutHandler))
	r.HandleFunc("/post/create", AuthWrapper(PostCreateHandler))

	// Subrouter
	//s := r.PathPrefix("/secure").Subrouter()
	//// /secure/test
	//s.HandleFunc("/test", AuthWrapper(HomeHandler))

	srv := &http.Server{
		Handler:      r,
		Addr:         config.GlogchainConfigGlobal.GlogchainWebAddr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe()) // Bind to a port and pass our router in

	return nil
}
