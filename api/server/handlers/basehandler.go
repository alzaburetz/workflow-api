package handlers

import ("gopkg.in/mgo.v2"
		"net/http"
		"log"
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
	var err error
	if database == nil {
		database, err = mgo.Dial("mongodb://admin:main123@ds163517.mlab.com:63517/heroku_gwrf0w5w")
		if err != nil {
			log.Println("AAAAAAAAAAAAAAAAAAAAAAAAAAAA BLYA" + err.Error())
			return nil
		}
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
