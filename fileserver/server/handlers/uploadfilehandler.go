package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// "gopkg.in/mgo.v2/bson"
// . "fileserver/server")

func UploadFile(w http.ResponseWriter, r *http.Request) {
	user := r.Header.Get("User")
	if len(user) == 0 {
		fmt.Fprintf(w, "No user in handler")
		return
	}

	r.ParseMultipartForm(10 << 20)
	file, _, err := r.FormFile("File")
	if err != nil {
		return
	}
	defer file.Close()
	fileName := user + string(time.Now().Unix())
	_ = os.Mkdir("./static", 0777)
	f, err := os.OpenFile("./static/"+fileName+".png", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()
	io.Copy(f, file)

}
