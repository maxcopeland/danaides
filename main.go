package main

import (
	"danaides/internal/s3"
)

func main() {

	s3.DownloadFile(64,
		4,
		"dbgap-maxcope",
		"data/SRR1219898/SRR1219898",
		"/home/ec2-user/")
}
