package main

import (
	"html/template"
	"net/http"
	"os"

	"log"

	uuid "github.com/satori/go.uuid"
)

func finishHandler(t *template.Template) http.Handler {
	hander := func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Header().Set("Allow", http.MethodPost)
			http.Error(w, http.StatusText(http.StatusMethodNotAllowed),
				http.StatusMethodNotAllowed)
			return
		}

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

		if err := saveTemplateToFile(data.HTMLPath(), t, "confirmed", data); err != nil {
			log.Printf("unable to save confirmation template to file %s: %v", data.HTMLPath(), err)
			http.Error(w, "Error processing your request", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/opn/"+data.ID, http.StatusMovedPermanently)
	}

	return http.HandlerFunc(hander)
}

func indexHandler(t *template.Template) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			if err := t.ExecuteTemplate(w, "app", nil); err != nil {
				log.Printf("unable to execute template 'index': %v", err)
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}
		case http.MethodPost:
			if err := r.ParseMultipartForm(2 * MB); err != nil {
				log.Printf("unable to parse multipart form: %v", err)
				http.Error(w, "Bad data provided", http.StatusBadRequest)
				return
			}

			data := PageData{
				ID:  uuid.NewV4().String(),
				URL: r.FormValue("url"),
			}

			if err := data.InitWorkDir(); err != nil {
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

			if err := saveFile(file, data.IconPath()); err != nil {
				log.Printf("unable to save icon: %v", err)
				http.Error(w, "Error processing your request", http.StatusInternalServerError)
				return
			}

			if err := saveTemplateToFile(data.HTMLPath(), t, "page", data); err != nil {
				log.Printf("unable to save confirmation template to file %s: %v", data.HTMLPath(), err)
				http.Error(w, "Error processing your request", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/opn/"+data.ID, http.StatusMovedPermanently)
		}
	}

	return http.HandlerFunc(handler)
}
