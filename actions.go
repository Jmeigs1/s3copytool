package main

import (
	"fmt"
	"os"
	"path"

	"github.com/atotto/clipboard"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

type action struct {
	name     string
	function func(state *state) error
}

var actionList []action = []action{
	{
		name:     "..",
		function: back,
	},
	{
		name:     "copy",
		function: copy,
	},
	{
		name:     "print",
		function: print,
	},
	{
		name:     "pbcopy",
		function: pbcopy,
	},
}

func (a action) String() string {
	return a.name
}

func back(state *state) error {
	state.key = ""
	return nil
}

func copy(state *state) error {
	// Create a file to write the S3 Object contents to.
	filename := path.Base(state.key)
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file %q, %v", filename, err)
	}

	// Write the contents of S3 Object to the file
	n, err := state.downloadService.Download(f, &s3.GetObjectInput{
		Bucket: aws.String(state.bucket),
		Key:    aws.String(state.key),
	})
	if err != nil {
		return fmt.Errorf("failed to download file, %v", err)
	}
	fmt.Printf("file downloaded, %d bytes\n", n)

	back(state)
	return nil
}

func print(state *state) error {

	fmt.Printf("s3://%s/%s\n", state.bucket, state.key)

	os.Exit(0)
	return nil
}

func pbcopy(state *state) error {

	path := fmt.Sprintf("s3://%s/%s", state.bucket, state.key)

	clipboard.WriteAll(path)
	os.Exit(0)
	return nil
}

func doAction(state *state, action int) error {

	err := actionList[action].function(state)
	if err != nil {
		return err
	}

	return nil
}
