package main

import (
	"html/template"
	"net/http"
	"os"

	uuid "github.com/satori/go.uuid"
)

func finishHandler() http.Handler {
	hander := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			r.ParseMultipartForm(256)
			data := PageDataFromRequest(r)
			os.Remove(data.IconPath())
			templateToFile("./tmpl/index_conf.html", data.HTMLPath(), data)
			http.Redirect(w, r, "/opn/"+data.ID, 302)
		}
	}

	return http.HandlerFunc(hander)
}

func indexHandler() http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			t, _ := template.ParseFiles("content/index.html")
			t.Execute(w, nil)
		case http.MethodPost:
			if err := r.ParseMultipartForm(2 * MB); err == nil {
				data := PageData{
					ID:  uuid.NewV4().String(),
					URL: r.FormValue("url"),
				}

				if err = data.InitWorkDir(); err == nil {
					if file, _, err := r.FormFile("uploadfile"); err == nil {
						defer file.Close()
						err = saveFile(file, data.IconPath())
						err = templateToFile("./tmpl/index.html", data.HTMLPath(), data)
						http.Redirect(w, r, "/opn/"+data.ID, 301)
					}
				}

			}

		}
	}

	return http.HandlerFunc(handler)
}
