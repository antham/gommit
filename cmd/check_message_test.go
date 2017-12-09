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

func TestCheckMessageWithErrors(t *testing.T) {
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
			"message",
		},
		{
			"check",
			"message",
			"test",
		},
		{
			"check",
			"message",
			"test",
		},
	}

	errors := []error{
		fmt.Errorf("One argument required : message"),
		fmt.Errorf(`At least one matcher must be defined`),
		fmt.Errorf(`At least one example must be defined`),
	}

	configs := []string{
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

func TestCheckMessage(t *testing.T) {
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

		os.Args = []string{"", "--config", path + "/../features/.gommit.toml", "check", "message", "feat(cmd) : everything is fine\n"}

		Execute()
	}()

	w.Wait()

	assert.EqualValues(t, 0, code, "Must exit without errors (exit 0)")
	assert.EqualValues(t, "Everything is ok", message, "Must return a message to inform everything is ok")
}

func TestCheckMessageWithBadMessage(t *testing.T) {
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
			if r := recover(); r != nil {
				code = r.(int)
			}

			w.Done()
		}()

		os.Args = []string{"", "--config", path + "/../features/.gommit.toml", "check", "message", "everything is fine\n"}

		Execute()
	}()

	w.Wait()

	assert.EqualValues(t, 1, code, "Must exit with errors (exit 1)")
	assert.Empty(t, message, "Must return no success message")
	assert.Len(t, *matchings, 1, "Must return 1 commits")
	assert.Len(t, examples, 3, "Must return 3 examples")
}
