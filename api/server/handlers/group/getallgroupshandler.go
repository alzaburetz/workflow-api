package group

import (
	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	"net/http"
)

func GetAllGroups(w http.ResponseWriter, r *http.Request) {
	var groups []Group

	var database = AccessDataStore()
	defer database.Close()

	if err := database.DB(DBNAME).C("Groups").Find(nil).All(&groups); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error fetching data from database", err.Error()}, 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, groups, []string{}, 200)
}
