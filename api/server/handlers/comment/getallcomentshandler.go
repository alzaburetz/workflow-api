package comment

import (
	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func GetAllComments(w http.ResponseWriter, r *http.Request) {
	var id_post = mux.Vars(r)
	var postid = id_post["post"]
	var comments []Comment
	database := AccessDataStore()
	defer database.Close()
	database.DB(DBNAME).C("Comments").Find(bson.M{"postid": postid}).All(&comments)
	w.WriteHeader(200)
	WriteAnswer(&w, comments, []string{}, 200)
}
