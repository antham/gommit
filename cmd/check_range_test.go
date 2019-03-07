package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	"github.com/antham/gommit/gommit"
)

func TestCheckRangeWithErrors(t *testing.T) {
	err := exec.Command("../features/repo.sh").Run()

	if err != nil {
		logrus.Fatal(err)
	}

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
		{
			"check",
			"range",
		},
		{
			"check",
			"range",
			"master~2",
		},
		{
			"check",
			"range",
			"master~1",
			"master~2",
			"test",
			"whatever",
		},
		{
			"check",
			"range",
			"master~1",
			"master~2",
			"whatever",
		},
		{
			"check",
			"range",
			"master~1",
			"master~2",
			"check.go",
		},
		{
			"check",
			"range",
			"whatever",
			"master",
			"testing-repository/",
		},
		{
			"check",
			"range",
			"master~2",
			"master~1",
			"testing-repository/",
		},
		{
			"check",
			"range",
			"master~2",
			"master~1",
			"testing-repository/",
		},
	}

	errors := []error{
		fmt.Errorf("Two arguments required : origin commit and end commit"),
		fmt.Errorf("Two arguments required : origin commit and end commit"),
		fmt.Errorf("3 arguments must be provided at most"),
		fmt.Errorf(`Ensure "whatever" directory exists`),
		fmt.Errorf(`"check.go" must be a directory`),
		fmt.Errorf(`Reference "whatever" can't be found in git repository`),
		fmt.Errorf(`At least one matcher must be defined`),
		fmt.Errorf(`At least one example must be defined`),
	}

	configs := []string{
		path + "/../features/.gommit.toml",
		path + "/../features/.gommit.toml",
		path + "/../features/.gommit.toml",
		path + "/../features/.gommit.toml",
		path + "/../features/.gommit.toml",
		path + "/../features/.gommit.toml",
		path + "/../features/.gommit-no-matchers.toml",
		path + "/../features/.gommit-no-examples.toml",
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

			os.Args = []string{"", "--config", configs[i]}
			os.Args = append(os.Args, a...)

			_ = RootCmd.Execute()
		}()

		w.Wait()

		assert.Error(t, errc, "Must return an error")
		assert.EqualError(t, errc, errors[i].Error(), "Must return an error : "+errors[i].Error())
	}
}

func TestCheckRangeWithBadCommitMessage(t *testing.T) {
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

	var matchings *[]*gommit.Matching
	var examples map[string]string

	renderMatchings = func(m *[]*gommit.Matching) {
		matchings = m
	}

	renderExamples = func(e map[string]string) {
		examples = e
	}

	w.Add(1)

	go func() {
		defer func() {
			if r := recover(); r != nil && r.(int) == 0 {
				matchings = &[]*gommit.Matching{}
				examples = map[string]string{}
			}

			w.Done()
		}()

		os.Args = []string{"", "--config", path + "/../features/.gommit.toml", "check", "range", "test~3", "test", path + "/testing-repository"}

		Execute()
	}()

	w.Wait()

	assert.Len(t, *matchings, 1, "Must return 1 commits")
	assert.Len(t, examples, 3, "Must return 3 examples")
}

func TestCheckRangeWithNoErrors(t *testing.T) {
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

		os.Args = []string{"", "--config", path + "/../features/.gommit.toml", "check", "range", "test~2", "test", path + "/testing-repository"}

		Execute()
	}()

	w.Wait()

	assert.EqualValues(t, 0, code, "Must exit without errors (exit 0)")
	assert.EqualValues(t, "Everything is ok", message, "Must return a message to inform everything is ok")
}
