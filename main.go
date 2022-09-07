package main

import "net/http"

func main() {
	http.HandleFunc("/", index)
	http.ListenAndServe(":7000", nil)
}

func index(w http.ResponseWriter, r *http.Request) {

}
