package group

import ("net/http"
		"gopkg.in/mgo.v2/bson"
		"strconv"
		"io/ioutil"
		"github.com/gorilla/mux"
		"encoding/json"
		."app/server/handlers")

func UpdateGroup(w http.ResponseWriter, r *http.Request) {
	var muxvars = mux.Vars(r)
	var err error
	var id int
	groupvar, _ := muxvars["id"]
	if id, err = strconv.Atoi(groupvar); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{err.Error()}, 400)
		return
	}

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
		WriteAnswer(&w, nil, []string{ err.Error()}, 400)
		return
	}

	var database = AccessDataStore()
	defer database.Close()

	err = database.DB("app").C("Groups").Update(bson.M{"_id_":id}, bson.M{"$set": bson.D{{"name",group.Name}, {"description",group.Description}}})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error updating document", err.Error()}, 500)
		return
	}

	var updated Group
	database.DB("app").C("Groups").Find(bson.M{"_id_":id}).One(&updated)

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, updated, []string{},200)


}