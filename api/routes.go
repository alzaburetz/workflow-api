package main

import ("github.com/gorilla/mux"
		. "app/handlers"
		"net/http")


func CreateRouter() *mux.Router{
	var r = mux.NewRouter()
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))
	r.HandleFunc("/uploadfile", UploadFile)

	var api = r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/user", GetUser).Methods("GET")
	api.HandleFunc("/user/register", RegisterUser).Methods("POST")
	api.HandleFunc("/user/login", Login).Methods("POST")
	api.HandleFunc("/user/update", UpdateUser).Methods("POST")

	var admin = r.PathPrefix("/admin").Subrouter()
	admin.HandleFunc("/wipe/{name}", DropDB)
	admin.HandleFunc("/get", GetAll)
	return r
}