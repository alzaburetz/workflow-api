package handlers

import "gopkg.in/mgo.v2"

var database *mgo.Session


type Resp struct {
	Code int `json:"code"`
	Errors []string `json:"errors"`
	Response interface{} `json:"response"`
}

type Msg struct {
	Message string `json:"message"`
	Reason string `json:"reason"`
}

type DataValidity interface {
	NotNil() bool
	HasRequiredFields() (bool, error)
}

func InitDatabase(session *mgo.Session) {
	database = session
}
