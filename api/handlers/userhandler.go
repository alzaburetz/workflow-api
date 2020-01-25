package handlers

import ("net/http"
		"io/ioutil"
		"app/middleware"
		"errors"
		"time"
		"strings"
		"gopkg.in/mgo.v2/bson"
		"golang.org/x/crypto/bcrypt"
		"encoding/json")

type User struct {
	Id int `json:"id" bson:"_id_"`
	Name string `json:"name" bson:"name"`
	Surname string `json:"surname" bson:"surname"`
	Workdays int `json:"workdays" bson:"workdays"`
	Weekdays int `json:"weekdays" bson:"weekdays"`
	FirstWorkDay time.Time `json:"firstwork" bson:"firstwork"`
	UserCreated time.Time `json:"-" bson:"created"`
	Email string `json:"email" bson:"email"`
	Phone string `json:"phone" bson:"phone"`
}

type UserAuth struct {
	Email string `json:"email" bson:"email"`
	Phone string `json:"phone" bson:"phone"`
	Name string `json:"name" bson:"name"`
	Password string `json:"password" bson:"password"`
}

func (ua UserAuth)HasRequiredFields() error {
	if len(ua.Name) == 0 {
		return errors.New("Name is emprty")
	} else if len(ua.Email) == 0 || !strings.Contains(ua.Email, "@") {
		return errors.New("Email not valid")
	} else if len(ua.Phone) < 11 {
		return errors.New("Phone is invalid")
	} else {
		return nil
	}
}

func (u User)HasRequiredFields()  error {
	
	if len(u.Name) == 0 {
			return  errors.New("Name is empty")
	} else if len(u.Surname) == 0 {
			return errors.New("Surname is empty")
	} else if u.Workdays < 1 || u.Weekdays < 1 {
		return errors.New("Graph not filled correctly")
	} else {
		return nil
	}
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	if err := middleware.CheckToken(r); err == nil {
		json.NewEncoder(w).Encode(Msg{Message: "Hooray, you passed middleware"})
	} else {
		json.NewEncoder(w).Encode(Msg{Message: err.Error()})
	}
}


//TODO: Check for existing user before overwriting one
func RegisterUser(w http.ResponseWriter, r *http.Request) {
	var auth *UserAuth
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		errs := make([]string, 1)
		errs[0] = "Could not read request body"
		message := Resp {
			Code: 500,
			Errors: errs,
			Response: err.Error(),
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(message)
	} else {
		if err = json.Unmarshal(data, &auth); err != nil {
			errs := make([]string, 1)
			errs[0] = "Could not unmarshal json data"
			message := Resp {
				Code: 500,
				Errors: errs,
				Response: err.Error(),
			}
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(message)
		} else if err = auth.HasRequiredFields(); err != nil {
			errs := make([]string, 1)
			errs[0] = err.Error()
			message := Resp {
				Code: 400,
				Errors: errs,
				Response: nil,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(message)
		} else {
			//FULL SUCCESS GO AHEAD AND CREATE USER, HOORAY
			var user User
			user.Id, _ = database.DB("app").C("Users").Count()
			user.UserCreated = time.Now()
			user.Name = auth.Name
			user.Email = auth.Email
			user.Phone = auth.Phone
			database.DB("app").C("Users").Insert(user)


			//Save user credentials 
			//Hash password
			passwd, err := bcrypt.GenerateFromPassword([]byte(auth.Password), bcrypt.DefaultCost)
			if err != nil {
				
			}
			auth.Password = string(passwd)
			database.DB("app").C("Credentials").Insert(auth)

			w.WriteHeader(http.StatusOK)
			var er = make([]string, 0)
			var resp = Resp {
				Code: 200,
				Errors: er,
				Response: user,
			}
			json.NewEncoder(w).Encode(resp)
		}
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	var auth UserAuth
	var resp Resp
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Code = 400
		resp.Response = ""
		json.NewEncoder(w).Encode(resp)
		return
	}

	if err = json.Unmarshal(data, &auth); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		resp.Code = 400
		resp.Response = ""
		json.NewEncoder(w).Encode(resp)
		return
	}

	var userExists UserAuth
	database.DB("app").C("Credentials").Find(bson.M{"email": auth.Email}).One(&userExists)
	if userExists.Email == "" {
		er := make([]string, 1)
		er[0] = "User doesn't exist!"
		w.WriteHeader(http.StatusBadRequest)
		resp.Code = 400
		resp.Errors = er
		json.NewEncoder(w).Encode(resp)
		return
	}

	

	err = bcrypt.CompareHashAndPassword([]byte(userExists.Password),[]byte(auth.Password))
	if err != nil {
		er := make([]string, 1)
		er[0] = "Password is incorrect!"
		w.WriteHeader(http.StatusUnauthorized)
		resp.Code = 401
		resp.Errors = er
		json.NewEncoder(w).Encode(resp)
		return
	} else {
		token, err := middleware.CreateToken(auth.Email)
		if err != nil {
			er := make([]string, 1)
			er[0] = "Error creating token"
			w.WriteHeader(http.StatusInternalServerError)
			resp.Code = 500
			resp.Errors = er
			json.NewEncoder(w).Encode(resp)
			return
		} else {
			resp.Code= 200
			resp.Response = token
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		}

	}

}

//This function Handles creation of user from data provided via json from request
//For admin and testing purposes only
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
				err = user.HasRequiredFields()
				if err != nil {
					message := Msg {
						Message: "Wrong data provided",
						Reason: err.Error(),
					}
					json.NewEncoder(w).Encode(message)
				} else {
					user.Id, _ = database.DB("app").C("Users").Count()
					user.UserCreated = time.Now()
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
	var errors []string
	response := Resp {
		200,
		errors,
		&users,
	}
	json.NewEncoder(w).Encode(response)
}
