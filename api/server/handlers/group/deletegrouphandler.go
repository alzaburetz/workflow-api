package group

import ("net/http"
		. "app/server/handlers"
		// . "app/server/handlers/user"
		. "app/server/middleware"
		"gopkg.in/mgo.v2/bson"
		"github.com/gorilla/mux")

func DeleteGroup(w http.ResponseWriter, r *http.Request) {
	var urlvars = mux.Vars(r)
	groupid, parsed := urlvars["id"]
	if !parsed {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"Error getting id from url", "Usage: /group/{id}/delete"}, 400)
		return
	}

	var database = AccessDataStore()
	defer database.Close()

	err, email := CheckToken(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error getting user token", err.Error()}, 500)
		return
	}

	var group Group

	iter := database.DB("app").C("Groups").Pipe([]bson.M{{"$match": bson.M{"$and": []bson.M{{"creator.email":email}, {"_id_":groupid}}}}}).Iter()
	var groupcount int 
	for iter.Next(&group) {
		groupcount++
	}
	if groupcount == 0 {
		w.WriteHeader(http.StatusUnauthorized)
		WriteAnswer(&w,nil,[]string{"You are not the owner"},401)
		return
	}
	_, err = database.DB("app").C("Users").UpdateAll(bson.M{"_id_":bson.M{"$gte":0}}, bson.M{"$pull":bson.M{"groups":groupid}})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{err.Error()}, 500)
		return
	}

	if err = database.DB("app").C("Groups").Remove(bson.M{"_id_":groupid}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error removing group", err.Error()}, 500)
		return
	}
	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, "Successfully deleted group", []string{}, 200)
	
}