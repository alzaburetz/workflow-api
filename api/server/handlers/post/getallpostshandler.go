package post

import (
	"net/http"
	"strconv"

	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	"github.com/alzaburetz/workflow-api/api/server/middleware"
	"github.com/alzaburetz/workflow-api/api/util"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
)

func GetAllPosts(w http.ResponseWriter, r *http.Request) {
	var urlvars = mux.Vars(r)
	group, _ := urlvars["id"]

	var database = AccessDataStore()
	defer database.Close()

	var posts []Post
	// err := database.DB(DBNAME).C("Posts").

	// if err != nil {
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	WriteAnswer(&w, nil, []string{"Error getting data from database", err.Error()}, 500)
	// 	return
	// }
	query, ok := r.URL.Query()["timestamp"]
	var stamp bson.M
	if ok {
		st, _ := strconv.Atoi(query[0])
		stamp = bson.M{"$lt": st}
	} else {
		stamp = bson.M{"$gt": 0}
	}

	tags, ok := r.URL.Query()["tags"]
	var tag bson.M
	if ok {
		tag = bson.M{"tags": bson.M{"$in": tags}}
	} else {
		tag = bson.M{}
	}

	database.DB(DBNAME).C("Posts").Pipe([]bson.M{{"$match": bson.M{"$and": []bson.M{{"group_id": group}, {"timestamp": stamp}, tag}}}, {"$sort": bson.M{"timestamp": -1}}, {"$limit": 10}}).All(&posts)

	_, user := middleware.CheckToken(r)

	for key, val := range posts {
		posts[key].LikesCount = len(val.Likes)
		posts[key].LikedByUser = util.Contains(val.Likes, user)
	}

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, posts, []string{}, 200)
}
