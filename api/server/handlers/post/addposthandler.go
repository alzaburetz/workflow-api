package post

import ("net/http"
		. "app/server/handlers"
		. "app/server/middleware"
		"strconv"
		"gopkg.in/mgo.v2/bson"
		"encoding/json"
		"io/ioutil"
		"time"
		"github.com/gorilla/mux")

func AddPost(w http.ResponseWriter, r *http.Request) {
	groupid, _ := mux.Vars(r)["id"]
	id, err := strconv.Atoi(groupid)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"Error getting group id", err.Error()}, 400)
		return
	}

	var post Post 
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"Error reading data from request body", err.Error()}, 400)
		return
	}

	if err = json.Unmarshal(body,&post); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"Error reading json", err.Error()},400)
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
	if err = database.DB("app").C("Users").Find(bson.M{"email":email}).One(&post.Author); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Database error", "Error getting user", err.Error()}, 500)
		return
	}

	post.GroupID = id
	post.Timestamp = time.Now().Unix()
	count, _ := database.DB("app").C("Posts").Find(bson.M{"group_id":id}).Count()
	post.Id = count

	if err = database.DB("app").C("Posts").Insert(post); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error inserting data to the collection"}, 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, post, []string{}, 200)
}