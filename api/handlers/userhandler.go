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
	} else if len(ua.Password) < 5{
		return errors.New("Password is invalid or missing")
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

//Gets user by token
func GetUser(w http.ResponseWriter, r *http.Request) {
	if err, userKey := middleware.CheckToken(r); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		WriteAnswer(&w, err.Error(), []string{"Wrong token, try relogin"}, 403)
	} else {
		w.WriteHeader(http.StatusOK)

		var user User
		var db = AccessDataStore()
		db.db.DB("app").C("Users").Find(bson.M{"email":userKey}).One(&user)
		WriteAnswer(&w, user, []string{}, 200)
	}
}

//Handle register
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
			
			//This part checks if user aready exists
			var userExists User
			database.DB("app").C("Users").Find(bson.M{"$or" :[]bson.M{ bson.M{"email": auth.Email}, bson.M{"phone":auth.Phone}}}).One(&userExists)
			if userExists.Email != "" || userExists.Phone != "" { //if user is found, return error
				w.WriteHeader(http.StatusBadRequest)
				WriteAnswer(&w, "", []string {"User already exists"},400)
				return
			} else { //if user is new, we save it
				var user User
				user.Id, _ = database.DB("app").C("Users").Count()
				user.UserCreated = time.Now()
				user.Name = auth.Name
				user.Email = auth.Email
				user.Phone = auth.Phone
				database.DB("app").C("Users").Insert(user)
	
	
				//Save user credentials 
				//Hash password
				passwd, _ := bcrypt.GenerateFromPassword([]byte(auth.Password), bcrypt.DefaultCost)
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
}


//Handles login
func Login(w http.ResponseWriter, r *http.Request) {
	var auth UserAuth
	var resp Resp
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, "", []string {"Couldn't read data from post form", "Expecing json"}, 400)
		return
	}

	if err = json.Unmarshal(data, &auth); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, "", []string{"Couldn't unmarshal data", err.Error(), "Expecting 'email', 'password' fields"}, 400)
		return
	}

	var userExists UserAuth
	database.DB("app").C("Credentials").Find(bson.M{"email": auth.Email}).One(&userExists)
	if userExists.Email == "" {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, "", []string{"Couldn't login", "User doesn't exist"}, 400)
		return
	}

	

	err = bcrypt.CompareHashAndPassword([]byte(userExists.Password),[]byte(auth.Password))
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		WriteAnswer(&w, "", []string{"Password is incorrect"}, 401)
		return
	} else {
		token, err := middleware.CreateToken(auth.Email)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			WriteAnswer(&w, "", []string{"Error creating token"}, 500)
			return
		} else {
			resp.Code= 200
			resp.Response = token
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(resp)
		}

	}

}