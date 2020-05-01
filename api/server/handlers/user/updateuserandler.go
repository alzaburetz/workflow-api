package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	. "github.com/alzaburetz/workflow-api/api/server/middleware"
	"github.com/alzaburetz/workflow-api/api/util"
	"gopkg.in/mgo.v2/bson"
)

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
	_ = db.DB(DBNAME).C("Users").Find(bson.M{"email": userKey}).One(&current)

	userUpdated.Id = current.Id
	userUpdated.UserCreated = current.UserCreated

	if len(current.Schedule) == 0 {
		userUpdated.Schedule = util.CalculateCalendar(time.Unix(userUpdated.FirstWorkDay, 0), userUpdated.Workdays, userUpdated.Weekdays)
	}

	// WriteAnswer(&w, current, []string{},200)
	// return

	err = db.DB(DBNAME).C("Users").Update(bson.M{"email": userKey}, userUpdated)
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
