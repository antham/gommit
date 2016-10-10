package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/antham/gommit/gommit"
)

func TestCheckWithErrors(t *testing.T) {
	path, err := os.Getwd()

	if err != nil {
		logrus.Fatal(err)
	}

	exitError = func() {
		panic(1)
	}

	exitSuccess = func() {
		panic(0)
	}

	var errc error

	failure = func(err error) {
		errc = err
	}

	arguments := [][]string{
		[]string{
			"check",
		},
		[]string{
			"check",
			"master~2",
		},
		[]string{
			"check",
			"master~1",
			"master~2",
			"test",
			"whatever",
		},
		[]string{
			"check",
			"master~1",
			"master~2",
			"whatever",
		},
		[]string{
			"check",
			"master~1",
			"master~2",
			"check.go",
		},
		[]string{
			"check",
			"master~15",
			"master",
		},
	}

	errors := []error{
		fmt.Errorf("Two arguments required : origin commit and end commit"),
		fmt.Errorf("Two arguments required : origin commit and end commit"),
		fmt.Errorf("3 arguments must be provided at most"),
		fmt.Errorf(`Ensure "whatever" directory exists`),
		fmt.Errorf(`"check.go" must be a directory`),
		fmt.Errorf(`Interval between "master~15" and "master" can't be fetched`),
	}

	for i, a := range arguments {
		var w sync.WaitGroup

		w.Add(1)

		go func() {
			defer func() {
				if r := recover(); r != nil && r.(int) == 0 {
					errc = nil
				}

				w.Done()
			}()

			os.Args = []string{"", "--config", path + "/../.gommit.toml"}
			os.Args = append(os.Args, a...)
			_ = RootCmd.Execute()
		}()

		w.Wait()

		assert.Error(t, errc, "Must return an error")
		assert.EqualError(t, errc, errors[i].Error(), "Must return an error : "+errors[i].Error())
	}
}

func TestCheckCommitsWithBadCommitMessage(t *testing.T) {
	path, err := os.Getwd()

	if err != nil {
		logrus.Fatal(err)
	}

	for _, filename := range []string{"../features/repo.sh", "../features/bad-commit.sh"} {

		err := exec.Command(filename).Run()

		if err != nil {
			logrus.Fatal(err)
		}
	}

	exitError = func() {
		panic(1)
	}

	exitSuccess = func() {
		panic(0)
	}

	var w sync.WaitGroup

	var errors *[]gommit.CommitError
	var examples map[string]string

	renderErrors = func(e *[]gommit.CommitError) {
		errors = e
	}

	renderExamples = func(e map[string]string) {
		examples = e
	}

	w.Add(1)

	go func() {
		defer func() {
			if r := recover(); r != nil && r.(int) == 0 {
				e := []gommit.CommitError{}
				errors = &e
				examples = map[string]string{}
			}

			w.Done()
		}()

		os.Args = []string{"", "--config", path + "/../features/.gommit.toml", "check", "master~3", "master", path + "/test"}

		Execute()
	}()

	w.Wait()

	assert.Len(t, *errors, 1, "Must return 1 commits")
	assert.Len(t, examples, 3, "Must return 3 examples")
}

func TestCheckCommitsWithNoErrors(t *testing.T) {
	path, err := os.Getwd()

	if err != nil {
		logrus.Fatal(err)
	}

	for _, filename := range []string{"../features/repo.sh"} {

		err := exec.Command(filename).Run()

		if err != nil {
			logrus.Fatal(err)
		}
	}

	var code int
	var message string
	var w sync.WaitGroup

	success = func(msg string) {
		message = msg
	}

	exitError = func() {
		panic(1)
	}

	exitSuccess = func() {
		panic(0)
	}

	w.Add(1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				code = r.(int)
			}

			w.Done()
		}()

		os.Args = []string{"", "--config", path + "/../features/.gommit.toml", "check", "master~2", "master", path + "/test"}

		Execute()
	}()

	w.Wait()

	assert.EqualValues(t, 0, code, "Must exit without errors (exit 0)")
	assert.EqualValues(t, "Everything is ok", message, "Must return a message to inform everything is ok")
}
