package s3

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func DownloadFile(chunk int64, conc int, bucketName, keyName, path string) {
	cfg, err := config.LoadDefaultConfig(context.TODO())

	if err != nil {
		panic(err)
	}

	filename := keyName[strings.LastIndex(keyName, "/")+1:]

	file, err := os.Create(filepath.Join(path, filename))

	if err != nil {
		panic(err)
	}

	defer file.Close()

	client := s3.NewFromConfig(cfg)

	fmt.Println("Chunk size", chunk)
	fmt.Println("Concurrency", conc)

	downloader := manager.NewDownloader(client, func(d *manager.Downloader) {
		d.PartSize = chunk * 1024 * 1024 // 8MB
		d.Concurrency = conc
	})

	start := time.Now()
	numBytes, err := downloader.Download(context.TODO(), file, &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(keyName),
	})
	t := time.Now()
	elapsed := t.Sub(start)

	if err != nil {
		panic(err)
	}

	fmt.Println("Elapsed time: ", elapsed)
	fmt.Println("Num bytes: ", numBytes)
	var numMB float64 = float64(numBytes) / (1024 * 1024)
	var numSeconds float64 = float64(elapsed) / 1e9

	fmt.Printf("Download rate: %.2f MB/s\n", numMB/numSeconds)
}
