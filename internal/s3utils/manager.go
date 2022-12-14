package s3utils

import (
	"bytes"
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"io"
	"log"
	"os"
)

type BucketManager struct { // TODO: Need better name for this obj
	S3Client *s3.Client
	// TODO: Should bucket name be tracked here instead of passed to methods?
}

func (mgr BucketManager) UploadFile(bucketName, objectKey, fileName string) error {
	file, err := os.Open(fileName)

	if err != nil {
		log.Printf("Error uploading file %v... %v", fileName, err)
	}
	defer file.Close()

	_, err = mgr.S3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   file,
	})
	if err != nil {
		log.Printf("Error uploading file %v... %v", fileName, err)
	}

	return err
}

func (mgr BucketManager) DownloadFile(bucketName, objectKey, fileName string) error {
	result, err := mgr.S3Client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		log.Printf("Error downloading object %v:%v... %v", bucketName, objectKey, err)
	}
	defer result.Body.Close()

	file, err := os.Create(fileName)
	if err != nil {
		log.Printf("Error creating file %v... %v", fileName, err)
	}
	defer file.Close()

	body, err := io.ReadAll(result.Body)

	if err != nil {
		log.Printf("Error reading body from %v... %v", objectKey, err)
	}

	_, err = file.Write(body)
	return err
}

func (mgr BucketManager) UploadLargeObject(bucketName, objectKey string, largeObject []byte) error {
	largeBuffer := bytes.NewReader(largeObject)

	var partMiBs int64 = 10

	uploader := manager.NewUploader(mgr.S3Client, func(u *manager.Uploader) {
		u.PartSize = partMiBs * 1024 * 1024
	})

	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		Body:   largeBuffer,
	})

	if err != nil {
		log.Printf("Error uploading large file to %v:%v... %v", bucketName, objectKey, err)
	}

	return err
}

func (mgr BucketManager) DownloadLargeOject(bucketName string, objectKey string) ([]byte, error) {
	var partMiBs int64 = 10

	downloader := manager.NewDownloader(mgr.S3Client, func(d *manager.Downloader) {
		d.PartSize = partMiBs * 1024 * 1024
		// TODO: Add in concurrency?
	})

	buffer := manager.NewWriteAtBuffer([]byte{})
	_, err := downloader.Download(context.TODO(), buffer, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
	})

	if err != nil {
		log.Printf("Error downloading object %v:%v... %v\n", bucketName, objectKey, err)
	}

	return buffer.Bytes(), err
}

func (mgr BucketManager) ListAllObjects(bucketName string) ([]types.Object, error) {

	var contents []types.Object
	var err error

	var token *string

	for {

		results, err := mgr.S3Client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
			Bucket:            aws.String(bucketName),
			ContinuationToken: token,
		})

		if err != nil {
			log.Printf("Can't list objects in bucket %v. Here's where: %v", bucketName, err)
		} else {
			contents = append(contents, results.Contents...)
			token = results.NextContinuationToken

			if !results.IsTruncated {
				break
			}
		}
	}

	return contents, err
}
