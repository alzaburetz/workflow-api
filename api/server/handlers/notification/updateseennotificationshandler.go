package notification

import (
	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	. "github.com/alzaburetz/workflow-api/api/server/middleware"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"time"
)

func UpdateNotifications(w http.ResponseWriter, r *http.Request) {
	_, user := CheckToken(r)

	var database = AccessDataStore()
	defer database.Close()

	var createdfilter = bson.M{"created": bson.M{"$lt": time.Now().Unix()}}
	var userfilter = bson.M{"useremail": user}

	var filter = bson.M{"$and": []bson.M{createdfilter, userfilter}}

	_, err := database.DB(DBNAME).C("Notifications").UpdateAll(filter, bson.M{"$set": bson.M{"seen": true}})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error updating resords", err.Error()}, 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, "Successfully updated notifications", []string{}, 200)
}
