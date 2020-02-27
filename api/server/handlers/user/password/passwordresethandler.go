package password

import (
	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	. "github.com/alzaburetz/workflow-api/api/server/handlers/user"
	"gopkg.in/mgo.v2/bson"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func ResetPassword(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	phone := query.Get("phone")
	if len(query) == 0 || len(phone) == 0  {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"No phone present"}, 400)
		return
	}

	var user UserAuth
	database := AccessDataStore()
	defer database.Close()
	log.Println("+"+phone)

	if err := database.DB(DBNAME).C("Credentials").Find(bson.M{"phone":"+" + strings.TrimSpace(phone)}).One(&user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error getting data from database", err.Error()}, 500)
		return
	}
	if len(user.Phone) == 0 {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"User doesn't exist"}, 400)
		return
	}

	code := rand.Intn(9999 - 1000) + 1000
	Storage[strings.TrimSpace(phone)] = strconv.Itoa(code)
	log.Println(code)

	login := os.Getenv("LOGINSMS")
	password := os.Getenv("PASSWORDSMS")
	message := strconv.Itoa(code)
	requeststr := "https://smsc.ru/sys/send.php?login=" + login + "&psw=" + password + "&phones=" + strings.TrimSpace(phone) + "&mes=" + message
	log.Println(requeststr)

	sendcode, err := http.Get(requeststr)
	if err != nil {
		w.WriteHeader(500)
		WriteAnswer(&w, nil, []string{"Error occured sending SMS", err.Error()}, 500)
		return
	}

	var resp []byte
	response, _ := sendcode.Body.Read(resp)

	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, response, []string{}, 200)

}
