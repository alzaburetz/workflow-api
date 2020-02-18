package group

import (
	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func GetGroup(w http.ResponseWriter, r *http.Request) {
	var muxvar = mux.Vars(r)
	var idurl string
	var err error
	var converted bool
	if idurl, converted = muxvar["id"]; converted != true {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"ID of group not provided", "Usage: /group/{id}"}, 400)
		return
	}
	var group Group
	var database = AccessDataStore()
	defer database.Close()
	if err := database.DB(DBNAME).C("Groups").Find(bson.M{"_id_": idurl}).One(&group); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Database error", "Error getting group from database", err.Error()}, 500)
		return
	}
	if err = database.DB(DBNAME).C("Posts").Pipe([]bson.M{bson.M{"$match": bson.M{"group_id": idurl}}, bson.M{"$limit": 2}}).All(&group.Posts); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Database error", "Error getting posts from database", err.Error()}, 500)
		return
	}
	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, group, []string{}, 200)
}
