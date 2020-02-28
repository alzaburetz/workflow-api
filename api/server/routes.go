package server

import (
	"encoding/json"
	. "github.com/alzaburetz/workflow-api/api/server/handlers/comment"
	. "github.com/alzaburetz/workflow-api/api/server/handlers/group"
	. "github.com/alzaburetz/workflow-api/api/server/handlers/notification"
	. "github.com/alzaburetz/workflow-api/api/server/handlers/post"
	. "github.com/alzaburetz/workflow-api/api/server/handlers/user"
	"github.com/alzaburetz/workflow-api/api/server/handlers/user/filehandlers"
	. "github.com/alzaburetz/workflow-api/api/server/handlers/user/password"
	. "github.com/alzaburetz/workflow-api/api/server/middleware"
	"github.com/gorilla/mux"
	"net/http"
)

var r *mux.Router

type Routes struct {
	Paths []string `json:"avaliable_routes"`
}

func CreateRouter() *mux.Router {
	r = mux.NewRouter()
	r.Use(commonMiddleware)
	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
	user.HandleFunc("/avatar", filehandlers.UploadAvatar).Methods("POST")
	user.HandleFunc("/login", Login).Methods("POST")
	user.HandleFunc("/update", UpdateUser).Methods("PUT")
	user.HandleFunc("/find", FindUsers).Methods("GET")
	user.HandleFunc("/notifications", GetNotifications).Methods("GET")
	user.HandleFunc("/notifications/update", UpdateNotifications).Methods("PUT", "POST")
	user.HandleFunc("/password/reset", ResetPassword).Methods("GET")
	user.HandleFunc("/password/checkcode", CheckCode).Methods("GET")
	user.HandleFunc("/password/set", SetPassword).Methods("POST", "PUT")

	var group = api.PathPrefix("/groups").Subrouter()
	group.HandleFunc("", GetAllGroups).Methods("GET")
	group.HandleFunc("/get", GetUserGroups).Methods("GET")
	group.HandleFunc("/create", CreateGroup).Methods("POST")
	group.HandleFunc("/{id}", GetGroup).Methods("GET")
	group.HandleFunc("/{id}/members", GetMembers).Methods("GET")
	group.HandleFunc("/{id}/update", UpdateGroup).Methods("PUT")
	group.HandleFunc("/{id}/enter", EnterGroup).Methods("POST")
	group.HandleFunc("/{id}/exit", ExitGroup).Methods("POST")
	group.HandleFunc("/{id}/delete", DeleteGroup).Methods("DELETE")

	var posts = group.PathPrefix("/{id}/posts").Subrouter()
	posts.HandleFunc("", GetAllPosts).Methods("GET")
	posts.HandleFunc("/add", AddPost).Methods("POST")
	posts.HandleFunc("/{post}", GetPost).Methods("GET")
	posts.HandleFunc("/{post}/delete", DeletePost).Methods("DELETE")
	posts.HandleFunc("/{post}/like", LikePost).Methods("PUT")

	var comments = posts.PathPrefix("/{post}/comments").Subrouter()
	comments.HandleFunc("", GetAllComments).Methods("GET")
	comments.HandleFunc("/add", CreateComment).Methods("POST")
	comments.HandleFunc("/{comment}/edit", EditComment).Methods("PUT")
	comments.HandleFunc("/{comment}/delete", DeleteComment).Methods("DELETE")

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
