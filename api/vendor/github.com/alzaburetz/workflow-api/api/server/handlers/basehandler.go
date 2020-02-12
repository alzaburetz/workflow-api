package handlers

import ("gopkg.in/mgo.v2"
		"net/http"
		"log"
		"os"
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


  func CreateDatabaseInstance() {
	dialinfo := &mgo.DialInfo{
		Addrs: []string{"ds163517.mlab.com:63517"},
		Database: "heroku_gwrf0w5w",
		Username: "admin",
		Password: "main123",
	}
	var err error

	if database, err = mgo.DialWithInfo(dialinfo); err != nil {
		log.Println("AAAAAAAAAAAA")
		log.Println(err.Error())
		os.Exit(1)
	}
  }
  

func AccessDataStore() *mgo.Session {
	if database == nil {
		CreateDatabaseInstance()
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
