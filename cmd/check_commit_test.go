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

func TestCheckCommitWithErrors(t *testing.T) {
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
		[]string{
			"check",
			"commit",
		},
		[]string{
			"check",
			"commit",
			"whatever",
		},
		[]string{
			"check",
			"commit",
			"826f193edd4ba9d6d1799b66fa64f9a84f1db3bf",
			"whatever",
		},
		[]string{
			"check",
			"commit",
			"826f193edd4ba9d6d1799b66fa64f9a84f1db3bf",
			"test",
		},
		[]string{
			"check",
			"commit",
			"826f193edd4ba9d6d1799b66fa64f9a84f1db3bf",
			"test",
		},
		[]string{
			"check",
			"commit",
			"826f193edd4ba9d6d1799b66fa64f9a84f1db3bf",
			"test",
		},
		[]string{
			"check",
			"commit",
			"826f193edd4ba9d6d1799b66fa64f9a84f1db3bf",
			"test",
		},
	}

	errors := []error{
		fmt.Errorf("One argument required : commit id"),
		fmt.Errorf(`Argument must be a valid commit id`),
		fmt.Errorf(`Ensure "whatever" directory exists`),
		fmt.Errorf(`object not found`),
		fmt.Errorf(`At least one matcher must be defined`),
		fmt.Errorf(`At least one example must be defined`),
		fmt.Errorf(`Regexp "**" identified by "all" is not a valid regexp, please check the syntax`),
	}

	configs := []string{
		path + "/../features/.gommit.toml",
		path + "/../features/.gommit.toml",
		path + "/../features/.gommit.toml",
		path + "/../features/.gommit.toml",
		path + "/../features/.gommit-no-matchers.toml",
		path + "/../features/.gommit-no-examples.toml",
		path + "/../features/.gommit-wrong-regexp.toml",
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

func TestCheckCommit(t *testing.T) {
	path, err := os.Getwd()

	if err != nil {
		logrus.Fatal(err)
	}

	for _, filename := range []string{"../features/repo.sh"} {
		err = exec.Command(filename).Run()

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

	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = "test"

	ID, err := cmd.Output()

	if err != nil {
		logrus.Fatal(err)
	}

	w.Add(1)

	go func() {
		defer func() {
			if r := recover(); r != nil {
				code = r.(int)
			}

			w.Done()
		}()

		os.Args = []string{"", "--config", path + "/../features/.gommit.toml", "check", "commit", string(ID[:len(ID)-1]), path + "/test"}

		Execute()
	}()

	w.Wait()

	assert.EqualValues(t, 0, code, "Must exit with no errors (exit 0)")
	assert.EqualValues(t, "Everything is ok", message, "Must return a message to inform everything is ok")
}

func TestCheckCommitWithBadMessage(t *testing.T) {
	path, err := os.Getwd()

	if err != nil {
		logrus.Fatal(err)
	}

	for _, filename := range []string{"../features/repo.sh", "../features/bad-commit.sh"} {
		err = exec.Command(filename).Run()

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

	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = "test"

	ID, err := cmd.Output()

	if err != nil {
		logrus.Fatal(err)
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

		os.Args = []string{"", "--config", path + "/../features/.gommit.toml", "check", "commit", string(ID[:len(ID)-1]), path + "/test"}

		Execute()
	}()

	w.Wait()

	assert.EqualValues(t, 1, code, "Must exit with errors (exit 1)")
	assert.Empty(t, message, "Must return no success message")
	assert.Len(t, *matchings, 1, "Must return 1 commits")
	assert.Len(t, examples, 3, "Must return 3 examples")
}
