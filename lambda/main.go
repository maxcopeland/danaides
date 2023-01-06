package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

// S3BucketEvent is the event that is passed to the Lambda function
type S3BucketEvent struct {
	Bucket string `json:"bucket"`
}

func listS3BucketContents(ctx context.Context, event S3BucketEvent) (string, error) {
	// Create an AWS session
	sess, err := session.NewSession()
	if err != nil {
		return "", fmt.Errorf("Error creating AWS session: %s", err)
	}

	// Create an S3 client
	svc := s3.New(sess)

	// List the objects in the bucket
	params := &s3.ListObjectsInput{
		Bucket: aws.String(event.Bucket),
	}
	resp, err := svc.ListObjects(params)
	if err != nil {
		return "", fmt.Errorf("Error listing objects in S3 bucket: %s", err)
	}

	// Print the object keys
	for _, item := range resp.Contents {
		fmt.Println("Name:         ", *item.Key)
		fmt.Println("Last modified:", *item.LastModified)
		fmt.Println("Size:         ", *item.Size)
		fmt.Println("Storage class:", *item.StorageClass)
		fmt.Println("")
	}

	return "Success", nil
}

func main() {
	lambda.Start(listS3BucketContents)
}
