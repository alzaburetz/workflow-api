package post

import ("net/http"
		"strconv"
		"gopkg.in/mgo.v2/bson"
		. "app/server/handlers"
		"github.com/gorilla/mux")

func DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	group, _ := vars["id"]
	post, _ := vars["post"]
	groupid, err := strconv.Atoi(group)
	postid, err :=strconv.Atoi(post)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"Error reading id", "Usage: /group/{id}/posts/{post}/delete", err.Error()}, 400)
		return
	}

	var database = AccessDataStore()
	defer database.Close()

	var post_ Post

	post_id_group := []bson.M{bson.M{"group_id":groupid},bson.M{"_id_":postid}}

	var and_operator = bson.M{"$and":post_id_group}

	var match = bson.M{"$match":and_operator}

	if err = database.DB("app").C("Posts").Find(match).One(&post_); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error getting post from database", err.Error()}, 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, post_, []string{}, 200)
}