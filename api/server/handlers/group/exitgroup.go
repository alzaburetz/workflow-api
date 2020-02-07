package group

import ("net/http"
		"strconv"
		. "app/server/handlers"
		. "app/server/middleware"
		. "app/server/handlers/user"
		. "app/util"
		"gopkg.in/mgo.v2/bson"
		"github.com/gorilla/mux")

func ExitGroup(w http.ResponseWriter, r *http.Request) {
	var muxvar = mux.Vars(r)
	muxid, converted := muxvar["id"]
	if !converted {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"Error reading id", "Usage: /group/{id}/exit"}, 400)
		return
	}

	id, err := strconv.Atoi(muxid)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error converting to int", err.Error()}, 500)
		return
	}

	var user User
	err, email := CheckToken(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error getting token", err.Error()}, 500)
		return
	} 

	var database = AccessDataStore()
	defer database.Close()

	err = database.DB("app").C("Users").Find(bson.M{"email":email}).One(&user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error fetching user from database", err.Error()}, 500)
		return
	}

	if !Contains(user.Groups, id) {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"User already exited group!"}, 400)
		return
	}

	err = database.DB("app").C("Users").Update(bson.M{"email":email}, bson.M{"$pull": bson.M{"groups":id}})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error updating user", err.Error()}, 500)
		return
	}

	err = database.DB("app").C("Groups").Update(bson.M{"_id_":id}, bson.M{"$inc": bson.M{"usercount": -1}})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error updating group", err.Error()}, 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, "Successfuly exited group", []string{}, 200)
}