package main

import (
	"net/http"
	"os"
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
