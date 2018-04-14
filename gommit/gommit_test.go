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

	commits, err := fetchCommits("testing-repository", "test~2", "test")
	assert.NoError(t, err, "Must return no errors")

	expected := []string{
		"feat(file8) : new file 8\n\ncreate a new file 8\n",
		"feat(file7) : new file 7\n\ncreate a new file 7\n",
	}

	assert.Len(t, *commits, 2, "Must contains 2 commits")

	for i, c := range *commits {
		assert.Equal(t, expected[i], c.Message, "Wrong commit fetched from repository")
	}
}

func TestFetchCommitsWithWrongRepository(t *testing.T) {
	_, err := fetchCommits("testtesttest", "4906f72818c0185162a3ec9c39a711d7c2842d40", "test")
	assert.Error(t, err, "Must return an error")
}

func TestFetchCommitsWithWrongInterval(t *testing.T) {
	_, err := fetchCommits("testing-repository", "4906f72818c0185162a3ec9c39a711d7c2842d40", "maste")
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

func TestMessageDoesntMatchTemplate(t *testing.T) {
	msg := "This is a test\n"
	msg += "=> an added reason\n"

	temp := "This is a test\n=> an added reaso\n"

	match := messageMatchTemplate(msg, temp)
	assert.False(t, match, "Message must not match template")
}

func TestMatchRangeCommitQuery(t *testing.T) {
	err := exec.Command("../features/repo.sh").Run()

	if err != nil {
		logrus.Fatal(err)
	}

	q := RangeCommitQuery{
		"testing-repository/",
		"test~2",
		"test",
		map[string]string{"simple": "(?:update|feat)\\(.*?\\) : .*?\\n\\n.*?\\n"},
		Options{
			CheckSummaryLength:  false,
			ExcludeMergeCommits: false,
			SummaryLength:       50,
		},
	}

	m, err := MatchRangeCommitQuery(q)

	assert.NoError(t, err, "Must return no errors")
	assert.Len(t, *m, 0, "Must return no items, match was successful for every commit")
}

func TestMatchRangeCommitQueryrWithAMessageErrorCommit(t *testing.T) {
	err := exec.Command("../features/repo.sh").Run()

	if err != nil {
		logrus.Fatal(err)
	}

	q := RangeCommitQuery{
		"testing-repository/",
		"test~2",
		"test",
		map[string]string{"simple": "(?:update)\\(.*?\\) : .*?\\n\\n.*?\\n"},
		Options{
			CheckSummaryLength:  false,
			ExcludeMergeCommits: false,
			SummaryLength:       50,
		},
	}

	m, err := MatchRangeCommitQuery(q)

	assert.NoError(t, err, "Must return no errors")
	assert.Len(t, *m, 2, "Must return two items")
	assert.Equal(t, "feat(file8) : new file 8\n\ncreate a new file 8\n", (*m)[0].Context["message"], "Must contains commit message")
	assert.EqualError(t, (*m)[0].MessageError, "No template match commit message", "Must contains commit message error")
	assert.NoError(t, (*m)[0].SummaryError, "Must not contains error")
}

func TestMatchRangeCommitQueryASummaryErrorCommit(t *testing.T) {
	for _, filename := range []string{"../features/repo.sh", "../features/bad-summary-message-commit.sh"} {
		err := exec.Command(filename).Run()

		if err != nil {
			logrus.Fatal(err)
		}
	}

	q := RangeCommitQuery{
		"testing-repository/",
		"test~1",
		"test",
		map[string]string{"simple": ".*\n"},
		Options{
			CheckSummaryLength:  true,
			ExcludeMergeCommits: false,
			SummaryLength:       50,
		},
	}

	m, err := MatchRangeCommitQuery(q)

	assert.NoError(t, err, "Must return no errors")
	assert.Len(t, *m, 1, "Must return one item")
	assert.Equal(t, "A very long summary commit greater than minimum length 50\n", (*m)[0].Context["message"], "Must contains commit message")
	assert.NoError(t, (*m)[0].MessageError, "Must not contains error")
	assert.EqualError(t, (*m)[0].SummaryError, "Commit summary length is greater than 50 characters", "Must contains summary message error")
}

func TestMatchRangeCommitWithAMessageErrorCommitWithoutMergeCommit(t *testing.T) {
	for _, filename := range []string{"../features/repo.sh"} {
		err := exec.Command(filename).Run()

		if err != nil {
			logrus.Fatal(err)
		}
	}

	q := RangeCommitQuery{
		"testing-repository/",
		"test^^^^",
		"test",
		map[string]string{"simple": "(?:update)\\(.*?\\) : .*?\\n\\n.*?\\n"},
		Options{
			CheckSummaryLength:  false,
			ExcludeMergeCommits: true,
			SummaryLength:       50,
		},
	}

	m, err := MatchRangeCommitQuery(q)

	assert.NoError(t, err, "Must return no errors")
	assert.Len(t, *m, 7)
	for i := 0; i < 7; i++ {
		assert.NotContains(t, (*m)[i].Context["message"], "Merge")
	}
	assert.EqualError(t, (*m)[0].MessageError, "No template match commit message", "Must contains commit message error")
	assert.NoError(t, (*m)[0].SummaryError, "Must not contains error")
}

func TestMatchRangeCommitQueryWithAMessageErrorCommitWithMergeCommits(t *testing.T) {
	for _, filename := range []string{"../features/repo.sh"} {
		err := exec.Command(filename).Run()

		if err != nil {
			logrus.Fatal(err)
		}
	}

	q := RangeCommitQuery{
		"testing-repository/",
		"test^^^^",
		"test",
		map[string]string{"simple": "(?:update)\\(.*?\\) : .*?\\n\\n.*?\\n"},
		Options{
			CheckSummaryLength:  false,
			ExcludeMergeCommits: false,
			SummaryLength:       50,
		},
	}

	m, err := MatchRangeCommitQuery(q)

	assert.NoError(t, err, "Must return no errors")
	assert.Len(t, *m, 9, "Must return two itesm")
	assert.EqualError(t, (*m)[0].MessageError, "No template match commit message", "Must contains commit message error")
	assert.NoError(t, (*m)[0].SummaryError, "Must not contains error")
}

func TestMatchRangeCommitWithAnUnexistingCommitRange(t *testing.T) {
	err := exec.Command("../features/repo.sh").Run()

	if err != nil {
		logrus.Fatal(err)
	}

	q := RangeCommitQuery{
		"testing-repository/",
		"test~15",
		"test",
		map[string]string{"simple": "(?:update)\\(.*?\\) : .*?\\n\\n.*?\\n"},
		Options{
			CheckSummaryLength:  false,
			ExcludeMergeCommits: false,
			SummaryLength:       50,
		},
	}

	m, err := MatchRangeCommitQuery(q)

	assert.EqualError(t, err, "Reference \"test~15\" can't be found in git repository")
	assert.Len(t, *m, 0, "Must return no item")
}

func TestMatchMessageQuery(t *testing.T) {
	err := exec.Command("../features/repo.sh").Run()

	if err != nil {
		logrus.Fatal(err)
	}

	q := MessageQuery{
		"update(file) : fix",
		map[string]string{"simple": "(?:update|feat)\\(.*?\\) : .*"},
		Options{
			CheckSummaryLength:  false,
			ExcludeMergeCommits: false,
		},
	}

	m, err := MatchMessageQuery(q)

	assert.NoError(t, err, "Must return no error")
	assert.Equal(t, Matching{}, *m, "Must return an empty matching struct")
}

func TestMatchMessageQueryWithAMessageThatDoesntMatchTemplate(t *testing.T) {
	err := exec.Command("../features/repo.sh").Run()

	if err != nil {
		logrus.Fatal(err)
	}

	q := MessageQuery{
		"update(file) :",
		map[string]string{"simple": "(?:update|feat)\\(.*?\\) : .*"},
		Options{
			CheckSummaryLength:  false,
			ExcludeMergeCommits: false,
			SummaryLength:       50,
		},
	}

	m, err := MatchMessageQuery(q)

	assert.NoError(t, err, "Must return no error")
	assert.EqualError(t, m.MessageError, "No template match commit message", "Must return a template message error")
	assert.NoError(t, m.SummaryError, "Must return no summary error")
	assert.Equal(t, "update(file) :", m.Context["message"], "Must contains original message")
}

func TestMatchMessageQueryWithAMessageThatDoesntFitSummaryLength(t *testing.T) {
	err := exec.Command("../features/repo.sh").Run()

	if err != nil {
		logrus.Fatal(err)
	}

	q := MessageQuery{
		"update(file) : test test test test test test test test test test test test test test",
		map[string]string{"simple": "(?:update|feat)\\(.*?\\) : .*"},
		Options{
			CheckSummaryLength:  true,
			ExcludeMergeCommits: false,
			SummaryLength:       50,
		},
	}

	m, err := MatchMessageQuery(q)

	assert.NoError(t, err, "Must return no error")
	assert.NoError(t, m.MessageError, "Must return no template message error")
	assert.EqualError(t, m.SummaryError, "Commit summary length is greater than 50 characters", "Must return a template message error")
	assert.Equal(t, "update(file) : test test test test test test test test test test test test test test", m.Context["message"], "Must contains original message")
}

func TestMatchCommitQuery(t *testing.T) {
	err := exec.Command("../features/repo.sh").Run()

	if err != nil {
		logrus.Fatal(err)
	}

	cmd := exec.Command("git", "rev-parse", "test")
	cmd.Dir = "testing-repository"

	ID, err := cmd.Output()

	if err != nil {
		logrus.Fatal(err)
	}

	q := CommitQuery{
		"testing-repository/",
		string(ID[:len(ID)-1]),
		map[string]string{"simple": "(?:update|feat)\\(.*?\\) : .*?\\n\\n.*?\\n"},
		Options{
			CheckSummaryLength:  true,
			ExcludeMergeCommits: false,
			SummaryLength:       50,
		},
	}

	m, err := MatchCommitQuery(q)

	assert.NoError(t, err, "Must return no error")
	assert.Equal(t, Matching{}, *m, "Must return an empty matching struct")
}

func TestMatchCommitQueryWithCommitMessageThatDoesntMatchTemplate(t *testing.T) {
	err := exec.Command("../features/repo.sh").Run()

	if err != nil {
		logrus.Fatal(err)
	}

	cmd := exec.Command("git", "rev-parse", "test")
	cmd.Dir = "testing-repository"

	ID, err := cmd.Output()

	if err != nil {
		logrus.Fatal(err)
	}

	q := CommitQuery{
		"testing-repository/",
		string(ID[:len(ID)-1]),
		map[string]string{"simple": "whatever"},
		Options{
			CheckSummaryLength:  true,
			ExcludeMergeCommits: false,
			SummaryLength:       50,
		},
	}

	m, err := MatchCommitQuery(q)

	assert.NoError(t, err, "Must return no error")
	assert.Equal(t, "feat(file8) : new file 8\n\ncreate a new file 8\n", m.Context["message"], "Must contains commit message")
	assert.Equal(t, string(ID[:len(ID)-1]), m.Context["ID"], "Must contains commit id")
	assert.EqualError(t, m.MessageError, "No template match commit message", "Must contains commit message error")
	assert.NoError(t, m.SummaryError, "Must not contains error")
}

func TestMatchCommitQueryWithWrongRepository(t *testing.T) {
	err := exec.Command("../features/repo.sh").Run()

	if err != nil {
		logrus.Fatal(err)
	}

	cmd := exec.Command("git", "rev-parse", "test")
	cmd.Dir = "testing-repository"

	ID, err := cmd.Output()

	if err != nil {
		logrus.Fatal(err)
	}

	q := CommitQuery{
		"testtestest/",
		string(ID[:len(ID)-1]),
		map[string]string{"simple": "whatever"},
		Options{
			CheckSummaryLength:  true,
			ExcludeMergeCommits: false,
			SummaryLength:       50,
		},
	}

	_, err = MatchCommitQuery(q)

	assert.Error(t, err, "Must return an error")
}

func TestMatchCommitQueryWithAnUnexistingCommit(t *testing.T) {
	err := exec.Command("../features/repo.sh").Run()

	if err != nil {
		logrus.Fatal(err)
	}

	q := CommitQuery{
		"testing-repository/",
		"4e1243bd22c66e76c2ba9eddc1f91394e57f9f83",
		map[string]string{"simple": "whatever"},
		Options{
			CheckSummaryLength:  true,
			ExcludeMergeCommits: false,
			SummaryLength:       50,
		},
	}

	_, err = MatchCommitQuery(q)

	assert.Error(t, err, "Must return an error")
}

func TestIsValidSummaryLengthWithCorrectSize(t *testing.T) {
	assert.True(t, isValidSummaryLength(50, "test"))
	assert.True(t, isValidSummaryLength(50, "a sequence which is 50 size long abcdefghijklmnopq"), "Must have a length which is exactly 50 characters")
	assert.True(t, isValidSummaryLength(72, "test"))
	assert.True(t, isValidSummaryLength(72, "a sequence which is 72 size long abcdefghijklmnopqrstuvwxyz abcdefghijkl"), "Must have a length which is exactly 72 characters")
}

func TestIsValidSummaryLengthWithInCorrectSize(t *testing.T) {
	assert.False(t, isValidSummaryLength(50, "a sequence which is 51 size long abcdefghijklmnopqr"))
	assert.False(t, isValidSummaryLength(72, "a sequence which is 73 size long abcdefghijklmnopqrstuvwxyz abcdefghijklm"))
}

func TestIsMergeCommitWithANonMergeCommit(t *testing.T) {
	err := exec.Command("../features/repo.sh").Run()

	if err != nil {
		logrus.Fatal(err)
	}

	commits, err := fetchCommits("testing-repository", "test~2", "test")

	if err != nil {
		logrus.Fatal(err)
	}

	assert.False(t, isMergeCommit((*commits)[0]), "Must return false with non merge commit")
}

func TestIsMergeCommitWithAMergeCommit(t *testing.T) {
	for _, filename := range []string{"../features/repo.sh"} {
		err := exec.Command(filename).Run()

		if err != nil {
			logrus.Fatal(err)
		}
	}

	commits, err := fetchCommits("testing-repository", "test~3", "test~2")

	if err != nil {
		logrus.Fatal(err)
	}

	assert.True(t, isMergeCommit((*commits)[0]), "Must return false with non merge commit")
}
