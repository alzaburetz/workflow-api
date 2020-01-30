package server

import ("github.com/gorilla/mux"
		. "app/server/handlers/user"
		_ "app/server/middleware"
		"net/http")


func CreateRouter() *mux.Router{
	var r = mux.NewRouter()
	r.Use(commonMiddleware)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	var api = r.PathPrefix("/api").Subrouter()
	api.HandleFunc("/user", GetUser).Methods("GET")
	api.HandleFunc("/user/register", RegisterUser).Methods("POST")
	api.HandleFunc("/user/login", Login).Methods("POST")
	api.HandleFunc("/user/update", UpdateUser).Methods("POST")
	api.HandleFunc("/user/find", FindUsers).Methods("GET")

	// var admin = r.PathPrefix("/admin").Subrouter()
	// admin.HandleFunc("/wipe/{name}", DropDB)
	// admin.HandleFunc("/get", GetAll)
	return r
}

func commonMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Add("Content-Type", "application/json")
        next.ServeHTTP(w, r)
    })
}