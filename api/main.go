package main

import ("log"
		//handl "github.com/gorilla/handlers"
		. "./server"
		"net/http")


func main() {
	var s Server
	s.ServerNew()
	log.Fatal(http.ListenAndServe(":3000", s.HTTP.Handler))
	//log.Fatal(http.ListenAndServeTLS(":3000", "https-server.crt" , "https-server.key", handl.CORS()(s.httpServer.Handler)))
}



