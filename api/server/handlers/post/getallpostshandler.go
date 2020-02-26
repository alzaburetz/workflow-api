package post

import (
	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	"github.com/alzaburetz/workflow-api/api/server/middleware"
	"github.com/alzaburetz/workflow-api/api/util"
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

	if err := database.DB(DBNAME).C("Posts").Find(bson.M{"group_id": group}).All(&posts); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error getting data from database", err.Error()}, 500)
		return
	}

	_, user := middleware.CheckToken(r)

	for key, val := range posts {
		posts[key].LikesCount = len(val.Likes)
		posts[key].LikedByUser = util.Contains(val.Likes, user)
	}

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, posts, []string{}, 200)
}
