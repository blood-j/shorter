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
	rtr.HandleFunc("/short", short)
	rtr.HandleFunc("/status/{req}", status)
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
	t, err := template.ParseFiles("template/index.html", "template/footer.html", "template/header.html")
	check(err)
	t.ExecuteTemplate(w, "index", nil)
}

// Process request
func req(w http.ResponseWriter, r *http.Request) {
	log.Printf("req %s", r.URL)
	// -- Request to Redis by URL and get long URL
	// -- Redirect to url
}

func short(w http.ResponseWriter, r *http.Request) {
	// -- Get url for shorting
	// -- Generate short url
	// -- Check new url in db for uniq
	// -- Save data in db
	// -- Redirect to status page
}

func status(w http.ResponseWriter, r *http.Request) {
	// -- Get url, separate short link
	// -- Get info, from db
	// -- Show status page
}
