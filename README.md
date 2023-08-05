# S3 Copy Tool

A simple golang cli tool to easily browse and download S3 Objects.

`s3copytool` invokes AWS API call via aws-sdk-go

All of the same environment variables and settings that power the AWS CLI are required/used for this application

For more info see: https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-configure.html

## Usage

General browsing of all account buckets

`s3copytool`

Browsing starting at a given s3 URI

`s3copytool s3://official-org-bucket/very-important-data/`

`s3copytool s3://official-org-bucket/very-important-data/results.csv`
