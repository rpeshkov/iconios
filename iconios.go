package main

import (
	"html/template"
	"io"
	"net/http"
	"os"
	"path"
)

// PageData contains main information that needed to render page
type PageData struct {
	ID  string
	URL string
}

// InitWorkDir returns
func (p *PageData) InitWorkDir() error {
	return os.MkdirAll(path.Join("storage", p.ID), 0744)
}

// IconPath returns
func (p *PageData) IconPath() string {
	return path.Join("storage", p.ID, "icon.png")
}

// HTMLPath returns
func (p *PageData) HTMLPath() string {
	return path.Join("storage", p.ID, "index.html")
}

// PageDataFromRequest creates new pageData from values provided in request
func PageDataFromRequest(r *http.Request) PageData {
	return PageData{
		ID:  r.PostFormValue("id"),
		URL: r.PostFormValue("url"),
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fs := http.FileServer(http.Dir("storage"))

	http.Handle("/", indexHandler())
	http.Handle("/finish/", finishHandler())
	http.Handle("/opn/", http.StripPrefix("/opn", fs))

	http.ListenAndServe(":"+port, nil)

}

func saveFile(src io.ReadCloser, dest string) (err error) {
	if f, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0666); err == nil {
		defer f.Close()
		_, err = io.Copy(f, src)
	}
	return
}

func templateToFile(templateFilename string, filename string, data interface{}) (err error) {
	if f, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666); err == nil {
		defer f.Close()

		if t, err := template.ParseFiles(templateFilename); err == nil {
			return t.Execute(f, data)
		}
	}
	return
}
