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

	expected := []string{
		"update(file1) : update file 1\n\nupdate file 1 with a text\n",
		"feat(file2) : new file 2\n\ncreate a new file 2\n",
	}

	assert.Len(t, *commits, 2, "Must contains 2 commits")

	for i, c := range *commits {
		assert.Equal(t, expected[i], c.Message, "Wrong commit fetched from repository")
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

	match := messageMatchTemplate(msg, temp)
	assert.True(t, match, "Message must match template")
}

func TestMessageMatchTemplate2(t *testing.T) {
	msg := "(feat) : Hello world !\n"
	msg += "* test1\n"
	msg += "* test2\n"
	msg += "* test3\n"

	temp := "\\((?:feat|test|bug)\\) : .*?\n(?:\\* .*?\n)+"

	match := messageMatchTemplate(msg, temp)
	assert.True(t, match, "Message must match template")
}

func TestDontMessageMatchTemplate(t *testing.T) {
	msg := "This is a test\n"
	msg += "=> an added reason\n"

	temp := "This is a test\n=> an added reaso\n"

	match := messageMatchTemplate(msg, temp)
	assert.False(t, match, "Message must not match template")
}

func TestRunMatching(t *testing.T) {
	err := exec.Command("../features/repo.sh").Run()

	if err != nil {
		logrus.Fatal(err)
	}

	q := Query{
		"test/",
		"master~2",
		"master",
		map[string]string{"simple": "(?:update|feat)\\(.*?\\) : .*?\\n\\n.*?\\n"},
		map[string]bool{
			"check-summary-length":  false,
			"exclude-merge-commits": false,
		},
	}

	m, err := RunMatching(q)

	assert.NoError(t, err, "Must return no errors")
	assert.Len(t, *m, 0, "Must return no items, match was successful for every commit")
}

func TestRunMatchingWithAMessageErrorCommit(t *testing.T) {
	err := exec.Command("../features/repo.sh").Run()

	if err != nil {
		logrus.Fatal(err)
	}

	q := Query{
		"test/",
		"master~2",
		"master",
		map[string]string{"simple": "(?:update)\\(.*?\\) : .*?\\n\\n.*?\\n"},
		map[string]bool{
			"check-summary-length":  false,
			"exclude-merge-commits": false,
		},
	}

	m, err := RunMatching(q)

	assert.NoError(t, err, "Must return no errors")
	assert.Len(t, *m, 1, "Must return one item")
	assert.Equal(t, "feat(file2) : new file 2\n\ncreate a new file 2\n", (*m)[0].Message, "Must contains commit message")
	assert.EqualError(t, (*m)[0].MessageError, "No template match commit message", "Must contains commit message error")
	assert.NoError(t, (*m)[0].SummaryError, "Must not contains error")
}

func TestRunMatchingWithASummaryErrorCommit(t *testing.T) {
	for _, filename := range []string{"../features/repo.sh", "../features/bad-summary-message-commit.sh"} {
		err := exec.Command(filename).Run()

		if err != nil {
			logrus.Fatal(err)
		}
	}

	q := Query{
		"test/",
		"master~1",
		"master",
		map[string]string{"simple": ".*\n"},
		map[string]bool{
			"check-summary-length":  true,
			"exclude-merge-commits": false,
		},
	}

	m, err := RunMatching(q)

	assert.NoError(t, err, "Must return no errors")
	assert.Len(t, *m, 1, "Must return one item")
	assert.Equal(t, "A very long summary commit greater than minimum length 50\n", (*m)[0].Message, "Must contains commit message")
	assert.NoError(t, (*m)[0].MessageError, "Must not contains error")
	assert.EqualError(t, (*m)[0].SummaryError, "Commit summary length is greater than 50 characters", "Must contains summary message error")
}

func TestRunMatchingWithAMessageErrorCommitWithoutMergeCommist(t *testing.T) {
	for _, filename := range []string{"../features/repo.sh", "../features/merge-commit.sh"} {
		err := exec.Command(filename).Run()

		if err != nil {
			logrus.Fatal(err)
		}
	}

	q := Query{
		"test/",
		"master~2",
		"master",
		map[string]string{"simple": "(?:update)\\(.*?\\) : .*?\\n\\n.*?\\n"},
		map[string]bool{
			"check-summary-length":  false,
			"exclude-merge-commits": true,
		},
	}

	m, err := RunMatching(q)

	assert.NoError(t, err, "Must return no errors")
	assert.Len(t, *m, 1, "Must return one item")
	assert.Equal(t, "feat(file) : new file 3\n\ncreate a new file 3\n", (*m)[0].Message, "Must contains commit message")
	assert.EqualError(t, (*m)[0].MessageError, "No template match commit message", "Must contains commit message error")
	assert.NoError(t, (*m)[0].SummaryError, "Must not contains error")
}

func TestRunMatchingWithAMessageErrorCommitWithMergeCommits(t *testing.T) {
	for _, filename := range []string{"../features/repo.sh", "../features/merge-commit.sh"} {
		err := exec.Command(filename).Run()

		if err != nil {
			logrus.Fatal(err)
		}
	}

	q := Query{
		"test/",
		"master~2",
		"master",
		map[string]string{"simple": "(?:update)\\(.*?\\) : .*?\\n\\n.*?\\n"},
		map[string]bool{
			"check-summary-length":  false,
			"exclude-merge-commits": false,
		},
	}

	m, err := RunMatching(q)

	assert.NoError(t, err, "Must return no errors")
	assert.Len(t, *m, 2, "Must return two itesm")
	assert.Equal(t, "Merge branch 'test'\n", (*m)[0].Message, "Must contains commit message")
	assert.EqualError(t, (*m)[0].MessageError, "No template match commit message", "Must contains commit message error")
	assert.NoError(t, (*m)[0].SummaryError, "Must not contains error")
}

func TestRunMatchingWithAnInvalidCommitRange(t *testing.T) {
	err := exec.Command("../features/repo.sh").Run()

	if err != nil {
		logrus.Fatal(err)
	}

	q := Query{
		"test/",
		"master",
		"master~2",
		map[string]string{"simple": "(?:update)\\(.*?\\) : .*?\\n\\n.*?\\n"},
		map[string]bool{
			"check-summary-length":  false,
			"exclude-merge-commits": false,
		},
	}

	m, err := RunMatching(q)

	assert.Error(t, err, "Must return an error")
	assert.EqualError(t, err, `No commits found between "master" and "master~2"`, "Must return an explicit message error")
	assert.Len(t, *m, 0, "Must return no item")
}

func TestRunMatchingWithAnUnexistingCommitRange(t *testing.T) {
	err := exec.Command("../features/repo.sh").Run()

	if err != nil {
		logrus.Fatal(err)
	}

	q := Query{
		"test/",
		"master~15",
		"master",
		map[string]string{"simple": "(?:update)\\(.*?\\) : .*?\\n\\n.*?\\n"},
		map[string]bool{
			"check-summary-length":  false,
			"exclude-merge-commits": false,
		},
	}

	m, err := RunMatching(q)

	assert.EqualError(t, err, "Can't find reference", "Must return an error")
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
