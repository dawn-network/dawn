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
	"strconv"
	"github.com/baabeetaa/glogchain/db"
	"strings"
	"github.com/tendermint/go-crypto"
	"encoding/binary"
	"bytes"
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
	categoryId, err := strconv.ParseInt(cat, 10, 64)
	if err != nil {
		panic(err)
	}

	posts, err := db.GetPostsByCategory(categoryId, 0, 20)
	if err != nil {
		panic(err)
	}

	render(w, "category", posts)
}

func ViewSinglePostHandler(w http.ResponseWriter, req *http.Request) {
	//context := Context{Title: "Welcome!"}
	//context.Static = "/static/"

	p := req.FormValue("p")

	postId, err := strconv.ParseInt(p, 10, 64)
	if err != nil {
		panic(err)
	}

	post, err := db.GetPost(postId)
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

	//key, err := byte_arr.(crypto.PubKeyEd25519)

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


// user.ID, user.Username, user.Pubkey, user.UserRegistered, user.DisplayName
	// save to db
	var user db.User
	user.ID = address
	user.Username = username
	user.Pubkey = pubkey
	user.UserRegistered = "2017-01-06 09:00:28"
	user.DisplayName = username

	err = db.CreateUser(user)
	if err != nil {
		render(w, "account_create", ActionResult{Status: "error", Message: err.Error(), Data: map[string]interface{}{ "username": username, "pubkey": pubkey}})
		return
	}


	render(w, "account_create", ActionResult{Status: "success", Message: "ok", Data: map[string]interface{}{ "username": username, "pubkey": pubkey}})
}

func AboutHandler(w http.ResponseWriter, req *http.Request) {
	context := Context{Title: "About"}
	render(w, "about", context)
}

func PostCreateHandler(w http.ResponseWriter, req *http.Request) {
	context := Context{Title: "PostCreate"}
	render(w, "PostCreate", context)
}

func PostCreateSave(w http.ResponseWriter, req *http.Request) {
	//var postOperation protocol.PostOperation
	Title := req.FormValue("Title")
	Author := req.FormValue("Author")
	Body := req.FormValue("Body")

	var newPostString string = "{\"Type\": \"PostOperation\" , \"Operation\" : {\"Title\": \"" +  Title + "\", \"Body\": \"" + Body + "\", \"Author\": \"" + Author + "\"} }"
	fmt.Printf("---------------------------------\n")
	fmt.Printf("newPostString: %s\n", newPostString)

	newPostStringHex := make([]byte, len([]byte(newPostString)) * 2)
	hex.Encode(newPostStringHex, []byte(newPostString))
	fmt.Println("newPostStringHex: %s\n", newPostStringHex)

	/// example
	// {"Type": "PostOperation" , "Operation" : {"Title": "aaa", "Body": "aaa", "Author": "aaa"} }
	// http://10.0.0.11:46657/broadcast_tx_commit?tx=%227b2254797065223a2022506f73744f7065726174696f6e22202c20224f7065726174696f6e22203a207b225469746c65223a2022616161222c2022426f6479223a2022616161222c2022417574686f72223a2022616161227d207d%22

	var url_request string = config.GlogchainConfigGlobal.TmRpcLaddr + "/broadcast_tx_commit?tx=%22" + string(newPostStringHex[:]) + "%22"
	log.Print("url_request: %#v\n", url_request)
	resp, err := http.Get(url_request)
	//req, err := http.NewRequest("GET", url_request, nil)
	//resp, err := http.PostForm(config.GlogchainConfigGlobal.TmRpcLaddr + "/broadcast_tx_commit", url.Values{"tx": {"%22" + string(newPostStringHex[:]) + "%22"}})
	if err != nil {
		log.Print("http.Get error: %#v\n", err)
		return;
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Print("ioutil.ReadAll error: %#v\n", err)
		return;
	}
	json_response_string := string(body[:])
	fmt.Println("json_response_string: %#v\n", json_response_string)

	// delay sometime to make sure hugo loading new page content
	time.Sleep(1000 * time.Millisecond) // 1s
	http.Redirect(w, req, config.GlogchainConfigGlobal.HugoBaseUrl + "/post/" + Title  + "/", http.StatusFound)
}



func StartWebServer() error  {
	gob.Register(&User{})

	r := mux.NewRouter()

	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/login", LoginHandler)
	r.HandleFunc("/logout", LogoutHandler)
	r.HandleFunc("/account/create", AccountCreateView)
	r.HandleFunc("/account/create_handler", AccountCreateHandler)

	r.HandleFunc("/category", CategoryHandler)
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
