package server

import "github.com/gorilla/mux"
import . "github.com/alzaburetz/workflow-api/fileserver/server/handlers"

func CreateRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/upload", UploadFile)
	return r
}
