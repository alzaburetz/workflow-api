package group

import ("net/http"
		. "app/server/handlers"
		"gopkg.in/mgo.v2/bson"
		"strconv"
		"github.com/gorilla/mux")

func GetGroup(w http.ResponseWriter, r *http.Request) {
	var muxvar = mux.Vars(r)
	var idurl string
	var err error
	var id int
	var converted bool
	if idurl, converted = muxvar["id"]; converted != true {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"ID of group not provided", "Usage: /group/{id}"}, 400)
		return
	}
	if id, err = strconv.Atoi(idurl); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{ err.Error() }, 500)
		return
	}
	var group Group
	var database = AccessDataStore()
	defer database.Close()
	if err := database.DB("app").C("Groups").Find(bson.M{"_id_":id}).One(&group); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Database error", "Error getting group from database", err.Error()}, 500)
		return
	}
	if err = database.DB("app").C("Posts").Pipe([]bson.M{bson.M{"$match": bson.M{"group_id":id}}, bson.M{"$limit":2}}).All(&group.Posts); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Database error", "Error getting posts from database", err.Error()}, 500)
		return
	}
	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, group, []string{}, 200)
}