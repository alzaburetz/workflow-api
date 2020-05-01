package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"gopkg.in/mgo.v2"
)

var (
	DBNAME     string
	DBUSER     string
	DBPASSWORD string
	DBADDR     string
)
var database *mgo.Session

type DataBase struct {
	db *mgo.Session
}

type Resp struct {
	Code     int         `json:"code"`
	Errors   []string    `json:"errors"`
	Response interface{} `json:"response"`
}

func InitDatabase(session *mgo.Session) {
	database = session
}

func GetPushService() string {
	return os.Getenv("FIREBASE")
}

func CreateDatabaseInstance() {
	DBNAME = os.Getenv("DBNAME")
	DBUSER = os.Getenv("DBUSER")
	DBPASSWORD = os.Getenv("DBPASSWORD")
	DBADDR = os.Getenv("DBADDR")
	dialinfo := &mgo.DialInfo{
		Addrs:    []string{DBADDR},
		Database: DBNAME,
		Username: DBUSER,
		Password: DBPASSWORD,
	}
	var err error

	if database, err = mgo.DialWithInfo(dialinfo); err != nil {
		os.Exit(1)
	}
}

func AccessDataStore() *mgo.Session {
	if database == nil {
		CreateDatabaseInstance()
	}
	return database.Copy()
}

func WriteAnswer(w *http.ResponseWriter, msg interface{}, httperrors []string, code int) {
	var response = Resp{
		Code:     code,
		Errors:   httperrors,
		Response: msg,
	}
	json.NewEncoder(*w).Encode(response)
}
