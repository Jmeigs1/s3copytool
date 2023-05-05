package main

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

const FileNameExample = "go-example.txt"

// Helper functions
func joinPrefixes(prefixes []string) string {
	if len(prefixes) == 0 {
		return ""
	}
	return strings.Join(prefixes, "/") + "/"
}

func getFileSize(bucketName string, prefix string) (int64, error) {

	s3Service, err := createS3ServiceForBucket(&bucketName)
	if err != nil {
		return 0, err
	}

	params := &s3.HeadObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(prefix),
	}

	resp, err := s3Service.HeadObject(params)
	if err != nil {
		return 0, err
	}

	return *resp.ContentLength, nil
}

func byteCountDecimal(b int64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(b)/float64(div), "kMGTPE"[exp])
}
