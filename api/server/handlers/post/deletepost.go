package post

import (
	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	. "github.com/alzaburetz/workflow-api/api/server/middleware"
	"github.com/gorilla/mux"
	"gopkg.in/mgo.v2/bson"
	"net/http"
)

func DeletePost(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	group, _ := vars["id"]
	post, _ := vars["post"]

	var database = AccessDataStore()
	defer database.Close()

	var post_ Post

	post_id_group := []bson.M{bson.M{"group_id": group}, bson.M{"_id_": post}}

	var and_operator = bson.M{"$and": post_id_group}

	//var match = bson.M{"$match":and_operator}

	if err := database.DB("app").C("Posts").Find(and_operator).One(&post_); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error getting post from database", err.Error()}, 500)
		return
	}

	err, email := CheckToken(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error checking token", err.Error()}, 500)
		return
	}

	if post_.Author.Email != email {
		w.WriteHeader(http.StatusUnauthorized)
		WriteAnswer(&w, nil, []string{"You are not the owner of the post"}, 401)
		return
	} else {
		if err = database.DB("app").C("Posts").Remove(and_operator); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			WriteAnswer(&w, nil, []string{"Database error", err.Error()}, 500)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, "Successfully removed post", []string{}, 200)
}
