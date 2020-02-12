package main

import ("log"
		"os"
		//handl "github.com/gorilla/handlers"
		. "github.com/alzaburetz/workflow-api/api/server"
		"net/http")


func main() {
	port := os.Getenv("PORT")
	var s Server
	s.ServerNew()
	log.Fatal(http.ListenAndServe(":"+port, s.HTTP.Handler))
	//log.Fatal(http.ListenAndServeTLS(":3000", "https-server.crt" , "https-server.key", handl.CORS()(s.httpServer.Handler)))
}



