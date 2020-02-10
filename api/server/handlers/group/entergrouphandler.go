package group

import (. "github.com/alzaburetz/workflow-api/api/server/handlers"
		. "github.com/alzaburetz/workflow-api/api/server/handlers/user"
		. "github.com/alzaburetz/workflow-api/api/server/middleware"
		. "github.com/alzaburetz/workflow-api/api/util"
		"gopkg.in/mgo.v2/bson"
		"github.com/gorilla/mux"
		"net/http")

func EnterGroup(w http.ResponseWriter, r *http.Request) {
	var muxvars = mux.Vars(r)
	idvar, converted := muxvars["id"]
	if !converted {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"Error getting variable from url", "Usage: /group/{id}/enter"}, 400)
		return
	}

	var database = AccessDataStore()
	defer database.Close()

	err, user := CheckToken(r)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"Error getting token", err.Error()}, 400)
		return
	}

	var usr User
	err = database.DB("app").C("Users").Find(bson.M{"email":user}).One(&usr)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error getting user", err.Error()},500)
		return
	}

	if Contains(usr.Groups, idvar) {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil,[]string{"User already in group"}, 400)
		return
	}

	err = database.DB("app").C("Users").Update(bson.M{"email":user}, bson.M{"$addToSet": bson.M{"groups":idvar}})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error updating user", err.Error()},500)
		return
	}

	err = database.DB("app").C("Groups").Update(bson.M{"_id_": idvar}, bson.M{"$inc":bson.M{"usercount":1}})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error updating group", err.Error()},500)
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, "Successfully entered group", []string{},200)

}