package user

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2/bson"
)

//Handle register
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var auth *UserAuth
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errs := make([]string, 1)
		errs[0] = "Could not read request body"
		message := Resp{
			Code:     500,
			Errors:   errs,
			Response: err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message)
		return
	}

	if err = json.Unmarshal(data, &auth); err != nil {
		errs := make([]string, 1)
		errs[0] = "Could not unmarshal json data"
		message := Resp{
			Code:     500,
			Errors:   errs,
			Response: err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message)
		return
	}

	if err = auth.HasRequiredFields(); err != nil {
		errs := make([]string, 1)
		errs[0] = err.Error()
		message := Resp{
			Code:     400,
			Errors:   errs,
			Response: nil,
		}
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(message)
		return
	}

	//FULL SUCCESS GO AHEAD AND CREATE USER, HOORAY

	//This part checks if user aready exists
	var userExists User
	database := AccessDataStore()
	if database == nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, "", []string{"Database is nil"}, 400)
		return
	}
	database.DB(DBNAME).C("Credentials").Find(bson.M{"$or": []bson.M{bson.M{"email": auth.Email}, bson.M{"phone": auth.Phone}}}).One(&userExists)
	if userExists.Email != "" || userExists.Phone != "" { //if user is found, return error
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, "", []string{"User already exists"}, 400)
		return
	}

	database.DB(DBNAME).C("Users").Find(bson.M{"$or": []bson.M{bson.M{"email": auth.Email}, bson.M{"phone": auth.Phone}}}).One(&userExists)
	if userExists.Email == auth.Email {
		var set = bson.M{"$set": bson.M{"email": auth.Email, "name": auth.Name, "phone": auth.Phone, "surname": auth.Surname}}
		database.DB(DBNAME).C("Users").Update(bson.M{"$or": []bson.M{{"email": auth.Email}, {"phone": auth.Phone}}}, set)
	}
	token := uuid.NewV4()

	var user User
	user.Id = token.String()
	user.UserCreated = time.Now().Unix()
	user.Name = auth.Name
	user.Email = auth.Email
	user.Phone = auth.Phone
	err = database.DB(DBNAME).C("Users").Insert(user)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error inserting user", err.Error()}, 500)
		return
	}

	//Save user credentials
	//Hash password
	passwd, _ := bcrypt.GenerateFromPassword([]byte(auth.Password), bcrypt.DefaultCost)
	auth.Password = string(passwd)
	err = database.DB(DBNAME).C("Credentials").Insert(auth)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error inserting user", err.Error()}, 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	var er = make([]string, 0)
	var resp = Resp{
		Code:     200,
		Errors:   er,
		Response: user,
	}
	json.NewEncoder(w).Encode(resp)

}
