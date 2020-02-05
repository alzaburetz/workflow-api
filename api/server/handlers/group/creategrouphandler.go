package group

import ("net/http"
		"gopkg.in/mgo.v2/bson"
		"io/ioutil"
		"encoding/json"
		. "app/server/middleware"
		. "app/server/handlers"
		. "app/server/handlers/user")


func CreateGroup(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{ "Error reading request body",err.Error()}, 500)
		return
	}

	var group Group

	if err = json.Unmarshal(body, &group); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{ "Error unmarshaling data",err.Error()}, 500)
		return
	}

	if err = group.HasRequiredFields(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{err.Error()}, 400)
		return
	}

	err, user := CheckToken(r) 
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{ err.Error()}, 500)
		return
	}

	var database = AccessDataStore()
	defer database.Close()

	var creator User
	if err = database.DB("app").C("Users").Find(bson.M{"email":user}).One(&creator); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{ err.Error()}, 500)
		return
	}

	group.Creator = creator
	group.Id, _ = database.DB("app").C("Groups").Count()
	group.UserCount = 1

	if err = database.DB("app").C("Groups").Insert(group); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{ err.Error()}, 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, group, []string{ }, 200)

} 
