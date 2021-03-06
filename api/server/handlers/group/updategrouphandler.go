package group

import (
	"encoding/json"
	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
)

func UpdateGroup(w http.ResponseWriter, r *http.Request) {
	var muxvars = mux.Vars(r)
	var err error
	var id int
	groupvar, _ := muxvars["id"]

	var body []byte
	if body, err = ioutil.ReadAll(r.Body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error getting data from request body", err.Error()}, 500)
		return
	}

	var group Group
	if err = json.Unmarshal(body, &group); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error unmarshaling json", err.Error()}, 500)
		return
	}

	if err = group.HasRequiredFields(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{err.Error()}, 400)
		return
	}

	var database = AccessDataStore()
	defer database.Close()

	err = database.DB(DBNAME).C("Groups").Update(bson.M{"_id_": id}, bson.M{"$set": bson.D{{"name", group.Name}, {"description", group.Description}}})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error updating document", err.Error()}, 500)
		return
	}

	var updated Group
	database.DB(DBNAME).C("Groups").Find(bson.M{"_id_": groupvar}).One(&updated)

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, updated, []string{}, 200)

}
