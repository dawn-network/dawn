package web

import (
	"fmt"
	"html/template"
	"net/http"
	"log"
	"errors"
	"github.com/baabeetaa/glogchain/db"
	"reflect"
	"github.com/baabeetaa/glogchain/config"
)

var funcMap template.FuncMap = template.FuncMap{
	//"GetFeaturedPosts": 	db.GetFeaturedPosts,
	"GetPost": 		db.GetPost,
	"GetCategoryOfPost": 	db.GetCategoryOfPost,
	//"GetPostThumbnail": 	db.GetPostThumbnail,
	"GetUser": 		db.GetUser,
	"GetPostsByCategory": 	db.GetPostsByCategory,
	"GetTopCategories":	db.GetTopCategories,
	"GetType": 		GetType,
	"Dict": 		Dict,
	"StringCut": 		StringCut}



// http://stackoverflow.com/questions/20170275/how-to-find-a-type-of-a-object-in-golang
func GetType(v interface{}) string {
	return reflect.TypeOf(v).String()
}

// http://stackoverflow.com/questions/18276173/calling-a-template-with-several-pipeline-parameters
func Dict(values ...interface{}) (map[string]interface{}, error) {
	if len(values)%2 != 0 {
		return nil, errors.New("invalid dict call")
	}

	dict := make(map[string]interface{}, len(values)/2)

	for i := 0; i < len(values); i+=2 {
		key, ok := values[i].(string)
		if !ok {
			return nil, errors.New("dict keys must be strings")
		}
		dict[key] = values[i+1]

		// fix wrong type for value; expected int64; got int
		switch dict[key].(type) {
		case int:
			dict[key] = int64(dict[key].(int))
		}
	}

	return dict, nil
}

func StringCut(str string, n int) string {
	if (n < len(str)) {
		return str[:n]
	}

	return str
}

func render(w http.ResponseWriter, tmpl string, data interface{}) {
	//context := Context{Title: "Welcome!"}
	//context.Static = "/static/"

	t := template.New("index.html")
	t = t.Funcs(funcMap)

	var tmpl_list = []string {
		config.GlogchainConfigGlobal.WebRootDir + "/web/templates/index.html",
		config.GlogchainConfigGlobal.WebRootDir + "/web/templates/header.html",
		config.GlogchainConfigGlobal.WebRootDir + "/web/templates/footer.html",
		config.GlogchainConfigGlobal.WebRootDir + "/web/templates/featured_posts.html",
		config.GlogchainConfigGlobal.WebRootDir + "/web/templates/highlighted_posts.html",
		config.GlogchainConfigGlobal.WebRootDir + "/web/templates/primary.html",
		config.GlogchainConfigGlobal.WebRootDir + "/web/templates/secondary.html",
		config.GlogchainConfigGlobal.WebRootDir + "/web/templates/widget_slider.html",
		config.GlogchainConfigGlobal.WebRootDir + "/web/templates/widget_featured_posts_vertical.html",
		config.GlogchainConfigGlobal.WebRootDir + "/web/templates/widget_featured_posts.html",
		config.GlogchainConfigGlobal.WebRootDir + "/web/templates/widget_728x90_advertisement.html",
		fmt.Sprintf(config.GlogchainConfigGlobal.WebRootDir + "/web/templates/%s.html", tmpl)}

	t, err := t.ParseFiles(tmpl_list...)
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}