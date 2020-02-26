package comment

import (
	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func DeleteComment(w http.ResponseWriter, r *http.Request) {
	var commentid = mux.Vars(r)["comment"]

	var database = AccessDataStore()
	defer database.Close()

	if err := database.DB(DBNAME).C("Comments").Remove(bson.M{"_id_":commentid}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error removing comment", err.Error()}, 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, "Successfully removed comment", []string{}, 200)
}
