package user

import ("net/http"
		"io/ioutil"
		. "app/server/middleware"
		. "app/server/handlers"
		"gopkg.in/mgo.v2/bson"
		"encoding/json")

func UpdateUser(w http.ResponseWriter, r *http.Request) {
	err, userKey := CheckToken(r)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		WriteAnswer(&w, "", []string{"Error getting token, try relogin", err.Error()}, 403)
		return
	}
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, "", []string{"Error reading request body", err.Error()}, 500)
		return
	}

	var userUpdated = User{}
	if err = json.Unmarshal(data, &userUpdated); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, "", []string{"Error unmarshaling json data", "Expecting json of type User", err.Error()}, 500)
		return
	}

	if err = userUpdated.HasRequiredFields(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, "", []string{err.Error()}, 400)
		return
	}
	
	var db = AccessDataStore()
	defer db.Close()

	var current User
	_ = db.DB("app").C("Users").Find(bson.M{"email":userKey}).One(&current)
	
	userUpdated.Id = current.Id
	userUpdated.UserCreated = current.UserCreated

	

	err = db.DB("app").C("Users").Update(current, userUpdated)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, "", []string{"Error updating data in database", err.Error()}, 500)
		return
	}

	if err = UpdateToken(userUpdated.Email, r.Header.Get("Token")); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, "", []string{"Redis error", err.Error()}, 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, userUpdated, []string{}, 200)
	

}