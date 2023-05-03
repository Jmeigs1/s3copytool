package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var DEFAULT_REGION string = "us-east-1"

func getBucketsListAWS() ([]string, error) {

	var err error
	AwsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(DEFAULT_REGION),
	})
	if err != nil {
		return nil, err
	}

	s3Service := s3.New(AwsSession)

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

func getObjectsListAWS(bucketName string, prefix string) ([]ListObj, error) {

	s3Service, err := createS3ServiceForBucket(bucketName)
	if err != nil {
		return nil, err
	}

	input := &s3.ListObjectsV2Input{
		Bucket:    aws.String(bucketName),
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

func createS3Service(region string) (*s3.S3, error) {
	AwsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}

	s3Service := s3.New(AwsSession)

	return s3Service, nil
}

func createS3ServiceForBucket(bucketName string) (*s3.S3, error) {
	s3Service, err := createS3Service(DEFAULT_REGION)
	if err != nil {
		return nil, err
	}

	bucketRegion, err := s3Service.GetBucketLocation(&s3.GetBucketLocationInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return nil, err
	}

	if *bucketRegion.LocationConstraint != DEFAULT_REGION {
		s3Service, err = createS3Service(*bucketRegion.LocationConstraint)
		if err != nil {
			return nil, err
		}
	}

	return s3Service, nil
}

func createDownloadService(region string) (*s3manager.Downloader, error) {
	AwsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err != nil {
		return nil, err
	}

	s3Service := s3manager.NewDownloader(AwsSession)

	return s3Service, nil
}

func createDownloadServiceForBucket(bucketName string) (*s3manager.Downloader, error) {

	s3Service, err := createS3Service(DEFAULT_REGION)
	if err != nil {
		return nil, err
	}

	bucketRegion, err := s3Service.GetBucketLocation(&s3.GetBucketLocationInput{
		Bucket: aws.String(bucketName),
	})
	if err != nil {
		return nil, err
	}

	AwsSession, err := session.NewSession(&aws.Config{
		Region: bucketRegion.LocationConstraint,
	})
	if err != nil {
		return nil, err
	}

	s3Downloader := s3manager.NewDownloader(AwsSession)

	return s3Downloader, nil
}
