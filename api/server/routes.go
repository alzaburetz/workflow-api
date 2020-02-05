package server

import ("github.com/gorilla/mux"
		. "app/server/handlers/user"
		. "app/server/handlers/group"
		. "app/server/middleware"
		"net/http")


func CreateRouter() *mux.Router{
	var r = mux.NewRouter()
	r.Use(commonMiddleware)
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	var api = r.PathPrefix("/api").Subrouter()
	api.Use(AuthMiddleware)

	var user = api.PathPrefix("/user").Subrouter()
	user.HandleFunc("", GetUser).Methods("GET")
	user.HandleFunc("/register", RegisterUser).Methods("POST")
	user.HandleFunc("/login", Login).Methods("POST")
	user.HandleFunc("/update", UpdateUser).Methods("PUT")
	user.HandleFunc("/find", FindUsers).Methods("GET")

	var group = api.PathPrefix("/group").Subrouter()
	group.HandleFunc("/create", CreateGroup).Methods("POST")

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