package web

import (
	"fmt"
	"html/template"
	"net/http"
	"log"
	"github.com/dawn-network/glogchain/core/db"
	"io"
)

var funcMap template.FuncMap = template.FuncMap {
	//"GetFeaturedPosts": 	db.GetFeaturedPosts,
	"GetPost": 		db.GetPost,
	"GetCategoryOfPost": 	db.GetCategoryOfPost,
	//"GetPostThumbnail": 	db.GetPostThumbnail,
	"GetCommentsByPost":	db.GetCommentsByPost,
	"GetUser": 		db.GetUser,
	"GetPostsByCategory": 	db.GetPostsByCategory,
	"GetTopCategories":	db.GetTopCategories,
	"GetType": 		GetType,
	"Dict": 		Dict,
	"StringCut": 		StringCut,
	"Config_IpFsGateway": 	Config_IpFsGateway }


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
	context := Context{Title: "Welcome!"}
	context.Static = "/webcontent/static/"
	//context.Request = req
	context.SessionValues = GetSession(req).Values

	p := req.FormValue("p")
	post, err := db.GetPost(p)
	if err != nil {
		panic(err)
	}

	context.Data = post
	render(w, "single_post", context)
}


func render(w http.ResponseWriter, tmpl string, data interface{}) {
	//context := Context{Title: "Welcome!"}
	//context.Static = "/static/"

	t := template.New("index.html")
	t = t.Funcs(funcMap)


	//////////////////////////
	// using file system
	//var tmpl_list = []string {
	//	app.GlogchainConfigGlobal.WebRootDir + "/web/webcontent/templates/index.html",
	//	app.GlogchainConfigGlobal.WebRootDir + "/web/webcontent/templates/header.html",
	//	app.GlogchainConfigGlobal.WebRootDir + "/web/webcontent/templates/footer.html",
	//	app.GlogchainConfigGlobal.WebRootDir + "/web/webcontent/templates/featured_posts.html",
	//	app.GlogchainConfigGlobal.WebRootDir + "/web/webcontent/templates/highlighted_posts.html",
	//	app.GlogchainConfigGlobal.WebRootDir + "/web/webcontent/templates/primary.html",
	//	app.GlogchainConfigGlobal.WebRootDir + "/web/webcontent/templates/secondary.html",
	//	app.GlogchainConfigGlobal.WebRootDir + "/web/webcontent/templates/widget_slider.html",
	//	app.GlogchainConfigGlobal.WebRootDir + "/web/webcontent/templates/widget_featured_posts_vertical.html",
	//	app.GlogchainConfigGlobal.WebRootDir + "/web/webcontent/templates/widget_featured_posts.html",
	//	app.GlogchainConfigGlobal.WebRootDir + "/web/webcontent/templates/widget_728x90_advertisement.html",
	//	fmt.Sprintf(app.GlogchainConfigGlobal.WebRootDir + "/web/webcontent/templates/%s.html", tmpl)}
	//
	//t, err := t.ParseFiles(tmpl_list...)
	//
	//if err != nil {
	//	log.Print("template parsing error: ", err)
	//}

	//////////////////////////
	// using bundle resources
	prefix_path := "webcontent/templates/"

	var tmpl_list = []string {
		prefix_path + "index.html",
		prefix_path + "header.html",
		prefix_path + "footer.html",
		prefix_path + "featured_posts.html",
		prefix_path + "highlighted_posts.html",
		prefix_path + "primary.html",
		prefix_path + "secondary.html",
		prefix_path + "widget_slider.html",
		prefix_path + "widget_featured_posts_vertical.html",
		prefix_path + "widget_featured_posts.html",
		prefix_path + "widget_728x90_advertisement.html",
		fmt.Sprintf(prefix_path + "%s.html", tmpl)}

	for _, v := range tmpl_list {
		source, _ := Asset(v)
		t, _ = t.Parse(string(source))
	}


	err := t.Execute(w, data)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}

/////////////////////////////////////////

func serve400(w http.ResponseWriter) {
	w.WriteHeader(http.StatusBadRequest)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, http.StatusText(http.StatusBadRequest))
}

func serve404(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNotFound)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	io.WriteString(w, http.StatusText(http.StatusNotFound))
}

