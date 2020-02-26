package group

import (
	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func GetAllGroups(w http.ResponseWriter, r *http.Request) {
	var groups []Group

	query, ok := r.URL.Query()["search"]
	var filter bson.M
	if ok {
		filter = bson.M{"name": bson.M{"$regex": ".*" + string(query[0]) + ".*"}}
	} else {
		filter = nil
	}

	var database = AccessDataStore()
	defer database.Close()

	database.DB(DBNAME).C("Groups").Find(filter).All(&groups)

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, groups, []string{}, 200)
}
