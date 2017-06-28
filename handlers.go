package main

import (
	"html/template"
	"net/http"
	"os"

	"log"

	uuid "github.com/satori/go.uuid"
)

func finishHandler() http.Handler {
	hander := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			if err := r.ParseMultipartForm(256); err != nil {
				log.Printf("can't parse multipart form: %v", err)
				http.Error(w, "Bad data provided", http.StatusBadRequest)
				return
			}

			data := PageDataFromRequest(r)

			if _, err := os.Stat(data.HTMLPath()); os.IsNotExist(err) {
				log.Printf("html file doesn't exist: %v", err)
				http.Error(w, "Wrong id provided", http.StatusBadRequest)
				return
			}

			if err := os.Remove(data.IconPath()); err != nil {
				log.Println(err)
			}

			if err := templateToFile("./tmpl/index_conf.html", data.HTMLPath(), data); err != nil {
				log.Printf("unable to save confirmation template to file %s: %v", data.HTMLPath(), err)
				http.Error(w, "Error processing your request", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/opn/"+data.ID, http.StatusMovedPermanently)
		}
	}

	return http.HandlerFunc(hander)
}

func indexHandler() http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			t, err := template.ParseFiles("content/index.html")
			if err != nil {
				log.Printf("unable to parse index.html template: %v", err)
				http.Error(w, "Index page not found", http.StatusNotFound)
				return
			}
			t.Execute(w, nil)
		case http.MethodPost:
			err := r.ParseMultipartForm(2 * MB)
			if err != nil {
				log.Printf("unable to parse multipart form: %v", err)
				http.Error(w, "Bad data provided", http.StatusBadRequest)
				return
			}

			data := PageData{
				ID:  uuid.NewV4().String(),
				URL: r.FormValue("url"),
			}

			err = data.InitWorkDir()
			if err != nil {
				log.Printf("unable to init workdir: %v", err)
				http.Error(w, "Error processing your request", http.StatusInternalServerError)
				return
			}

			file, _, err := r.FormFile("uploadfile")
			if err != nil {
				log.Printf("unable to read file from POST: %v", err)
				http.Error(w, "Error processing your request", http.StatusInternalServerError)
				return
			}
			defer file.Close()

			err = saveFile(file, data.IconPath())
			if err != nil {
				log.Printf("unable to save icon: %v", err)
				http.Error(w, "Error processing your request", http.StatusInternalServerError)
				return
			}

			err = templateToFile("./tmpl/index.html", data.HTMLPath(), data)
			if err != nil {
				log.Printf("unable to save confirmation template to file %s: %v", data.HTMLPath(), err)
				http.Error(w, "Error processing your request", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/opn/"+data.ID, http.StatusMovedPermanently)
		}
	}

	return http.HandlerFunc(handler)
}
