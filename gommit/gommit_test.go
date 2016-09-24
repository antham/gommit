package gommit

import (
	"os/exec"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestFetchCommitsWithValidInterval(t *testing.T) {
	err := exec.Command("./repo.sh").Run()

	if err != nil {
		logrus.Fatal(err)
	}

	commits, err := FetchCommits("test", "master~2", "master")
	assert.NoError(t, err, "Must return no errors")

	expected := []string{"update(file1) : update file 1", "feat(file2) : new file 2"}

	assert.Len(t, *commits, 2, "Must contains 2 commits")

	for i, c := range *commits {
		assert.Equal(t, expected[i], c.Summary(), "Wrong commit fetched from repository")
	}
}

func TestFetchCommitsWithWrongRepository(t *testing.T) {
	_, err := FetchCommits("testtesttest", "4906f72818c0185162a3ec9c39a711d7c2842d40", "master")
	assert.Error(t, err, "Must return an error")
}

func TestFetchCommitsWithWrongInterval(t *testing.T) {
	_, err := FetchCommits("test", "4906f72818c0185162a3ec9c39a711d7c2842d40", "maste")
	assert.Error(t, err, "Must return an error")
}
