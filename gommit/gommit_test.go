package gommit

import (
	"os/exec"
	"testing"

	"github.com/Sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestFetchCommitsWithValidInterval(t *testing.T) {
	err := exec.Command("../features/repo.sh").Run()

	if err != nil {
		logrus.Fatal(err)
	}

	commits, err := fetchCommits("test", "master~2", "master")
	assert.NoError(t, err, "Must return no errors")

	expected := []string{"update(file1) : update file 1", "feat(file2) : new file 2"}

	assert.Len(t, *commits, 2, "Must contains 2 commits")

	for i, c := range *commits {
		assert.Equal(t, expected[i], c.Summary(), "Wrong commit fetched from repository")
	}
}

func TestFetchCommitsWithWrongRepository(t *testing.T) {
	_, err := fetchCommits("testtesttest", "4906f72818c0185162a3ec9c39a711d7c2842d40", "master")
	assert.Error(t, err, "Must return an error")
}

func TestFetchCommitsWithWrongInterval(t *testing.T) {
	_, err := fetchCommits("test", "4906f72818c0185162a3ec9c39a711d7c2842d40", "maste")
	assert.Error(t, err, "Must return an error")
}

func TestMessageMatchTemplate1(t *testing.T) {
	msg := "(feat) : Hello world !"
	temp := "\\((?:feat|test|bug)\\) : .*"

	match, extractedGroup := messageMatchTemplate(msg, temp)
	assert.True(t, match, "Message must match template")
	assert.Equal(t, msg, extractedGroup, "Must return extracted group")
}

func TestMessageMatchTemplate2(t *testing.T) {
	msg := "(feat) : Hello world !\n"
	msg += "* test1\n"
	msg += "* test2\n"
	msg += "* test3\n"

	temp := "\\((?:feat|test|bug)\\) : .*?\n(?:\\* .*?\n)+"

	match, extractedGroup := messageMatchTemplate(msg, temp)
	assert.True(t, match, "Message must match template")
	assert.Equal(t, msg, extractedGroup, "Must return extracted group")
}

func TestDontMessageMatchTemplate(t *testing.T) {
	msg := "This is a test\n"
	msg += "=> an added reason\n"

	temp := "This is a test\n=> an added reaso\n"

	match, extractedGroup := messageMatchTemplate(msg, temp)
	assert.False(t, match, "Message must not match template")
	assert.NotEqual(t, msg, extractedGroup, "Must return extracted group")
}

func TestRunMatching(t *testing.T) {
	m, err := RunMatching("test/", "master~2", "master", map[string]string{"simple": "(?:update|feat)\\(.*?\\) : .*?\\n\\n.*?\\n"})

	assert.NoError(t, err, "Must return no errors")
	assert.Len(t, *m, 0, "Must return no items, match was successful for every commit")
}

func TestRunMatchingWithAnErrorCommit(t *testing.T) {
	m, err := RunMatching("test/", "master~2", "master", map[string]string{"simple": "(?:update)\\(.*?\\) : .*?\\n\\n.*?\\n"})

	assert.NoError(t, err, "Must return no errors")
	assert.Len(t, *m, 1, "Must return one item")
	assert.Equal(t, (*m)[0]["message"], "feat(file2) : new file 2\n\ncreate a new file 2\n", "Must contains commit message")
}

func TestRunMatchingWithAnInvalidCommitRange(t *testing.T) {
	m, err := RunMatching("test/", "master", "master~2", map[string]string{"simple": "(?:update)\\(.*?\\) : .*?\\n\\n.*?\\n"})

	assert.Error(t, err, "Must return an error")
	assert.EqualError(t, err, `No commits found between "master" and "master~2"`, "Must return an explicit message error")
	assert.Len(t, *m, 0, "Must return no item")
}

func TestRunMatchingWithAnUnexistingCommitRange(t *testing.T) {
	m, err := RunMatching("test/", "master~15", "master", map[string]string{"simple": "(?:update)\\(.*?\\) : .*?\\n\\n.*?\\n"})

	assert.Error(t, err, "Must return an error")
	assert.EqualError(t, err, `Interval between "master~15" and "master" can't be fetched`, "Must return an explicit message error")
	assert.Len(t, *m, 0, "Must return no item")
}

func TestIsValidSummaryLengthWithCorrectSize(t *testing.T) {
	assert.True(t, isValidSummaryLength("test"), "Must have a length lower than 50 characters")
}

func TestIsValidSummaryLengthWithInCorrectSize(t *testing.T) {
	assert.False(t, isValidSummaryLength("ttttttttttttttttttttttttttttttttttttttttttttttttttt"), "Must have a length lower than 50 characters")
}

func TestIsMergeCommitWithANonMergeCommit(t *testing.T) {
	err := exec.Command("../features/repo.sh").Run()

	if err != nil {
		logrus.Fatal(err)
	}

	commits, err := fetchCommits("test", "master~2", "master")

	if err != nil {
		logrus.Fatal(err)
	}

	assert.False(t, isMergeCommit((*commits)[0]), "Must return false with non merge commit")
}

func TestIsMergeCommitWithAMergeCommit(t *testing.T) {
	for _, filename := range []string{"../features/repo.sh", "../features/merge-commit.sh"} {
		err := exec.Command(filename).Run()

		if err != nil {
			logrus.Fatal(err)
		}
	}

	commits, err := fetchCommits("test", "master~2", "master")

	if err != nil {
		logrus.Fatal(err)
	}

	assert.True(t, isMergeCommit((*commits)[0]), "Must return false with non merge commit")
}
