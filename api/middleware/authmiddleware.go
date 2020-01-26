package middleware

import ("net/http"
		"errors"
		"time"
		"fmt"
		"github.com/satori/go.uuid"
		"github.com/gomodule/redigo/redis")


var conn redis.Conn

func RedisInit() error {
	var duration = time.Second * 5
	con, err := redis.DialTimeout("tcp", "redis:6379", duration, duration, duration)
	conn = con
	return err
}

func CreateToken(login string) (string, error) {
	token, err := uuid.NewV4()
	
	_, err = conn.Do("SET", token, login)
	
	return token.String(), err
}

func CheckToken(r *http.Request) (error, string) {
	token := r.Header.Get("Token")
	if (len(token) > 0) {
		user, _ := conn.Do("GET",token)
		return nil, fmt.Sprintf("%s", user)
	} else {
		return errors.New("Token invalid"), ""
	}
}

// func AuthMiddleware(next http.Handler) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		token := r.Header.Get("Token")
// 		if (len(token) == 0) {
// 			http.Error(w, http.StatusText(403), 403)
// 			return
// 		} else {
// 			next.ServeHTTP(w,r)
// 		}
// 	}
// }