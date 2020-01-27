package handlers

import ("net/http"
		"github.com/gorilla/mux")

func DropDB(w http.ResponseWriter, r *http.Request) {
	dbname := mux.Vars(r)
	err := database.DB(dbname["name"]).DropDatabase()
	if err != nil {
		WriteAnswer(&w, "", []string{err.Error()}, 500)
		return
	}
	WriteAnswer(&w, "Database wiped", []string{}, 200)
}

func GetAll(w http.ResponseWriter, r *http.Request) {
	var users []User
		if err := database.DB("app").C("Users").Find(nil).All(&users); err != nil{
			w.WriteHeader(http.StatusInternalServerError)
			WriteAnswer(&w, "", []string{err.Error()}, 500)
			return
		}
		WriteAnswer(&w, users, []string{}, 200)
}