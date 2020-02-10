package user

import ("net/http"
		"encoding/json"
		. "github.com/alzaburetz/workflow-api/api/server/handlers"
		"gopkg.in/mgo.v2/bson"
		"io/ioutil")

func FindUsers(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, "", []string{"Error reading data from request body", err.Error()}, 500)
		return
	}

	var phones []string

	if err = json.Unmarshal(data, &phones); err!= nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, "", []string{"Error unmarshalling to json", err.Error()}, 500)
		return
	}

	var db = AccessDataStore()
	defer db.Close()

	iter := db.DB("app").C("Users").Find(bson.M{"phone":bson.M{"$in":phones}}).Iter();
	var founduser User
	var foundUsers []User

	for iter.Next(&founduser) {
		foundUsers = append(foundUsers, founduser)
	}

	WriteAnswer(&w, foundUsers, []string{}, 200)
	return 

	if err = db.DB("app").C("Users").Find(bson.M{"phone":bson.M{"$in":phones}}).All(&foundUsers); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, "", []string{"Error reading data from database", err.Error()}, 500)
		return
	}



	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, foundUsers, []string{}, 200)
}