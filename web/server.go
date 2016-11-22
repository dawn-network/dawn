package web

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"time"
	"glogchain/config"
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

