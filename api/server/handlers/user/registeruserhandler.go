package user

import ("net/http"
		"io/ioutil"
		"time"
		"gopkg.in/mgo.v2/bson"
		. "github.com/alzaburetz/workflow-api/api/server/handlers"
		"github.com/satori/go.uuid"
		"golang.org/x/crypto/bcrypt"
		"encoding/json")

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
		return
	} 

	
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
			return
		}
		
	if err = auth.HasRequiredFields(); err != nil {
			errs := make([]string, 1)
			errs[0] = err.Error()
			message := Resp {
				Code: 400,
				Errors: errs,
				Response: nil,
			}
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(message)
			return
		}
		
			//FULL SUCCESS GO AHEAD AND CREATE USER, HOORAY
			
			//This part checks if user aready exists
			var userExists User
			var database = AccessDataStore()
			
			defer database.Close()
			database.DB("app").C("Users").Find(bson.M{"$or" :[]bson.M{ bson.M{"email": auth.Email}, bson.M{"phone":auth.Phone}}}).One(&userExists)
			if userExists.Email != "" || userExists.Phone != "" { //if user is found, return error
				w.WriteHeader(http.StatusBadRequest)
				WriteAnswer(&w, "", []string {"User already exists"},400)
				return
			} 
		token, err := uuid.NewV4()

				var user User
				user.Id = token.String()
				user.UserCreated = time.Now().Unix()
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
