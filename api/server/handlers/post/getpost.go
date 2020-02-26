package post

import (
	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func GetPost(w http.ResponseWriter, r *http.Request) {
	var postid = mux.Vars(r)["post"]

	var database = AccessDataStore()
	defer database.Close()

	var post Post
	database.DB(DBNAME).C("Posts").Find(bson.M{"_id_": postid}).One(&post)
	database.DB(DBNAME).C("Comments").Find(bson.M{"postid":postid}).All(&post.Comments)

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, post, []string{}, 200)
}
