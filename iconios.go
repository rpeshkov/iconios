package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	t, err := createTemplates()

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load templates: %v", err)
		os.Exit(2)
	}

	fs := http.FileServer(http.Dir("storage"))
	staticHandler := http.FileServer(http.Dir("static"))

	http.Handle("/", indexHandler(t))
	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))
	http.Handle("/finish/", finishHandler(t))
	http.Handle("/opn/", http.StripPrefix("/opn", fs))

	http.ListenAndServe(":"+port, nil)
}

func createTemplates() (*template.Template, error) {
	t, err := template.ParseGlob(filepath.Join("tmpl", "*.gohtml"))

	if err != nil {
		return nil, fmt.Errorf("unable to load index template: %v", err)
	}

	return t, nil
}
