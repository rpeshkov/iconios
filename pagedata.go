package main

import (
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
