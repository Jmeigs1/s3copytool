package main

import (
	"fmt"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/common-nighthawk/go-figure"
	"github.com/manifoldco/promptui"
)

var awsRegion string = "us-east-1"

type ListObj struct {
	Name   string
	IsDir  bool
	IsBack bool
}

type state struct {
	session         *session.Session
	s3Service       *s3.S3
	downloadService *s3manager.Downloader
	bucket          string
	prefixes        []string
	key             string
	action          int
}

func main() {

	defer func() {
		if err := recover(); err != nil {
			castErr, ok := err.(error)
			if ok && castErr.Error() == "^C" {

			} else {
				fmt.Println("General error: ", err)
			}
		}
	}()

	myFigure := figure.NewFigure("S3 Copy Tool", "", true)
	myFigure.Print()

	fmt.Println("")

	state := state{
		bucket:   "",
		prefixes: []string{},
		key:      "",
		action:   -1,
	}

	state.session = session.New(&aws.Config{
		Region: aws.String(awsRegion),
	})

	state.s3Service = s3.New(state.session)

	state.downloadService = s3manager.NewDownloader(state.session)

	for {
		if state.bucket == "" {
			err := setBucketState(&state)
			if err != nil {
				panic(err)
			}
		} else if state.key == "" {
			err := setObjectState(&state)
			if err != nil {
				panic(err)
			}
		} else if state.action < 0 {
			action, err := setAction(&state)
			if err != nil {
				panic(err)
			}
			doAction(&state, action)
		}
	}
}

// State functions

func setBucketState(state *state) error {

	buckets, err := getBucketsListAWS(state.s3Service)
	if err != nil {
		return err
	}

	prompt := promptui.Select{
		Label: "Select Bucket",
		Items: buckets,
		Searcher: func(input string, index int) bool {
			return strings.Contains(buckets[index], input)
		},
		StartInSearchMode: true,
		Stdout:            &bellSkipper{},
	}

	_, state.bucket, err = prompt.Run()

	if err != nil {
		return err
	}

	return nil
}

func setObjectState(state *state) error {

	objects, err := getObjectsListAWS(state.s3Service, state.bucket, joinPrefixes(state.prefixes))
	if err != nil {
		return err
	}

	prompt := promptui.Select{
		Label: "",
		Items: objects,
		Searcher: func(input string, index int) bool {

			lowerInput := strings.ToLower(input)
			lowerName := strings.ToLower(objects[index].Name)

			return strings.Contains(lowerName, lowerInput)
		},
		StartInSearchMode: true,
		Templates: &promptui.SelectTemplates{
			Inactive: fmt.Sprintf("{{if .IsDir}} {{ .Name | bold }} {{- else}} {{ .Name }} {{- end}}"),
			Active:   fmt.Sprintf("%s {{if .IsDir}} {{ .Name | bold }} {{- else}} {{ .Name }} {{- end}}", promptui.IconSelect),
			Selected: fmt.Sprintf(`{{ "%s" | green }} {{ "s3://%s/" | faint }}{{.Name | faint}}`, promptui.IconGood, state.bucket),
		},
		Stdout: &bellSkipper{},
	}

	selectedIndex, _, err := prompt.Run()
	if err != nil {
		return err
	}

	obj := objects[selectedIndex]

	if obj.IsBack {
		if len(state.prefixes) > 0 {
			state.prefixes = state.prefixes[:len(state.prefixes)-1]
		} else {
			state.bucket = ""
		}
	} else if obj.IsDir {
		state.prefixes = append(state.prefixes, path.Base(obj.Name))
	} else {
		state.key = obj.Name
	}

	return nil
}

func setAction(state *state) (int, error) {

	prompt := promptui.Select{
		Label: "Select Action",
		Items: actionList,
		Searcher: func(input string, index int) bool {
			return strings.Contains(actionList[index].name, input)
		},
		StartInSearchMode: true,
	}

	i, _, err := prompt.Run()
	if err != nil {
		return -1, err
	}

	return i, nil
}
