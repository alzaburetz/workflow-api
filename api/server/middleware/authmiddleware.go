package middleware

import ("net/http"
		"errors"
		//"time"
		"fmt"
		"log"
		"github.com/satori/go.uuid"
		"github.com/gomodule/redigo/redis")



// func RedisInit() error {
// 	//var duration = time.Second * 5
// 	con, err := redis.DialURL("redis://redistogo:c7ec584512cad0331e2d71355fadb333@pike.redistogo.com:10201/")
// 	//con, err := redis.DialURL("redis://h:p7f2f70abac018d527c3476ba62ed68ff8aaab05a7376c6f994ad6fd6a5fcce5a@ec2-34-200-118-77.compute-1.amazonaws.com:11849")
	
// 	return err
// }

func CreateToken(login string) (string, error) {
	
	token := uuid.NewV4()
	conn, err := AccessRedis()
	if err != nil {
		return "" , err
	}
	defer conn.Close()
	_, err = conn.Do("SET", token, login)
	return token.String(), err
}

func UpdateToken(login, token string) error {
	
	conn, err := AccessRedis()
	if err != nil {
		return err
	}
	defer conn.Close()
	_, err  = conn.Do("SET", token, login)
	return err
}

func CheckToken(r *http.Request) (error, string) {
	token := r.Header.Get("Token")
	conn, err := AccessRedis()
	if err != nil {
		return err, ""
	}
	defer conn.Close()
	if (len(token) > 0) {
		user, _ := conn.Do("GET",token)
		return nil, fmt.Sprintf("%s", user)
	} else {
		return errors.New("Token invalid"), ""
	}
}

func AccessRedis() (redis.Conn, error){
		con, err := redis.DialURL("redis://redistogo:c7ec584512cad0331e2d71355fadb333@pike.redistogo.com:10201/")
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		return con, nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/user/login" || r.URL.Path == "/api/user/register" {
			next.ServeHTTP(w,r)
			return
		}
		token := r.Header.Get("Token")
		
		if (len(token) == 0) {
			http.Error(w, http.StatusText(403), 403)
			return
		} else {
			next.ServeHTTP(w,r)
		}
	})
}