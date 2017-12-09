package cmd

import (
	"os"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	var code int
	var message string
	var w sync.WaitGroup

	info = func(msg string) {
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

		os.Args = []string{"", "version"}

		Execute()
	}()

	w.Wait()

	assert.EqualValues(t, 0, code, "Must exit without errors (exit 0)")
	assert.EqualValues(t, "v2.0.0", message, "Must return app version")
}
