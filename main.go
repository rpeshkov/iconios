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

	ps, err := NewRedisPagestore("tcp", "localhost:6379")

	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize pagestore: %v", err)
		os.Exit(2)
	}

	staticHandler := http.FileServer(http.Dir("static"))
	storageHandler := http.FileServer(http.Dir("storage"))

	http.Handle("/static/", http.StripPrefix("/static/", staticHandler))
	http.Handle("/storage/", http.StripPrefix("/storage/", storageHandler))
	http.Handle("/finish/", finishHandler(t, ps))
	http.Handle("/opn/", openHandler(t, ps))
	http.Handle("/", indexHandler(t, ps))

	http.ListenAndServe(":"+port, nil)
}

func createTemplates() (*template.Template, error) {
	t, err := template.ParseGlob(filepath.Join("tmpl", "*.gohtml"))

	if err != nil {
		return nil, fmt.Errorf("unable to load index template: %v", err)
	}

	return t, nil
}
