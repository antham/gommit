package reference

import (
	"os"
	"os/exec"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-git.v4"
)

// Result correctness is checked against git log

func fetchCommitFromAGivenInterval(from string, to string) ([]*git.Commit, error) {
	for _, filename := range []string{"../features/repo.sh", "../features/merge-commits.sh"} {
		err := exec.Command(filename).Run()

		if err != nil {
			logrus.Fatal(err)
		}
	}

	path, err := os.Getwd()

	if err != nil {
		logrus.Fatal(err)
	}

	repo, err := git.NewFilesystemRepository(path + "/test/.git/")

	if err != nil {
		logrus.Fatal(err)
	}

	commits, err := FetchCommitInterval(repo, from, to)

	return *commits, err
}

func TestFetchCommitInterval(t *testing.T) {

	commits, err := fetchCommitFromAGivenInterval("master~5", "master")

	assert.NoError(t, err, "Must return no error")

	expected := []string{
		"feat(file10) : new file 10\n\ncreate a new file 10\n",
		"feat(file9) : new file 9\n\ncreate a new file 9\n",
		"Merge branch 'test'\n",
		"Merge branch 'test1' into test\n",
		"Merge branch 'test2' into test1\n",
		"feat(file8) : new file 8\n\ncreate a new file 8\n",
		"feat(file7) : new file 7\n\ncreate a new file 7\n",
		"feat(file6) : new file 6\n\ncreate a new file 6\n",
		"feat(file5) : new file 5\n\ncreate a new file 5\n",
		"feat(file4) : new file 4\n\ncreate a new file 4\n",
		"feat(file3) : new file 3\n\ncreate a new file 3\n",
		"update(file1) : update file 1\n\nupdate file 1 with a text\n",
		"feat(file2) : new file 2\n\ncreate a new file 2\n",
	}

	results := []string{}

	for _, c := range commits {
		results = append(results, c.Message)
	}

	assert.Equal(t, expected, results, "Must return all commit hsitory")
}

func TestFetchCommitIntervalWithSubtrees1(t *testing.T) {
	commits, err := fetchCommitFromAGivenInterval("master~2^2^2^2~", "master")

	assert.NoError(t, err, "Must return no error")

	expected := []string{
		"feat(file10) : new file 10\n\ncreate a new file 10\n",
		"feat(file9) : new file 9\n\ncreate a new file 9\n",
		"Merge branch 'test'\n",
		"Merge branch 'test1' into test\n",
		"Merge branch 'test2' into test1\n",
		"feat(file8) : new file 8\n\ncreate a new file 8\n",
	}

	results := []string{}

	for _, c := range commits {
		results = append(results, c.Message)
	}

	assert.Equal(t, expected, results, "Must return a commit history subtree")
}

func TestFetchCommitIntervalWithSubtrees2(t *testing.T) {
	commits, err := fetchCommitFromAGivenInterval("master~2^2~", "master")

	assert.NoError(t, err, "Must return no error")

	expected := []string{
		"feat(file10) : new file 10\n\ncreate a new file 10\n",
		"feat(file9) : new file 9\n\ncreate a new file 9\n",
		"Merge branch 'test'\n",
		"Merge branch 'test1' into test\n",
		"Merge branch 'test2' into test1\n",
		"feat(file8) : new file 8\n\ncreate a new file 8\n",
		"feat(file7) : new file 7\n\ncreate a new file 7\n",
		"feat(file6) : new file 6\n\ncreate a new file 6\n",
		"feat(file5) : new file 5\n\ncreate a new file 5\n",
	}

	results := []string{}

	for _, c := range commits {
		results = append(results, c.Message)
	}

	assert.Equal(t, expected, results, "Must return a commit history subtree")
}

func TestFetchCommitIntervalWithSubtrees3(t *testing.T) {
	commits, err := fetchCommitFromAGivenInterval("master~", "master")

	assert.NoError(t, err, "Must return no error")

	expected := []string{
		"feat(file10) : new file 10\n\ncreate a new file 10\n",
	}

	results := []string{}

	for _, c := range commits {
		results = append(results, c.Message)
	}

	assert.Equal(t, expected, results, "Must return a commit history subtree")
}

func TestFetchCommitIntervalWithAnArbitrayRange1(t *testing.T) {
	commits, err := fetchCommitFromAGivenInterval("master~5", "master~3")

	assert.NoError(t, err, "Must return no error")

	expected := []string{
		"update(file1) : update file 1\n\nupdate file 1 with a text\n",
		"feat(file2) : new file 2\n\ncreate a new file 2\n",
	}

	results := []string{}

	for _, c := range commits {
		results = append(results, c.Message)
	}

	assert.Equal(t, expected, results, "Must return a commit history subtree")
}

func TestFetchCommitIntervalWithAnArbitrayRange2(t *testing.T) {
	commits, err := fetchCommitFromAGivenInterval("master~4", "master~5")

	assert.NoError(t, err, "Must return no error")

	expected := []string{}

	results := []string{}

	for _, c := range commits {
		results = append(results, c.Message)
	}

	assert.Equal(t, expected, results, "Must return no commit history")
}

func TestFetchCommitIntervalWithAnArbitrayRange3(t *testing.T) {
	commits, err := fetchCommitFromAGivenInterval("master~5", "master~3")

	assert.NoError(t, err, "Must return no error")

	expected := []string{
		"update(file1) : update file 1\n\nupdate file 1 with a text\n",
		"feat(file2) : new file 2\n\ncreate a new file 2\n",
	}

	results := []string{}

	for _, c := range commits {
		results = append(results, c.Message)
	}

	assert.Equal(t, expected, results, "Must return a commit history subtree")
}

func TestFetchCommitIntervalWithAnArbitrayRange4(t *testing.T) {
	commits, err := fetchCommitFromAGivenInterval("master~2^2^", "master~2")

	assert.NoError(t, err, "Must return no error")

	expected := []string{
		"Merge branch 'test'\n",
		"Merge branch 'test1' into test\n",
		"Merge branch 'test2' into test1\n",
		"feat(file8) : new file 8\n\ncreate a new file 8\n",
		"feat(file7) : new file 7\n\ncreate a new file 7\n",
		"feat(file6) : new file 6\n\ncreate a new file 6\n",
		"feat(file5) : new file 5\n\ncreate a new file 5\n",
	}

	results := []string{}

	for _, c := range commits {
		results = append(results, c.Message)
	}

	assert.Equal(t, expected, results, "Must return a commit history subtree")
}

func TestFetchCommitIntervalWithUnexistingRange(t *testing.T) {
	commits, err := fetchCommitFromAGivenInterval("master~25", "master~30")

	assert.EqualError(t, err, "Can't find reference", "Must return an error, interval doesn't exist")
	assert.Equal(t, []*git.Commit{}, commits, "Must contains no datas")
}
