package user

import (
	"encoding/json"
	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"time"
)

func AddUser(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"Error reading request body", err.Error()}, 400)
		return
	}

	var user User
	if err = json.Unmarshal(body, &user); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"Error unmarshaling data", err.Error()}, 400)
		return
	}

	if err = user.HasRequiredFields(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"Error reading user object", "Missing data",err.Error()},400)
		return
	}

	var database = AccessDataStore()
	defer database.Close()

	var checkuser User
	database.DB(DBNAME).C("Users").Find(bson.M{"$or":[]bson.M{{"email":user.Email}, {"phone":user.Phone}}}).One(&checkuser)

	if checkuser.Email == user.Email || checkuser.Phone == user.Phone{
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"User already exists"}, 400)
		return
	}

	user.Id = uuid.NewV4().String()
	user.UserCreated = time.Now().Unix()

	if err = database.DB(DBNAME).C("Users").Insert(user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error inserting data to database", err.Error()}, 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, user, []string{}, 200)
}
