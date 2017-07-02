package main

import (
	"net/http"
	"path"

	uuid "github.com/satori/go.uuid"
)

// PageData contains main information that needed to render page
type PageData struct {
	ID       uuid.UUID
	URL      string
	Finished bool
}

// IconPath returns
func (p *PageData) IconPath() string {
	return path.Join("storage", p.ID.String()+".png")
}

// PageDataFromRequest creates new pageData from values provided in request
func PageDataFromRequest(r *http.Request) PageData {
	return PageData{
		ID:  uuid.FromStringOrNil(r.PostFormValue("id")),
		URL: r.PostFormValue("url"),
	}
}
