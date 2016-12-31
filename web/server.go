package web

import (
	"fmt"
	"html/template"
	//"io"
	"log"
	"net/http"
	"time"
	"github.com/baabeetaa/glogchain/config"
//	"github.com/baabeetaa/glogchain/protocol"
	//"github.com/baabeetaa/glogchain/blog"
	//"github.com/tendermint/tmsp/client"
	"encoding/hex"
	//"github.com/tendermint/tmsp/types"
	"io/ioutil"
	//"net/url"
	"strings"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

// simple web server
// ref https://reinbach.com/golang-webapps-1.html

//const STATIC_URL string = "/static/"
//const STATIC_ROOT string = "web/static/"

type Context struct {
	Title  string
	Static string
}

var store = sessions.NewCookieStore([]byte("something-very-secret"))

func HomeHandler(w http.ResponseWriter, req *http.Request) {
	context := Context{Title: "Welcome!"}
	render(w, "index", context)
}

func About(w http.ResponseWriter, req *http.Request) {
	context := Context{Title: "About"}
	render(w, "about", context)
}

func PostCreate(w http.ResponseWriter, req *http.Request) {
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

	//client, err := tmspcli.NewGRPCClient(config.GlogchainConfigGlobal.TmspAddr, false)
	//client.Start()
	//defer client.Stop()
	//
	//if err == nil {
	//	txBytes := stringOrHexToBytes(newPostString)
	//	res := client.AppendTxSync(txBytes)
	//
	//	if ( !res.IsOK()) {
	//		log.Print("PostCreateSave AppendTxSync error: ", res.Code)
	//	}
	//
	//	client.CommitSync()
	//} else {
	//	log.Print("PostCreateSave error: ", err)
	//}

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

	// If the login succeeded
	//res.Write([]byte("Hello " + databaseUsername))
	http.Redirect(w, req, "/", http.StatusFound)
}


func render(w http.ResponseWriter, tmpl string, context Context) {
	context.Static = "/static/"
	tmpl_list := []string{"web/templates/base.html", fmt.Sprintf("web/templates/%s.html", tmpl)}
	t, err := template.ParseFiles(tmpl_list...)
	if err != nil {
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w, context)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

//func StaticHandler(w http.ResponseWriter, req *http.Request) {
//	static_file := req.URL.Path[len(STATIC_URL):]
//	if len(static_file) != 0 {
//		f, err := http.Dir(STATIC_ROOT).Open(static_file)
//		if err == nil {
//			content := io.ReadSeeker(f)
//			http.ServeContent(w, req, static_file, time.Now(), content)
//			return
//		}
//	}
//	http.NotFound(w, req)
//}

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

		// Retrieve our struct and type-assert it
		val := session.Values["person"]

		if (val == nil) {
			LoginHandler(w, r)
			return
		}

		//var person = &Person{}
		//if person, ok := val.(*Person); !ok {
		//	// Handle the case that it's not an expected type
		//}
		//
		// Now we can use our person object

		fn(w, r)
	}
}

func StartWebServer() error  {
	//http.HandleFunc("/", HomeHandler)
	//http.HandleFunc("/login", LoginHandler)
	//http.HandleFunc("/about/", About)
	//http.HandleFunc("/post/create", PostCreate)
	//http.HandleFunc("/post/create/save", PostCreateSave)
	//http.HandleFunc(STATIC_URL, StaticHandler)
	//err := http.ListenAndServe(config.GlogchainConfigGlobal.GlogchainWebAddr, nil)
	//if err != nil {
	//	log.Fatal("ListenAndServe: ", err)
	//}

	// use https://github.com/gorilla/mux
	r := mux.NewRouter()

	// This will serve files under http://localhost:8000/static/<filename>
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static/"))))

	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/login", LoginHandler)
	//r.HandleFunc("/products", ProductsHandler)
	//r.HandleFunc("/articles", ArticlesHandler)
	//http.Handle("/", r)

	s := r.PathPrefix("/secure").Subrouter()
	// /secure/test
	s.HandleFunc("/test", AuthWrapper(HomeHandler))

	srv := &http.Server{
		Handler:      r,
		Addr:         config.GlogchainConfigGlobal.GlogchainWebAddr,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Bind to a port and pass our router in
	log.Fatal(srv.ListenAndServe())

	return nil
}

//func main() {
//	StartWebServer()
//}


////////////
//// NOTE: s is interpreted as a string unless prefixed with 0x
//func stringOrHexToBytes(s string) []byte {
//	if len(s) > 2 && s[:2] == "0x" {
//		b, err := hex.DecodeString(s[2:])
//		if err != nil {
//			fmt.Println("Error decoding hex argument:", err.Error())
//		}
//		return b
//	}
//	return []byte(s)
//}

