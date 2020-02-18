package server

import (
	"crypto/tls"
	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	"gopkg.in/mgo.v2"
	"net/http"
	"time"
)

type Server struct {
	HTTP *http.Server
	db   *mgo.Session
}

func (s *Server) ServerNew() {
	s.HTTP = CreateHTTPServer()
	CreateDatabaseInstance()
}

func CreateHTTPServer() *http.Server {
	var Router = CreateRouter()
	Config, err := ConfigHTTPS()
	if err != nil {
		panic("Could not create config")
	}
	return &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		Handler:      Router,
		TLSConfig:    Config,
	}
}

func ConfigHTTPS() (*tls.Config, error) {
	var err error
	return &tls.Config{
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}, err
}
