package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	IP_Addr = "127.0.0.1"
	IP_Port = 7000
)

func check(err error) {
	if err != nil {
		log.Panicln(err.Error())
		fmt.Println(err.Error())
		panic(1)
	}
}

func main() {
	rtr := mux.NewRouter()
	rtr.HandleFunc("/", index)
	rtr.HandleFunc("/{req}", req)

	http.Handle("/", rtr)
	addr := fmt.Sprintf("%s:%d", IP_Addr, IP_Port)
	log.Printf("Serv at addr: %s", addr)
	http.ListenAndServe(addr, nil)
}

// Process index page
func index(w http.ResponseWriter, r *http.Request) {
	log.Printf("idx %s", r.URL)
	// -- show start page
	t, err := template.ParseFiles("template/index.html")
	check(err)
	t.ExecuteTemplate(w, "index", nil)
}

// Process request
func req(w http.ResponseWriter, r *http.Request) {
	log.Printf("req %s", r.URL)
	// -- Request to Redis by URL and get long URL
	// -- Redirect to url
}
