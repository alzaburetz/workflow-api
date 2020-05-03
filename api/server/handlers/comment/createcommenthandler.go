package comment

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"time"

	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	. "github.com/alzaburetz/workflow-api/api/server/handlers/notification"
	. "github.com/alzaburetz/workflow-api/api/server/middleware"
	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
)

func CreateComment(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"Error reading request body", err.Error()}, 400)
		return
	}

	var comment Comment
	if err := json.Unmarshal(body, &comment); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"Error unmarshaling data", err.Error()}, 400)
		return
	}

	if !comment.BodyNotEmpty() {
		w.WriteHeader(http.StatusNoContent)
		WriteAnswer(&w, nil, []string{"Body is empty"}, 204)
		return
	}

	var vars = mux.Vars(r)
	postid := vars["post"]

	comment.PostId = postid
	comment.Created = time.Now().Unix()
	comment.Uid = uuid.NewV4().String()

	database := AccessDataStore()
	defer database.Close()

	err, user := CheckToken(r)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error reading token", err.Error()}, 500)
		return
	}

	if err := database.DB(DBNAME).C("Users").Find(bson.M{"email": user}).One(&comment.Author); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error getting user from database", err.Error()}, 500)
		return
	}

	comment.Author.Schedule = nil
	comment.Author.Groups = nil

	if len(comment.At.Email) > 0 {
		var notification Notification
		notification.Object = comment.Body
		notification.Seen = false
		notification.Uid = uuid.NewV4().String()
		notification.Created = time.Now().Unix()
		notification.UserEmail = comment.At.Email
		database.DB(DBNAME).C("Notifications").Insert(notification)
		go SendNotification(comment.At.PushToken, "Ответ на коментарий", comment.Body)
	}

	if err = database.DB(DBNAME).C("Comments").Insert(comment); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error inserting data into database", err.Error()}, 500)
		return
	}

	if err = database.DB(DBNAME).C("Posts").Update(bson.M{"_id_": postid}, bson.M{"$inc": bson.M{"comments": 1}}); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error updating record in the database", err.Error()}, 500)
		return
	}

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, comment, []string{}, 200)
}

func SendNotification(push, title, body string) {
	request, _ := json.Marshal(map[string]string{
		"title": title,
		"body":  body,
		"token": push,
	})
	_, _ = http.Post(GetPushService(), "application/json", bytes.NewBuffer(request))
}
