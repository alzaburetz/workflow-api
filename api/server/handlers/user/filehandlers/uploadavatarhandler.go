package filehandlers

import (
	"bytes"
	. "github.com/alzaburetz/workflow-api/api/server/handlers"
	. "github.com/alzaburetz/workflow-api/api/server/handlers/user"
	"github.com/alzaburetz/workflow-api/api/server/middleware"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	uuid "github.com/satori/go.uuid"
	"gopkg.in/mgo.v2/bson"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func UploadAvatar(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(4096000)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"No file attached", err.Error()},400)
		return
	}

	file, fileHeader, _ := r.FormFile("avatar")
	fileHeaderBuffer := make([]byte,fileHeader.Size)
	file.Read(fileHeaderBuffer)

	var mime = http.DetectContentType(fileHeaderBuffer)
	if !strings.Contains(mime, "image") {
		w.WriteHeader(http.StatusBadRequest)
		WriteAnswer(&w, nil, []string{"You can only upload images!"}, 500)
		return
	}



	str := strings.Split(mime, "/")
	filetype := str[1]

	timestamp := strconv.Itoa(int(time.Now().Unix()))

	filename := uuid.NewV4().String() + timestamp + "." + filetype

	_, err = s3.New(sess).PutObject(&s3.PutObjectInput{
		Body: bytes.NewReader(fileHeaderBuffer),
		Bucket: aws.String("workflow-2020-filestorage"),
		Key: aws.String(filename),
		ACL: aws.String("public-read"),
	})


	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"error saving file", err.Error()}, 500)
		return
	}

	database := AccessDataStore()
	defer database.Close()

	_, email := middleware.CheckToken(r)
	var user User

	if err = database.DB(DBNAME).C("Users").Find(bson.M{"email":email}).One(&user); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		WriteAnswer(&w, nil, []string{"Error getting user", err.Error()}, 500)
		return
	}

	//If user alredy had avatar, delete previous
	if len(user.Avatar) > 0 {
		key := strings.TrimPrefix(user.Avatar, bucket)
		log.Println(key)
		request, _ := s3.New(sess).DeleteObjectRequest(&s3.DeleteObjectInput{
			Bucket:                    aws.String("workflow-2020-filestorage"),
			Key:                       aws.String(key),
		})

		go request.Send()
	}
	database.DB(DBNAME).C("Users").Update(bson.M{"email":email},bson.M{"$set":bson.M{"avatar": bucket + filename}})
	w.WriteHeader(http.StatusOK)
	WriteAnswer(&w, "Successfully uploaded image", []string{}, 200)
}
