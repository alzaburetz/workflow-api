package notification

import (
	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	. "github.com/alzaburetz/workflow-api/api/server/middleware"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func GetNotifications(w http.ResponseWriter, r *http.Request) {
	_, user := CheckToken(r)

	var database = AccessDataStore()
	defer database.Close()

	var userfilter = bson.M{"useremail": user}
	var seen = bson.M{"seen": false}

	var notifications []Notification
	database.DB(DBNAME).C("Notifications").Find(bson.M{"$and": []bson.M{userfilter, seen}}).All(&notifications)

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, notifications, []string{}, 200)
}