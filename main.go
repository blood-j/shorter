package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index)
	rtr.HandleFunc("/{req}", req)

	http.Handle("/", rtr)
	http.ListenAndServe(":7000", nil)
}

// Process index page
func index(w http.ResponseWriter, r *http.Request) {
	log.Printf("idx %s", r.URL)
}

// Process request
func req(w http.ResponseWriter, r *http.Request) {
	log.Printf("req %s", r.URL)
}
