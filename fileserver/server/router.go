package server

import "github.com/gorilla/mux"
import . "app/server/handlers"

func CreateRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/upload", UploadFile)
	return r
}