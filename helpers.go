package main

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

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

type ParsedS3Url struct {
	Bucket   string
	Prefixes []string
	Key      string
}

func ParseS3Url(s3Url string) (*ParsedS3Url, error) {
	parsed, err := url.Parse(s3Url)
	if err != nil {
		return nil, err
	}

	if parsed.Scheme != "s3" {
		return nil, fmt.Errorf("Starting url is not a valid s3 path: %s", s3Url)
	}

	// Remove leading slash
	path := parsed.Path[1:]

	splitPath := strings.Split(path, "/")

	retVal := &ParsedS3Url{
		Bucket:   parsed.Host,
		Prefixes: splitPath[:len(splitPath)-1],
		Key:      "",
	}

	// assume "directory" paths will end in / and last index will be ""
	if splitPath[len(splitPath)-1] != "" {
		retVal.Key = path
	}

	return retVal, nil
}
