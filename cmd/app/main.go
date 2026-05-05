package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/joho/godotenv"

	"partyinvites.org/internal/domain/entity"
)

// responses holds the list of RSVP responses.
// slice empty but has space for 10 elements before needing to resize.
var responses = make([]*entity.RSVP, 0, 10)

// templates holds the parsed HTML templates for the application.
var templates = make(map[string]*template.Template, 3)

func main() {
	loadTemplates()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the port from the environment variable, default to 8080 if not set.
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/list", listHandler)
	http.HandleFunc("/form", formHandler)

	fmt.Println("Starting server on :", port)
	err = http.ListenAndServe(":"+port, nil)
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

// formHandler serves the RSVP form when the "/form" URL is accessed with a GET request.
func formHandler(writer http.ResponseWriter, request *http.Request) {
	if request.Method == http.MethodGet {
		templates["form"].Execute(writer, entity.FormData{
			RSVP:   &entity.RSVP{},
			Errors: []string{},
		})
	} else if request.Method == http.MethodPost {
		request.ParseForm()
		responseData := &entity.RSVP{
			Name:       request.Form.Get("name"),
			Email:      request.Form.Get("email"),
			Phone:      request.Form.Get("phone"),
			WillAttend: request.Form.Get("willattend") == "true",
		}

		errors := []string{}
		if responseData.Name == "" {
			errors = append(errors, "Please enter your name")
		}
		if responseData.Email == "" {
			errors = append(errors, "Please enter your email address")
		}
		if responseData.Phone == "" {
			errors = append(errors, "Please enter your phone number")
		}
		if len(errors) > 0 {
			templates["form"].Execute(writer, entity.FormData{
				RSVP:   responseData,
				Errors: errors,
			})
			return
		}

		responses = append(responses, responseData)

		if responseData.WillAttend {
			templates["thanks"].Execute(writer, responseData)
		} else {
			templates["sorry"].Execute(writer, responseData)
		}
	}
}

// rootPath returns the root path of the project by using runtime.
// Caller to get the current file's location and navigating up the directory structure.
func rootPath(path string) string {
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "../..")
	fmt.Println("Project root:", projectRoot)
	return filepath.Join(projectRoot, path)
}
