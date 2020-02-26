package group

import (
	. "github.com/alzaburetz/workflow-api/api/server/handlers/user"
	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	"github.com/alzaburetz/workflow-api/api/server/middleware"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func GetUserGroups(w http.ResponseWriter, r *http.Request) {
	_, email := middleware.CheckToken(r)

	var user User

	var database = AccessDataStore()
	defer database.Close()

	database.DB(DBNAME).C("Users").Find(bson.M{"email":email}).One(&user)

	var groups []Group
	database.DB(DBNAME).C("Groups").Find(bson.M{"_id_":bson.M{"$in":user.Groups}}).All(&groups)

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, groups, []string{}, 200)
}

