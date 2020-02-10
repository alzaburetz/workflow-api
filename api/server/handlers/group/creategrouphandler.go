package group

import ("net/http"
		"gopkg.in/mgo.v2/bson"
		"io/ioutil"
		"encoding/json"
		"github.com/satori/go.uuid"
		. "github.com/alzaburetz/workflow-api/api/server/middleware"
		. "github.com/alzaburetz/workflow-api/api/server/handlers"
		. "github.com/alzaburetz/workflow-api/api/server/handlers/user")


func CreateGroup(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{ "Error reading request body",err.Error()}, 500)
		return
	}

	var group Group

	if err = json.Unmarshal(body, &group); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{ "Error unmarshaling data",err.Error()}, 500)
		return
	}

	if err = group.HasRequiredFields(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{err.Error()}, 400)
		return
	}

	err, user := CheckToken(r) 
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error getting token", err.Error()}, 500)
		return
	}

	var database = AccessDataStore()
	defer database.Close()

	var creator User
	if err = database.DB("app").C("Users").Find(bson.M{"email":user}).One(&creator); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Database error", "Error checking user", err.Error()}, 500)
		return
	}


	var updateduser User
	updateduser = creator

	token, _ := uuid.NewV4()
	group.Id = token.String()
	updateduser.Groups = append(updateduser.Groups, group.Id)

	if err = database.DB("app").C("Users").Update(bson.M{"email":user}, bson.M{"$addToSet":bson.M{"groups":group.Id}}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Database error", "Error updating userdata", err.Error()}, 500)
		return
	}


	creator.Groups = nil
	group.Creator = creator
	group.UserCount = 1
	if err = database.DB("app").C("Groups").Insert(group); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Database error", "Error inserting data", err.Error()}, 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, group, []string{ }, 200)

} 

