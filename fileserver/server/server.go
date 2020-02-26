package server

import (
	"gopkg.in/mgo.v2"
	"log"
	"net/http"
)

type Server struct {
	HTTP *http.Server
	db   *mgo.Session
}

func (s *Server) FileServerInit() {
	s.HTTP = CreateHttpServer()
	s.InitDatabase()
	log.Fatal(s.HTTP.ListenAndServe())
	log.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
}

func CreateHttpServer() *http.Server {
	return &http.Server{
		Addr:    ":8080",
		Handler: CreateRouter(),
	}
}
