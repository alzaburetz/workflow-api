package post

import (
	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	var urlvars = mux.Vars(r)
	group, _ := urlvars["id"]

	var database = AccessDataStore()
	defer database.Close()

	var posts []Post

	if err := database.DB("app").C("Posts").Find(bson.M{"group_id": group}).All(&posts); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error getting data from database", err.Error()}, 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, posts, []string{}, 200)
}
