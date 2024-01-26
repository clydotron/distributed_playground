package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
)

var webPort = "8090" // TODO env / command line
// var pathToTemplates = "./cmd/web/templates/"
var pathToTemplates = "templates/"

// simple web front end:
func main() {

	// initialize the templates?

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		render(w, "home.page.gohtml")
	})

	fmt.Printf("Starting front end service on port %s\n", webPort)
	err := http.ListenAndServe(fmt.Sprintf(":%s", webPort), nil)
	if err != nil {
		log.Panic(err)
	}
}

//go:embed templates
var templateFS embed.FS

func render(w http.ResponseWriter, templateToShow string) {

	// explore the other ways to do this:

	// use fmt.Sprintf instead?
	partials := []string{
		pathToTemplates + "base.layout.gohtml",
		pathToTemplates + "header.partial.gohtml",
		pathToTemplates + "footer.partial.gohtml",
	}

	templates := []string{fmt.Sprintf("%s%s", pathToTemplates, templateToShow)}
	templates = append(templates, partials...)

	tmpl, err := template.ParseFS(templateFS, templates...)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if err := tmpl.Execute(w, nil); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//TODO add logic to detect shutdown...
}
