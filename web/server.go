package web

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
	"github.com/baabeetaa/glogchain/config"
//	"github.com/baabeetaa/glogchain/protocol"
	//"github.com/baabeetaa/glogchain/blog"
	"github.com/tendermint/tmsp/client"
	"encoding/hex"
	//"github.com/tendermint/tmsp/types"
)

// simple web server
// ref https://reinbach.com/golang-webapps-1.html

const STATIC_URL string = "/static/"
const STATIC_ROOT string = "web/static/"

type Context struct {
	Title  string
	Static string
}

func Home(w http.ResponseWriter, req *http.Request) {
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

	client, err := tmspcli.NewGRPCClient(config.GlogchainConfigGlobal.TmspAddr, false)
	client.Start()
	defer client.Stop()

	if err == nil {
		txBytes := stringOrHexToBytes(newPostString)
		res := client.AppendTxSync(txBytes)

		if ( !res.IsOK()) {
			log.Print("PostCreateSave AppendTxSync error: ", res.Code)
		}
	} else {
		log.Print("PostCreateSave error: ", err)
	}

	// delay sometime to make sure hugo loading new page content
	time.Sleep(1000 * time.Millisecond) // 1s

	http.Redirect(w, req, "http://localhost:1313/post/" + Title  + "/", http.StatusFound)
}

func render(w http.ResponseWriter, tmpl string, context Context) {
	context.Static = STATIC_URL
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

func StaticHandler(w http.ResponseWriter, req *http.Request) {
	static_file := req.URL.Path[len(STATIC_URL):]
	if len(static_file) != 0 {
		f, err := http.Dir(STATIC_ROOT).Open(static_file)
		if err == nil {
			content := io.ReadSeeker(f)
			http.ServeContent(w, req, static_file, time.Now(), content)
			return
		}
	}
	http.NotFound(w, req)
}

func StartWebServer() error  {
	http.HandleFunc("/", Home)
	http.HandleFunc("/about/", About)
	http.HandleFunc("/post/create", PostCreate)
	http.HandleFunc("/post/create/save", PostCreateSave)
	http.HandleFunc(STATIC_URL, StaticHandler)
	err := http.ListenAndServe(config.GlogchainConfigGlobal.GlogchainWebAddr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	return nil
}

func main() {
	StartWebServer()
}


////////////
// NOTE: s is interpreted as a string unless prefixed with 0x
func stringOrHexToBytes(s string) []byte {
	if len(s) > 2 && s[:2] == "0x" {
		b, err := hex.DecodeString(s[2:])
		if err != nil {
			fmt.Println("Error decoding hex argument:", err.Error())
		}
		return b
	}
	return []byte(s)
}

