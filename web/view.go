package web

import (
	"fmt"
	"html/template"
	"net/http"
	"log"
	"errors"
	"github.com/baabeetaa/glogchain/db"
)

var funcMap template.FuncMap = template.FuncMap{
	"GetFeaturedPosts": 	db.GetFeaturedPosts,
	"GetPost": 		db.GetPost,
	"GetCategoryOfPost": 	db.GetCategoryOfPost,
	"GetPostThumbnail": 	db.GetPostThumbnail,
	"GetUser": 		db.GetUser,
	"GetPostsByCategory": 	db.GetPostsByCategory,
	"Dict": 		Dict,
	"StringCut": 		StringCut}

//var templates map[string]*template.Template
//
//// Load templates on program initialisation
//func init() {
//	if templates == nil {
//		templates = make(map[string]*template.Template)
//	}
//
//	templatesDir := "tmpl"
//
//	layouts, err := filepath.Glob(templatesDir + "layouts/*.tmpl")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	includes, err := filepath.Glob(templatesDir + "includes/*.tmpl")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	// Generate our templates map from our layouts/ and includes/ directories
//	for _, layout := range layouts {
//		files := append(includes, layout)
//		templates[filepath.Base(layout)] = template.Must(template.ParseFiles(files...))
//	}
//}
//
//// renderTemplate is a wrapper around template.ExecuteTemplate.
//func renderTemplate(w http.ResponseWriter, name string, data map[string]interface{}) error {
//	// Ensure the template exists in the map.
//	tmpl, ok := templates[name]
//	if !ok {
//		return fmt.Errorf("The template %s does not exist.", name)
//	}
//
//	w.Header().Set("Content-Type", "text/html; charset=utf-8")
//	return tmpl.ExecuteTemplate(w, "base", data)
//}

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
	firstnChars := str[:n]
	return firstnChars
}

func render(w http.ResponseWriter, tmpl string, data interface{}) {
	//context := Context{Title: "Welcome!"}
	//context.Static = "/static/"

	t := template.New("index.html")
	t = t.Funcs(funcMap)

	var tmpl_list = []string {
		"web/templates/index.html",
		"web/templates/header.html",
		"web/templates/footer.html",
		"web/templates/featured_posts.html",
		"web/templates/highlighted_posts.html",
		"web/templates/primary.html",
		"web/templates/secondary.html",
		"web/templates/widget_slider.html",
		"web/templates/widget_featured_posts_vertical.html",
		"web/templates/widget_featured_posts.html",
		"web/templates/widget_728x90_advertisement.html",
		fmt.Sprintf("web/templates/%s.html", tmpl)}

	t, err := t.ParseFiles(tmpl_list...)
	if err != nil {
		log.Print("template parsing error: ", err)
	}

	err = t.Execute(w, data)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}