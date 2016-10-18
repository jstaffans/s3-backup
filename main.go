package main

import (
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"io/ioutil"
	"os"
	"time"
)

func main() {
	accessKeyPtr := flag.String("accessKey", "", "Access key")
	secretKeyPtr := flag.String("secretKey", "", "Secret key")
	bucketPtr := flag.String("bucket", "", "Bucket")
	folderPtr := flag.String("folder", "backups", "Folder")

	flag.Parse()

	folderToBackup := flag.Args()[0]

	var err error
	var sess *session.Session

	if *accessKeyPtr != "" && *secretKeyPtr != "" {
		sess, err = session.NewSession(&aws.Config{
			Region:      aws.String("eu-central-1"),
			Credentials: credentials.NewStaticCredentials(*accessKeyPtr, *secretKeyPtr, ""),
		})
	} else {
		sess, err = session.NewSession(&aws.Config{
			Region: aws.String("eu-central-1"),
		})
	}

	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	files, err := ioutil.ReadDir(folderToBackup)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	svc := s3.New(sess)

	for _, fileInfo := range files {
		if !fileInfo.IsDir() {
			file, _ := os.Open(fmt.Sprintf("%s/%s", folderToBackup, fileInfo.Name()))

			target := fmt.Sprintf("%s/%s/%s",
				*folderPtr,
				time.Now().Format(time.RFC3339),
				fileInfo.Name())

			_, err := svc.PutObject(&s3.PutObjectInput{
				Bucket: bucketPtr,
				Key:    &target,
				Body:   file,
			})

			if err != nil {
				fmt.Println(err.Error())
				return
			}
		}
	}
}
