package main

import (
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/manifoldco/promptui"
)

type appState struct {
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
				fmt.Println("Error:", err)
			}
		}
	}()

	myFigure := figure.NewFigure("S3 Copy Tool", "", true)
	myFigure.Print()

	fmt.Println("")

	state := appState{
		bucket:   "",
		prefixes: []string{},
		key:      "",
		action:   -1,
	}

	if len(os.Args) >= 2 {

		p, err := ParseS3Url(os.Args[1])
		if err != nil {
			panic(err)
		}

		state = appState{
			bucket:   p.Bucket,
			prefixes: p.Prefixes,
			key:      p.Key,
			action:   -1,
		}
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

func setBucketState(state *appState) error {

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
		Templates: &promptui.SelectTemplates{
			Active:   fmt.Sprintf("%s  {{ . | underline }}", promptui.IconSelect),
			Selected: fmt.Sprintf(`{{ "%s" | green }} {{ "s3://" | faint }}{{. | faint}}`, promptui.IconGood),
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

func setObjectState(state *appState) error {

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
			Inactive: " {{if .IsDir}} {{ .Name }}/ {{- else}} {{ .Name }} {{- end}}",
			Active:   fmt.Sprintf("%s  {{if .IsDir -}} {{ .Name | underline }}/ {{- else -}} {{ .Name | underline }} {{- end}}", promptui.IconSelect),
			Selected: fmt.Sprintf(`{{ "%s" | green }} {{ "s3://%s/" | faint }}{{.Value | faint}}`, promptui.IconGood, state.bucket),
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
		state.prefixes = append(state.prefixes, path.Base(obj.Value))
	} else {
		state.key = obj.Value
	}

	return nil
}

func setAction(state *appState) (int, error) {

	prompt := promptui.Select{
		Label: "Select Action",
		Items: actionList,
		Searcher: func(input string, index int) bool {
			return strings.Contains(actionList[index].name, input)
		},
		Templates: &promptui.SelectTemplates{
			Active: fmt.Sprintf("%s  {{ . | underline }}", promptui.IconSelect),
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
