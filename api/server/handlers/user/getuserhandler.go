package user

import ("net/http"
		. "app/server/middleware"
		. "app/server/handlers"
		"gopkg.in/mgo.v2/bson")

//Gets user by token
func GetUser(w http.ResponseWriter, r *http.Request) {
	if err, userKey := CheckToken(r); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		WriteAnswer(&w, err.Error(), []string{"Wrong token, try relogin"}, 403)
	} else {
		w.WriteHeader(http.StatusOK)
		var user User
		var database = AccessDataStore()
		if err = database.DB("app").C("Users").Find(bson.M{"email":userKey}).One(&user); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			WriteAnswer(&w, "", []string{err.Error()}, 500)
			return
		}
		WriteAnswer(&w, user, []string{}, 200)
	}
}