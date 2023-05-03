package main

import (
	"fmt"
	"path"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/manifoldco/promptui"
)

type ListObj struct {
	Name   string
	IsDir  bool
	IsBack bool
}

type state struct {
	bucket   string
	prefixes []string
	key      string
	action   int
}

func main() {

	defer func() {
		if err := recover(); err != nil {
			castErr, isErr := err.(error)
			if isErr && castErr.Error() == "^C" {
				fmt.Println("Exiting")

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

	buckets, err := getBucketsListAWS()
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

	objects, err := getObjectsListAWS(state.bucket, joinPrefixes(state.prefixes))
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
			Inactive: "{{if .IsDir}} {{ .Name | bold }} {{- else}} {{ .Name }} {{- end}}",
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
		Stdout:            &bellSkipper{},
	}

	i, _, err := prompt.Run()
	if err != nil {
		return -1, err
	}

	return i, nil
}
