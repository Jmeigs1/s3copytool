package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func getBucketsListAWS(s3Service *s3.S3) ([]string, error) {
	input := &s3.ListBucketsInput{}

	result, err := s3Service.ListBuckets(input)
	if err != nil {
		return nil, err
	}

	bucketNames := []string{}

	for _, b := range result.Buckets {
		bucketNames = append(bucketNames, *b.Name)
	}

	return bucketNames, nil
}

func getObjectsListAWS(s3Service *s3.S3, bucket string, prefix string) ([]ListObj, error) {
	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucket),
		Prefix:    aws.String(prefix),
		Delimiter: aws.String("/"),
	}

	result, err := s3Service.ListObjectsV2(input)
	if err != nil {
		return nil, err
	}

	resultObjects := []ListObj{}

	for _, o := range result.Contents {
		resultObjects = append(resultObjects, ListObj{
			Name:   *o.Key,
			IsDir:  false,
			IsBack: false,
		})
	}

	for _, o := range result.CommonPrefixes {
		resultObjects = append(resultObjects, ListObj{
			Name:   *o.Prefix,
			IsDir:  true,
			IsBack: false,
		})
	}

	back := ListObj{
		Name:   prefix + "..",
		IsDir:  false,
		IsBack: true,
	}

	resultObjects = append([]ListObj{back}, resultObjects...)

	return resultObjects, nil
}
