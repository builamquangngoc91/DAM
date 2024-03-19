package main

import (
	"context"
	"fmt"
	"net/http"
	"path/filepath"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/google/uuid"
)

func uploadFileHandler(w http.ResponseWriter, r *http.Request) {
	backgroundCtx := context.Background()
	sdkConfig, err := config.LoadDefaultConfig(
		backgroundCtx,
		config.WithRegion("us-east-1"),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("", "", "")))
	_, _ = sdkConfig, err

	s3Client := s3.NewFromConfig(sdkConfig)

	// Parse the multipart form, 10 << 20 specifies a maximum
	// upload of 10 MB files.

	// FormFile returns the first file for the given key `file`
	file, fileHeader, err := r.FormFile("file")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Create a new file in the server's upload directory
	_, err = s3Client.PutObject(backgroundCtx, &s3.PutObjectInput{
		Bucket: aws.String("anhdong1996"),
		Key:    aws.String(uuid.NewString() + filepath.Ext(fileHeader.Filename)),
		Body:   file,
	})
	if err != nil {
		fmt.Println("error", err)
	}

	fmt.Fprintf(w, "Successfully Uploaded File\n")
}

func main() {
	// backgroundCtx := context.Background()
	// sdkConfig, err := config.LoadDefaultConfig(
	// 	backgroundCtx,
	// 	config.WithRegion("us-east-1"),
	// 	config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("", "", "")))
	// _, _ = sdkConfig, err

	// s3Client := s3.NewFromConfig(sdkConfig)

	// output, err := s3Client.ListObjectsV2(backgroundCtx, &s3.ListObjectsV2Input{
	// 	Bucket: aws.String("anhdong1996"),
	// })
	// fmt.Println("Error", err)

	// if output != nil {
	// 	for _, object := range output.Contents {
	// 		fmt.Printf("Object: %s, Size: %d\n", *object.Key, *object.Size)
	// 	}
	// }

	// file, err := os.Open("./a.txt")
	// if err != nil {
	// 	log.Printf("Couldn't open file %v to upload. Here's why: %v\n", file.Name(), err)
	// } else {
	// 	defer file.Close()
	// 	_, err = s3Client.PutObject(backgroundCtx, &s3.PutObjectInput{
	// 		Bucket: aws.String("anhdong1996"),
	// 		Key:    aws.String("a.txt"),
	// 		Body:   file,
	// 	})
	// 	if err != nil {
	// 		fmt.Println("error", err)
	// 	}
	// }

	// s3PresignClient := s3.NewPresignClient(s3Client)
	// result, err := s3PresignClient.PresignGetObject(backgroundCtx, &s3.GetObjectInput{
	// 	Bucket: aws.String("anhdong1996"),
	// 	Key:    aws.String("a.txt"),
	// }, func(options *s3.PresignOptions) {
	// 	options.Expires = 3600 * time.Second
	// })
	// _ = err

	// fmt.Println("URL: ", result.URL)

	http.HandleFunc("/upload", uploadFileHandler)
	http.ListenAndServe(":8082", nil)
}
