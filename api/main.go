package main

import ("log"
		"os"
		"app/handlers"
		handl "github.com/gorilla/handlers"
		. "runtime"
		"app/middleware"
		"net/http")

var s server

func main() {
	_ = GOMAXPROCS(10)
	s.Server()
	handlers.InitDatabase(s.db)
	if err := middleware.RedisInit(); err != nil {
		os.Exit(1)
	}
	log.Fatal(http.ListenAndServe(":3000", handl.CORS()(s.httpServer.Handler)))
	//log.Fatal(http.ListenAndServeTLS(":3000", "https-server.crt" , "https-server.key", handl.CORS()(s.httpServer.Handler)))
}



