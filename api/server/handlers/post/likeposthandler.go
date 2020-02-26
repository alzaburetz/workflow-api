package post

import (
	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	. "github.com/alzaburetz/workflow-api/api/server/middleware"
	"github.com/alzaburetz/workflow-api/api/util"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func LikePost(w http.ResponseWriter, r *http.Request) {
	var postid = mux.Vars(r)["post"]
	_, user := CheckToken(r)

	var database = AccessDataStore()
	defer database.Close()

	var post Post

	var arrOp string

	database.DB(DBNAME).C("Posts").Find(bson.M{"_id_": postid}).One(&post)
	if util.Contains(post.Likes, user) {
		arrOp = "$pull"
	} else {
		arrOp = "$addToSet"
	}

	if err := database.DB(DBNAME).C("Posts").Update(bson.M{"_id_": postid}, bson.M{arrOp: bson.M{"likes": user}}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error liking post", err.Error()}, 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, "Successfully liked/disliked post", []string{}, 200)
}
