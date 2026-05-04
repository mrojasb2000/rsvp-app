package main

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"runtime"

	"partyinvites.org/internal/domain/entity"
)

// responses holds the list of RSVP responses.
// slice empty but has space for 10 elements before needing to resize.
var responses = make([]*entity.RSVP, 0, 10)

// templates holds the parsed HTML templates for the application.
var templates = make(map[string]*template.Template, 3)

func main() {
	loadTemplates()

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/list", listHandler)

	fmt.Println("Starting server on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

// loadTemplates parses the HTML templates and stores them in the templates map.
func loadTemplates() {
	templateNames := [5]string{"welcome", "form", "thanks", "sorry", "list"}
	pathRoot := rootPath("templates")
	layOutPath := fmt.Sprintf("%s/layout.html", pathRoot)

	for index, name := range templateNames {
		tmpl, err := template.ParseFiles(layOutPath, fmt.Sprintf("%s/%s.html", pathRoot, name))
		if err != nil {
			panic(err)
		}
		templates[name] = tmpl
		fmt.Println("Loaded template:", index, name)
	}
}

// welcomeHandler serves the welcome page when the root URL is accessed.
func welcomeHandler(writer http.ResponseWriter, request *http.Request) {
	templates["welcome"].Execute(writer, nil)
}

// listHandler serves the list of RSVP responses when the "/list" URL is accessed.
func listHandler(writer http.ResponseWriter, request *http.Request) {
	templates["list"].Execute(writer, responses)
}

// rootPath returns the root path of the project by using runtime.
// Caller to get the current file's location and navigating up the directory structure.
func rootPath(path string) string {
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "../..")
	fmt.Println("Raíz del proyecto:", projectRoot)
	return filepath.Join(projectRoot, path)
}
