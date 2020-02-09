package main

import ("fmt"
	   "github.com/gorilla/mux"
	   "net/http")

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello fileserver")
	})

	http.ListenAndServe(":8080", r)
}