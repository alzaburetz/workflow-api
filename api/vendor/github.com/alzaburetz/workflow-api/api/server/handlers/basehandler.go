package handlers

import ("gopkg.in/mgo.v2"
		"net/http"
		"encoding/json")

var database *mgo.Session

type DataBase struct {
	db *mgo.Session
}


type Resp struct {
	Code int `json:"code"`
	Errors []string `json:"errors"`
	Response interface{} `json:"response"`
}

func InitDatabase(session *mgo.Session) {
	database = session
}

func AccessDataStore() *mgo.Session {
	if database == nil {
		database, _ = mgo.Dial("mongo:27017")
	}
	return database.Copy()
}

func WriteAnswer(w *http.ResponseWriter, msg interface{}, httperrors []string, code int) error {
	var response = Resp {
		Code: code,
		Errors: httperrors,
		Response: msg,
	}
	return json.NewEncoder(*w).Encode(response)
}
