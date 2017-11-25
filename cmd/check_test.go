package cmd

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/antham/gommit/gommit"
)

func TestBuildOptionsWithDefaultValues(t *testing.T) {
	opts := buildOptions()

	assert.Equal(t, gommit.Options{SummaryLength: 50, CheckSummaryLength: false, ExcludeMergeCommits: false}, opts)
}
