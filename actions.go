package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"path"
	"syscall"

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
		name:     "copy file",
		function: copy,
	},
	{
		name:     "print path",
		function: print,
	},
	{
		name:     "path to clipboard",
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

	size, err := getFileSize(state.bucket, state.key)
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting download, size:", byteCountDecimal(size))

	cwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	temp, err := ioutil.TempFile(cwd, "getObjWithProgress-tmp-")
	if err != nil {
		panic(err)
	}
	tempfileName := temp.Name()

	// Handle sigint and cleanup temp files
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Remove(tempfileName)
		os.Exit(0)
	}()

	writer := &progressWriter{writer: temp, size: size, written: 0}

	downloader, err := createDownloadServiceForBucket(&state.bucket)
	if err != nil {
		return err
	}

	_, err = downloader.Download(writer, &s3.GetObjectInput{
		Bucket: aws.String(state.bucket),
		Key:    aws.String(state.key),
	})
	if err != nil {
		fmt.Printf("Download failed! Deleting tempfile: %s", tempfileName)
		os.Remove(tempfileName)
		return fmt.Errorf("failed to download file, %v", err)
	}

	if err := temp.Close(); err != nil {
		panic(err)
	}

	if err := os.Rename(temp.Name(), filename); err != nil {
		panic(err)
	}

	os.Chmod(filename, 0644)

	fmt.Println("File downloaded! Available at:", filename)

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
