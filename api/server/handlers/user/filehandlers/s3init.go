package filehandlers

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"log"
	"os"
)

var sess *session.Session
const bucket = "https://workflow-2020-filestorage.s3.eu-central-1.amazonaws.com/"
const region = "eu-central-1"

func S3Init() {
	var err error
	awssecret := os.Getenv("AWSSECRET")
	awsid := os.Getenv("AWSID")
	sess, err = session.NewSession(&aws.Config{
		Credentials:                       credentials.NewStaticCredentials(awsid,awssecret, ""),
		Region:                            aws.String(region),
	})
	if err != nil {
		log.Fatal(err.Error())
	}
}
