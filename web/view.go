package web

import (
	"fmt"
	"html/template"
	"net/http"
	"log"
	"github.com/dawn-network/glogchain/app"
	"github.com/dawn-network/glogchain/db"
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
	context.Static = app.GlogchainConfigGlobal.WebRootDir + "/static/"
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

	var tmpl_list = []string {
		app.GlogchainConfigGlobal.WebRootDir + "/web/templates/index.html",
		app.GlogchainConfigGlobal.WebRootDir + "/web/templates/header.html",
		app.GlogchainConfigGlobal.WebRootDir + "/web/templates/footer.html",
		app.GlogchainConfigGlobal.WebRootDir + "/web/templates/featured_posts.html",
		app.GlogchainConfigGlobal.WebRootDir + "/web/templates/highlighted_posts.html",
		app.GlogchainConfigGlobal.WebRootDir + "/web/templates/primary.html",
		app.GlogchainConfigGlobal.WebRootDir + "/web/templates/secondary.html",
		app.GlogchainConfigGlobal.WebRootDir + "/web/templates/widget_slider.html",
		app.GlogchainConfigGlobal.WebRootDir + "/web/templates/widget_featured_posts_vertical.html",
		app.GlogchainConfigGlobal.WebRootDir + "/web/templates/widget_featured_posts.html",
		app.GlogchainConfigGlobal.WebRootDir + "/web/templates/widget_728x90_advertisement.html",
		fmt.Sprintf(app.GlogchainConfigGlobal.WebRootDir + "/web/templates/%s.html", tmpl)}

	t, err := t.ParseFiles(tmpl_list...)
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}