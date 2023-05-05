package main

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func getBucketsListAWS() ([]string, error) {
	AwsSession, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		err = handAwsError(err)
		return nil, err
	}

	s3Service := s3.New(AwsSession)

	input := &s3.ListBucketsInput{}

	result, err := s3Service.ListBuckets(input)
	if err != nil {
		err = handAwsError(err)
		return nil, err
	}

	bucketNames := []string{}

	for _, b := range result.Buckets {
		bucketNames = append(bucketNames, *b.Name)
	}

	return bucketNames, nil
}

func getObjectsListAWS(bucketName string, prefix string) ([]ListObj, error) {

	s3Service, err := createS3ServiceForBucket(&bucketName)
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
		err = handAwsError(err)
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

func createS3Service(region *string) (*s3.S3, error) {
	AwsSession, err := session.NewSession(&aws.Config{
		Region: region,
	})
	if err != nil {
		return nil, err
	}

	s3Service := s3.New(AwsSession)

	return s3Service, nil
}

func createS3ServiceForBucket(bucketName *string) (*s3.S3, error) {
	region, err := findBucketRegion(bucketName)
	if err != nil {
		return nil, err
	}

	s3Service, err := createS3Service(region)
	if err != nil {
		return nil, err
	}

	return s3Service, nil
}

func createDownloadService(region *string) (*s3manager.Downloader, error) {
	AwsSession, err := session.NewSession(&aws.Config{
		Region: region,
	})
	if err != nil {
		return nil, err
	}

	s3Downloader := s3manager.NewDownloader(AwsSession)

	return s3Downloader, nil
}

func createDownloadServiceForBucket(bucketName *string) (*s3manager.Downloader, error) {

	region, err := findBucketRegion(bucketName)
	if err != nil {
		return nil, err
	}

	s3Downloader, err := createDownloadService(region)
	if err != nil {
		return nil, err
	}

	return s3Downloader, nil
}

func findBucketRegion(bucketName *string) (*string, error) {
	AwsSession, err := session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	})
	if err != nil {
		return nil, err
	}

	s3Service := s3.New(AwsSession)

	bucketRegion, err := s3Service.GetBucketLocation(&s3.GetBucketLocationInput{
		Bucket: bucketName,
	})
	if err != nil {
		err = handAwsError(err)
		return nil, err
	}

	// This is null for us-east-1.  Maybe a better way out there
	region := bucketRegion.LocationConstraint
	if region == nil {
		region = aws.String("us-east-1")
	}

	return region, nil
}

func handAwsError(err error) error {
	if err != nil {
		if awserr, ok := err.(awserr.Error); ok {
			if awserr == aws.ErrMissingRegion {
				newError := fmt.Errorf("no region configuration found. set AWS_REGION or configure aws cli")
				return newError
			}
		}
	}
	return err
}
