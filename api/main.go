package main

import ("log"
		"app/handlers"
		"net/http")

var s server

func main() {
	s.Server()
	handlers.InitDatabase(s.db)
	log.Fatal(http.ListenAndServe(":3000", s.httpServer.Handler))
	//log.Fatal(http.ListenAndServeTLS(":3000", "https-server.crt" , "https-server.key", s.httpServer.Handler))
}