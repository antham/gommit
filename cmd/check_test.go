package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/antham/gommit/gommit"
)

func TestBuildOptionsWithDefaultValues(t *testing.T) {
	opts := buildOptions()

	assert.Equal(t, gommit.Options{SummaryLength: 50, CheckSummaryLength: false, ExcludeMergeCommits: false}, opts)
}

func TestParseDirectoryWithErrors(t *testing.T) {
	_, err := parseDirectory("/test")

	assert.EqualError(t, err, `Ensure "/test" directory exists`)

	_, err = os.Create("/tmp/file")

	assert.NoError(t, err)

	_, err = parseDirectory("/tmp/file")

	assert.EqualError(t, err, `"/tmp/file" must be a directory`)
}

func TestParseDirectorydWithErrors(t *testing.T) {
	path, err := parseDirectory("")

	assert.NoError(t, err)
	assert.Contains(t, path, "antham/gommit")
}
