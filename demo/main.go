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

	bucketName := *flag.String("bucketname", "dummy-data-maxcope",
		"S3 destination bucket name (e.g. \"dummy-data-maxcope\"")
	bucketDir := *flag.String("bucketdir", "data",
		"Destination folder in S3 bucket (e.g. \"data\" or \"tmp/test\"")
	nFiles := *flag.Int("n", 5, "Number of files to create")
	fileSize := *flag.Int("size", 1, "Size of each file (in MiB)")

	flag.Parse()

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithSharedConfigProfile("mle"), // Only necessary for non-default profile
		config.WithRegion("us-east-1"),        // Only necessary for non-default profile
	)

	if err != nil {
		log.Fatalf("error loading config: %v, %v", cfg, err)
	}

	client := s3.NewFromConfig(cfg)

	mgr := s3utils.BucketManager{
		S3Client: client,
	}

	dir, err := os.MkdirTemp("", "dummy_data")

	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)

	for i := 0; i < nFiles; i++ {

		absPath := filepath.Join(dir, "tmpfile")

		err = CreateTextFile(absPath, fileSize*1024*1024)

		if err != nil {
			log.Fatal(err)
		}

		objectKey := filepath.Join(bucketDir, fmt.Sprintf("test/file%d", i))

		err = mgr.UploadFile(bucketName, objectKey, absPath)

		if err != nil {
			log.Fatalf("error uploading %v to s3://%v/%v, %v", absPath, bucketName, objectKey, err)
		}

	}
	fmt.Println("success")

}
