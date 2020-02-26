package group

import (
	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	. "github.com/alzaburetz/workflow-api/api/server/handlers/user"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func GetMembers(w http.ResponseWriter, r *http.Request) {
	var groupid = mux.Vars(r)["id"]

	var database = AccessDataStore()
	defer database.Close()

	var group Group
	if err := database.DB(DBNAME).C("Groups").Find(bson.M{"_id_":groupid}).One(&group); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error fetching data from database", err.Error()},500)
		return
	}

	var users []User
	if err := database.DB(DBNAME).C("Users").Find(bson.M{"email": bson.M{"$in":group.Members}}).All(&users); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error getting users", err.Error()}, 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, users, []string{}, 200)
}
