package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"

	"path"

	"github.com/satori/go.uuid"
)

// PageData contains main information that needed to render page
type PageData struct {
	ID  string
	URL string
}

// PageDataFromRequest creates new pageData from values provided in request
func PageDataFromRequest(r *http.Request) PageData {
	return PageData{
		ID:  r.PostFormValue("id"),
		URL: r.PostFormValue("url"),
	}
}

const (
	_        = iota             // ignore first value by assigning to blank identifier
	KB int64 = 1 << (10 * iota) // 1 << (10*1)
	MB                          // 1 << (10*2)
	GB                          // 1 << (10*3)
	TB                          // 1 << (10*4)
	PB                          // 1 << (10*5)
	EB                          // 1 << (10*6)
)

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

func finishHandler() http.Handler {
	hander := func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			r.ParseMultipartForm(256)
			data := PageDataFromRequest(r)
			wrkDir := path.Join("storage", data.ID)
			os.Remove(path.Join(wrkDir, "icon.png"))
			templateToFile("./tmpl/index_conf.html", path.Join(wrkDir, "index.html"), data)
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
			err := r.ParseMultipartForm(2 * MB)
			if err != nil {
				return
			}

			data := PageData{ID: uuid.NewV4().String(), URL: r.FormValue("url")}
			wrkDir := path.Join("storage", data.ID)

			err = os.MkdirAll(wrkDir, 0744)
			iconPath := path.Join(wrkDir, "icon.png")
			htmlPath := path.Join(wrkDir, "index.html")

			file, _, err := r.FormFile("uploadfile")
			if err != nil {
				log.Println(err)
				return
			}
			defer file.Close()

			err = saveFile(file, iconPath)
			err = templateToFile("./tmpl/index.html", htmlPath, data)

			http.Redirect(w, r, "/opn/"+data.ID, 301)
		}
	}

	return http.HandlerFunc(handler)
}

func saveFile(src io.ReadCloser, dest string) error {
	f, err := os.OpenFile(dest, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, src)
	return err
}

func templateToFile(templateFilename string, filename string, data interface{}) error {
	hfile, err := os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer hfile.Close()

	t, err := template.ParseFiles(templateFilename)
	if err != nil {
		return err
	}

	return t.Execute(hfile, data)
}
