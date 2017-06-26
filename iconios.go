package main

import (
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

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

	fs := http.FileServer(http.Dir("test"))

	http.HandleFunc("/", encodeHandler)
	http.Handle("/test/", http.StripPrefix("/test", fs))

	http.ListenAndServe(":"+port, nil)
}

func encodeHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t, _ := template.ParseFiles("content/index.html")
		t.Execute(w, nil)
	} else {
		var data struct {
			ID  string
			URL string
		}
		data.ID = "test"
		r.ParseMultipartForm(2 * MB)
		file, _, err := r.FormFile("uploadfile")
		if err != nil {
			log.Println(err)
			return
		}

		defer file.Close()

		f, err := os.OpenFile("./test/icon.png", os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println(err)
			return
		}
		defer f.Close()

		io.Copy(f, file)

		hfile, err := os.OpenFile("./test/index.html", os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println(err)
			return
		}
		defer hfile.Close()

		data.URL = r.FormValue("url")

		t, _ := template.ParseFiles("./tmpl/index.html")
		t.Execute(hfile, data)

		http.Redirect(w, r, "/", 301)
	}

}
