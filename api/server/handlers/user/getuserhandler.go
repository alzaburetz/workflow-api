package user

import (
	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	. "github.com/alzaburetz/workflow-api/api/server/middleware"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

//Gets user by token
func GetUser(w http.ResponseWriter, r *http.Request) {
	query, ok := r.URL.Query()["search"]
	_, userKey := CheckToken(r)
		var user User
		var database = AccessDataStore()
		var find bson.M
		if ok {
			find = bson.M{"$or":[]bson.M{bson.M{"email":string(query[0])}, bson.M{"phone":string(query[0])}}}
		} else {
			find = bson.M{"email": userKey}
		}
		database.DB(DBNAME).C("Users").Find(find).One(&user)
		w.WriteHeader(http.StatusOK)
		WriteAnswer(&w, user, []string{}, 200)
}
