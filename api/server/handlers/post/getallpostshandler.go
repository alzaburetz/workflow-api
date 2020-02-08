package post

import ("net/http"
		"gopkg.in/mgo.v2/bson"
		"strconv"
		"github.com/gorilla/mux"
		. "app/server/handlers")

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	var urlvars = mux.Vars(r)
	group, _ := urlvars["id"]
	groupid, err := strconv.Atoi(group)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"Error convering id to int", err.Error()}, 400)
		return
	}

	var database = AccessDataStore()
	defer database.Close()

	var posts []Post

	if err = database.DB("app").C("Posts").Find(bson.M{"group_id":groupid}).All(&posts); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error getting data from database", err.Error()}, 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, posts, []string{}, 200)
}