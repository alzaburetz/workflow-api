package handlers

import ("net/http"
		"io/ioutil"
		"errors"
		"encoding/json")

type User struct {
	Name string `json:"name" bson:"name"`
	Surname string `json:"surname" bson:"surname"`
}

func (u User)NotNil() bool {
	if (u == User{}) {
		return false
	} else {
		return true
	}
}

func (u User)HasRequiredFields() (bool, error) {
	if len(u.Name) == 0 || len(u.Surname) == 0{
		if len(u.Name) == 0 {
			return false, errors.New("Name is empty")
		} else {
			return false, errors.New("Surname is empty")
		}
	} else {
		return true, nil
	}
}

//This function Handles creation of user from data provided via json from request
func CreateUser(w http.ResponseWriter, r *http.Request) {
		data, err := ioutil.ReadAll(r.Body)
		//Handle error of reading data from request body
		if err != nil {
			message := Msg {
				Message: "Could not read data from request",
				Reason: err.Error(),
			}
			json.NewEncoder(w).Encode(message)
		} else {
			var user User
			if err = json.Unmarshal(data, &user); err != nil {
				message := Msg {
					Message: "Could not unmarshal json data",
					Reason: err.Error(),
				}
				json.NewEncoder(w).Encode(message)
			} else {
				//Check if data correct
				_ , err = user.HasRequiredFields()
				if err != nil {
					message := Msg {
						Message: "Wrong data provided",
						Reason: err.Error(),
					}
					json.NewEncoder(w).Encode(message)
				} else {
					database.DB("app").C("Users").Insert(user)
					message := Msg {
						Message: "Successfully inserted data",
						Reason: "",
					}
					json.NewEncoder(w).Encode(message)
				}
			}
		}
}

//This function handles listing of all users
func ListUsers(w http.ResponseWriter, r *http.Request) {
	users := []User{}
	if err := database.DB("app").C("Users").Find(nil).All(&users); err != nil {
		message := Msg {
			Message: "Error getting data",
			Reason: err.Error(),
		}
		json.NewEncoder(w).Encode(message)
	}
	//var errors []string
	//response := Resp {
	//	200,
	//	errors,
	//	&users,
	//}
	json.NewEncoder(w).Encode(users)
}
