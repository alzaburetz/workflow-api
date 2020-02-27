package password

import (
	"encoding/json"
	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"io/ioutil"
	"net/http"
	"strings"
)

type ResetCredentials struct {
	Phone string `json:"phone"`
	Code string `json:"code"`
	Password string `json:"password"`
}

func SetPassword(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{err.Error(), "Error retrieving data!", "Body should be present", "Usage: {'phone':'...','code':...,'password':'...'"}, 400)
		return
	}

	var reset ResetCredentials
	if err = json.Unmarshal(body, &reset); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error unmarshaling json", err.Error()}, 500)
		return
	}

	codefromstorage, ok := Storage[reset.Phone]
	if !ok {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error getting code from storage"}, 500)
		return
	}

	if codefromstorage != reset.Code {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"Wrong code provided!"}, 400)
		return
	}

	database := AccessDataStore()
	defer database.Close()

	newpass, err := bcrypt.GenerateFromPassword([]byte(reset.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error hashing password", err.Error()}, 500)
		return
	}

	database.DB(DBNAME).C("Credentials").Update(bson.M{"phone": "+" + strings.TrimSpace(reset.Phone)}, bson.M{"$set":bson.M{"password":string(newpass)}})
	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, "Successfully changed password", []string{}, 200)
}
