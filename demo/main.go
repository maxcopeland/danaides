package main

import (
	"context"
	"danaides/internal/s3utils"
	"flag"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func CreateTextFile(name string, size int) error {
	b := make([]byte, size)

	err := os.WriteFile(name, b, 0644)

	if err != nil {
		return err
	}

	return nil

}

func main() {

	var bucketName, bucketDir string
	var nFiles, fileSize int

	flag.StringVar(&bucketName, "bucketname", "dummy-data-maxcope",
		"S3 destination bucket name (e.g. \"dummy-data-maxcope\"")
	flag.StringVar(&bucketDir, "bucketdir", "data",
		"Destination folder in S3 bucket (e.g. \"data\" or \"tmp/test\"")
	flag.IntVar(&nFiles, "n", 10, "Number of files to create")
	flag.IntVar(&fileSize, "size", 1, "Size of each file (in MiB)")

	flag.Parse()

	//fmt.Printf("bucketname: %s\n", bucketName)
	//fmt.Printf("buckdir: %s\n", bucketDir)
	//fmt.Printf("n: %d\n", nFiles)
	//fmt.Printf("size: %d\n", fileSize)
	//fmt.Println("ya")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile("mle"), // Only necessary for non-default profile
		config.WithRegion("us-east-1"),        // Only necessary for non-default profile
	)

	if err != nil {
		log.Fatalf("error loading config: %v, %v", cfg, err)
	}

	mgr := s3utils.BucketManager{
		S3Client: s3.NewFromConfig(cfg),
	}

	dir, err := os.MkdirTemp("", "dummy_data")

	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	zeroPadding := len(strconv.Itoa(nFiles)) - 1

	for i := 0; i < nFiles; i++ {

		done := make(chan struct{})

		go func() {
			fName := fmt.Sprintf("tmpfile%0*d", zeroPadding, i)
			absPath := filepath.Join(dir, fName)

			err = CreateTextFile(absPath, fileSize*1024*1024)

			if err != nil {
				log.Fatal(err)
			}

			objectKey := filepath.Join(bucketDir, fmt.Sprintf("test/file%0*d", zeroPadding, i))

			err = mgr.UploadFile(bucketName, objectKey, absPath)
			if err != nil {
				log.Fatalf("error uploading %v to s3://%v/%v, %v", absPath, bucketName, objectKey, err)
			}
			done <- struct{}{}

		}()

		<-done

	}
	fmt.Println("success")

}
