package main

import ("github.com/gorilla/mux"
		"app/handlers"
		"fmt"
		"net/http")

func CreateRouter() *mux.Router{
	var r = mux.NewRouter()
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world")
	})
	r.HandleFunc("/user/create", handlers.CreateUser).Methods("POST")
	r.HandleFunc("/user/list", handlers.ListUsers).Methods("GET")
	return r
}