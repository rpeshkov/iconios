package main

import (
	"html/template"
	"net/http"
	"os"

	"log"
)

func openHandler(t *template.Template, ps Pagestore) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/opn/"):]
		pd, err := ps.GetPage(id)

		if err != nil {
			log.Printf("unable to load page: %v", err)
			http.Error(w, "Unknown page", http.StatusNotFound)
			return
		}

		if pd.Finished {
			http.Redirect(w, r, pd.URL, http.StatusMovedPermanently)
			return
		}

		if err := t.ExecuteTemplate(w, "page", pd); err != nil {
			log.Printf("unable to render page: %v", err)
			http.Error(w, "Unable to render page", http.StatusInternalServerError)
			return
		}
	}

	return http.HandlerFunc(handler)
}

func finishHandler(t *template.Template, ps Pagestore) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
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

		if err := os.Remove(data.IconPath()); err != nil {
			log.Println(err)
		}

		if err := ps.FinishPage(&data); err != nil {
			log.Printf("finished status update failed: %v", err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/opn/"+data.ID.String(), http.StatusMovedPermanently)
	}

	return http.HandlerFunc(handler)
}

func indexHandler(t *template.Template, ps Pagestore) http.Handler {
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

			pd, err := ps.NewPage(r.FormValue("url"), false)
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

			if err := saveFile(file, pd.IconPath()); err != nil {
				log.Printf("unable to save icon: %v", err)
				http.Error(w, "Error processing your request", http.StatusInternalServerError)
				return
			}

			http.Redirect(w, r, "/opn/"+pd.ID.String(), http.StatusMovedPermanently)
		}
	}

	return http.HandlerFunc(handler)
}
