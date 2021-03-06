package post

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	_ "github.com/alzaburetz/workflow-api/api/server/handlers/user/filehandlers"
	. "github.com/alzaburetz/workflow-api/api/server/middleware"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
)

func AddPost(w http.ResponseWriter, r *http.Request) {
	groupid, _ := mux.Vars(r)["id"]

	var post Post
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"Error reading data from request body", err.Error()}, 400)
		return
	}

		if err = json.Unmarshal(body, &post); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			WriteAnswer(&w, nil, []string{"Error reading json", err.Error()}, 400)
			return
		}

	if err = post.HasRequiredFields(); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"Post error", err.Error()}, 400)
		return
	}

	var database = AccessDataStore()
	defer database.Close()

	err, email := CheckToken(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error getting token", err.Error()}, 500)
		return
	}
	if err = database.DB(DBNAME).C("Users").Find(bson.M{"email": email}).One(&post.Author); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Database error", "Error getting user", err.Error()}, 500)
		return
	}
	post.Author.Schedule = nil
	post.Author.Groups = nil

	post.GroupID = groupid
	post.Timestamp = time.Now().Unix()
	//count, _ := database.DB(DBNAME).C("Posts").Find(bson.M{"group_id":id}).Count()
	PostUUID := uuid.NewV4()
	post.Id = PostUUID.String()
	if err = database.DB(DBNAME).C("Posts").Insert(post); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error inserting data to the collection"}, 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, post, []string{}, 200)
}
