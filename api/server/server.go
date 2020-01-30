 package server

import ("time"
		"crypto/tls"
		"gopkg.in/mgo.v2"
		"net/http")

type Server struct {
	HTTP *http.Server
	db *mgo.Session
}

func (s *Server) ServerNew() {
	s.HTTP= CreateHTTPServer()
	db, err := mgo.Dial("mongo:27017")
	if err != nil {
		panic("Could not resolve connection to database")
	}
	s.db = db
}

func CreateHTTPServer() *http.Server  {
	var Router = CreateRouter()
	Config, err := ConfigHTTPS()
	if err != nil {
		panic("Could not create config")
	}
	return &http.Server {
		ReadTimeout: 5*time.Second,
		WriteTimeout: 5*time.Second,
		Handler: Router,
		TLSConfig: Config,
	}
}

func ConfigHTTPS() (*tls.Config, error) {
	var err error
	return &tls.Config {
		PreferServerCipherSuites: true,
		CipherSuites: []uint16 {
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		},
	}, err
}

