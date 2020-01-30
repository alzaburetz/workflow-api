package user

import ("net/http"
		"io/ioutil"
		. "app/server/middleware"
		. "app/server/handlers"
		"gopkg.in/mgo.v2/bson"
		"golang.org/x/crypto/bcrypt"
		"encoding/json")

//Handles login
func Login(w http.ResponseWriter, r *http.Request) {
	var auth UserAuth
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, "", []string {"Couldn't read data from post form", "Expecing json"}, 400)
		return
	}

	if err = json.Unmarshal(data, &auth); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, "", []string{"Couldn't unmarshal data", err.Error(), "Expecting 'email', 'password' fields"}, 400)
		return
	}

	var userExists UserAuth
	var db = AccessDataStore()
	defer db.Close()
	db.DB("app").C("Credentials").Find(bson.M{"email": auth.Email}).One(&userExists)
	if userExists.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, "", []string{"Couldn't login", "User doesn't exist"}, 400)
		return
	}

	

	err = bcrypt.CompareHashAndPassword([]byte(userExists.Password),[]byte(auth.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		WriteAnswer(&w, "", []string{"Password is incorrect"}, 401)
		return
	} else {
		token, err := CreateToken(auth.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			WriteAnswer(&w, "", []string{"Error creating token"}, 500)
			return
		} else {
			w.WriteHeader(http.StatusOK)
			WriteAnswer(&w, token, []string{}, 200)
		}

	}

}