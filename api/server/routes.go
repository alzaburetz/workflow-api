package server

import ("github.com/gorilla/mux"
		. "app/server/handlers/user"
		. "app/server/handlers/group"
		. "app/server/handlers/post"
		. "app/server/middleware"
		"encoding/json"
		"net/http")

var r *mux.Router

type Routes struct {
	Paths []string `json:"avaliable_routes"`
}

func CreateRouter() *mux.Router{
	r = mux.NewRouter()
	r.Use(commonMiddleware)
	r.HandleFunc("/", func (w http.ResponseWriter, r *http.Request)  {
		var routes Routes
		routes.Paths = GetRoutes()
		json.NewEncoder(w).Encode(routes)
	})
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static/"))))

	var api = r.PathPrefix("/api").Subrouter()
	api.Use(AuthMiddleware)

	var user = api.PathPrefix("/user").Subrouter()
	user.HandleFunc("", GetUser).Methods("GET")
	user.HandleFunc("/register", RegisterUser).Methods("POST")
	user.HandleFunc("/login", Login).Methods("POST")
	user.HandleFunc("/update", UpdateUser).Methods("PUT")
	user.HandleFunc("/find", FindUsers).Methods("GET")

	var group = api.PathPrefix("/groups").Subrouter()
	group.HandleFunc("", GetAllGroups).Methods("GET")
	group.HandleFunc("/create", CreateGroup).Methods("POST")
	group.HandleFunc("/{id}", GetGroup).Methods("GET")
	group.HandleFunc("/{id}/update", UpdateGroup).Methods("PUT")
	group.HandleFunc("/{id}/enter", EnterGroup).Methods("POST")
	group.HandleFunc("/{id}/exit", ExitGroup).Methods("POST")
	group.HandleFunc("/{id}/delete", DeleteGroup).Methods("DELETE")

	var posts = group.PathPrefix("/{id}/posts").Subrouter()
	posts.HandleFunc("", GetAllPosts).Methods("GET")
	posts.HandleFunc("/add", AddPost).Methods("POST")
	posts.HandleFunc("/{post}/delete", DeletePost).Methods("DELETE")

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

func GetRoutes() []string {
	var routes []string
	r.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
        t, err := route.GetPathTemplate()
        if err != nil {
            return err
        }
		routes = append(routes, t)
		return nil
	})
	return routes
}