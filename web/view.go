package web

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"log"
	"errors"
)

var templates map[string]*template.Template

// Load templates on program initialisation
func init() {
	if templates == nil {
		templates = make(map[string]*template.Template)
	}

	templatesDir := "tmpl"

	layouts, err := filepath.Glob(templatesDir + "layouts/*.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	includes, err := filepath.Glob(templatesDir + "includes/*.tmpl")
	if err != nil {
		log.Fatal(err)
	}

	// Generate our templates map from our layouts/ and includes/ directories
	for _, layout := range layouts {
		files := append(includes, layout)
		templates[filepath.Base(layout)] = template.Must(template.ParseFiles(files...))
	}
}

// renderTemplate is a wrapper around template.ExecuteTemplate.
func renderTemplate(w http.ResponseWriter, name string, data map[string]interface{}) error {
	// Ensure the template exists in the map.
	tmpl, ok := templates[name]
	if !ok {
		return fmt.Errorf("The template %s does not exist.", name)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	return tmpl.ExecuteTemplate(w, "base", data)
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
	firstnChars := str[:n]
	return firstnChars
}
