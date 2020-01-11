package main

import ("log"
		"io/ioutil"
		"net/http"
		"encoding/json"
		"gopkg.in/mgo.v2"
		"github.com/gorilla/mux")


type User struct {
	Name string
	Surname string
}

var col *mgo.Collection
func main() {
	session, err := mgo.Dial("mongo:27017")
	if err != nil {
		log.Println("ERROR CONNECTING TO DATABASE")
	}
	col = session.DB("app").C("test")
	r := mux.NewRouter()
	r.HandleFunc("/" , HandleFunc).Methods("GET")
	r.HandleFunc("/makeuser", TestDatabase).Methods("POST")
	r.HandleFunc("/deleteuser", DeleteTest).Methods("DELETE")
	http.Handle("/", r)
	http.ListenAndServe(":3000", nil)
	log.Println("Listening to port 3000")
}

func HandleFunc(w http.ResponseWriter, r *http.Request) {
	result := []User{}
	if err := col.Find(nil).All(&result); err != nil {
		json.NewEncoder(w).Encode(nil)
	} else {
		json.NewEncoder(w).Encode(result)
	}
	

}

func TestDatabase(w http.ResponseWriter, r *http.Request) {
	data, _ := ioutil.ReadAll(r.Body)
	var test User
	json.Unmarshal(data, &test)
	col.Insert(test)
}

func DeleteTest(w http.ResponseWriter, r *http.Request) {
	col.Remove(r.Body)
}