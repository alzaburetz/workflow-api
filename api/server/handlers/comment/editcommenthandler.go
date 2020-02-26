package comment

import (
	"encoding/json"
	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"time"
)

func EditComment(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"Error geting request body", err.Error()}, 400)
		return
	}

	var comment Comment
	if err = json.Unmarshal(body, &comment); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"Error unmarshaling data", err.Error()}, 400)
		return
	}

	var urlvars = mux.Vars(r)
	var commentid = urlvars["comment"]

	var database = AccessDataStore()
	defer database.Close()

	if err = database.DB(DBNAME).C("Comments").Update(bson.M{"_id_": commentid}, bson.M{"$set": bson.M{"edited": time.Now().Unix(), "body": comment.Body}}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error updateing data in database", err.Error()}, 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, "Successfully updated comment", []string{}, 200)

}
