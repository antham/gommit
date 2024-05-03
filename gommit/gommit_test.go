package gommit

import (
	"os/exec"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

func TestFetchCommits(t *testing.T) {
	type scenario struct {
		name      string
		arguments func() (repoPath string, from string, to string)
		test      func(*[]*object.Commit, error)
	}

	scenarios := []scenario{
		{
			"Fetch valid interval",
			func() (string, string, string) {
				return "testing-repository", "test~2", "test"
			},
			func(commits *[]*object.Commit, err error) {
				assert.NoError(t, err)
				assert.Len(t, *commits, 2)

				expected := []string{
					"feat(file8) : new file 8\n\ncreate a new file 8\n",
					"feat(file7) : new file 7\n\ncreate a new file 7\n",
				}

				for i, c := range *commits {
					assert.Equal(t, expected[i], c.Message)
				}
			},
		},
		{
			"Fetch wrong repository",
			func() (string, string, string) {
				return "testtesttest", "4906f72818c0185162a3ec9c39a711d7c2842d40", "test"
			},
			func(commits *[]*object.Commit, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, "repository does not exist")
				assert.Nil(t, commits)
			},
		},
		{
			"Fetch with wrong interval",
			func() (string, string, string) {
				return "testing-repository", "4906f72818c0185162a3ec9c39a711d7c2842d40", "maste"
			},
			func(commits *[]*object.Commit, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, `Reference "maste" can't be found in git repository`)
				assert.Nil(t, commits)
			},
		},
	}

	for _, s := range scenarios {
		t.Run(s.name, func(t *testing.T) {
			assert.NoError(t, exec.Command("../features/repo.sh").Run())
			s.test(fetchCommits(s.arguments()))
			assert.NoError(t, exec.Command("../features/repo-teardown.sh").Run())
		})
	}
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

func TestMatchRangeQuery(t *testing.T) {
	err := exec.Command("../features/repo.sh").Run()
	if err != nil {
		logrus.Fatal(err)
	}

	q := RangeQuery{
		Path:     "testing-repository/",
		From:     "test~2",
		To:       "test",
		Matchers: map[string]string{"simple": "(?:update|feat)\\(.*?\\) : .*?\\n\\n.*?\\n"},
		Options: Options{
			CheckSummaryLength:  false,
			ExcludeMergeCommits: false,
			SummaryLength:       50,
		},
	}

	m, err := MatchRangeQuery(q)

	assert.NoError(t, err, "Must return no errors")
	assert.Len(t, *m, 0, "Must return no items, match was successful for every commit")
}

func TestMatchRangeQueryrWithAMessageErrorCommit(t *testing.T) {
	err := exec.Command("../features/repo.sh").Run()
	if err != nil {
		logrus.Fatal(err)
	}

	q := RangeQuery{
		Path:     "testing-repository/",
		From:     "test~2",
		To:       "test",
		Matchers: map[string]string{"simple": "(?:update)\\(.*?\\) : .*?\\n\\n.*?\\n"},
		Options: Options{
			CheckSummaryLength:  false,
			ExcludeMergeCommits: false,
			SummaryLength:       50,
		},
	}

	m, err := MatchRangeQuery(q)

	assert.NoError(t, err, "Must return no errors")
	assert.Len(t, *m, 2, "Must return two items")
	assert.Equal(t, "feat(file8) : new file 8\n\ncreate a new file 8\n", (*m)[0].Context["message"], "Must contains commit message")
	assert.EqualError(t, (*m)[0].MessageError, "No template match commit message", "Must contains commit message error")
	assert.NoError(t, (*m)[0].SummaryError, "Must not contains error")
}

func TestMatchRangeQueryWithASummaryErrorCommit(t *testing.T) {
	for _, filename := range []string{"../features/repo.sh", "../features/bad-summary-message-commit.sh"} {
		err := exec.Command(filename).Run()
		if err != nil {
			logrus.Fatal(err)
		}
	}

	q := RangeQuery{
		Path:     "testing-repository/",
		From:     "test~1",
		To:       "test",
		Matchers: map[string]string{"simple": ".*\n"},
		Options: Options{
			CheckSummaryLength:  true,
			ExcludeMergeCommits: false,
			SummaryLength:       50,
		},
	}

	m, err := MatchRangeQuery(q)

	assert.NoError(t, err, "Must return no errors")
	assert.Len(t, *m, 1, "Must return one item")
	assert.Equal(t, "A very long summary commit greater than minimum length 50\n", (*m)[0].Context["message"], "Must contains commit message")
	assert.NoError(t, (*m)[0].MessageError, "Must not contains error")
	assert.EqualError(t, (*m)[0].SummaryError, "Commit summary length is greater than 50 characters", "Must contains summary message error")
}

func TestMatchRangeWithAMessageErrorCommitWithoutMergeCommit(t *testing.T) {
	for _, filename := range []string{"../features/repo.sh"} {
		err := exec.Command(filename).Run()
		if err != nil {
			logrus.Fatal(err)
		}
	}

	q := RangeQuery{
		Path:     "testing-repository/",
		From:     "test^^^^",
		To:       "test",
		Matchers: map[string]string{"simple": "(?:update)\\(.*?\\) : .*?\\n\\n.*?\\n"},
		Options: Options{
			CheckSummaryLength:  false,
			ExcludeMergeCommits: true,
			SummaryLength:       50,
		},
	}

	m, err := MatchRangeQuery(q)

	assert.NoError(t, err, "Must return no errors")
	assert.Len(t, *m, 7)
	for i := 0; i < 7; i++ {
		assert.NotContains(t, (*m)[i].Context["message"], "Merge")
	}
	assert.EqualError(t, (*m)[0].MessageError, "No template match commit message", "Must contains commit message error")
	assert.NoError(t, (*m)[0].SummaryError, "Must not contains error")
}

func TestMatchRangeQueryWithAMessageErrorCommitWithMergeCommits(t *testing.T) {
	for _, filename := range []string{"../features/repo.sh"} {
		err := exec.Command(filename).Run()
		if err != nil {
			logrus.Fatal(err)
		}
	}

	q := RangeQuery{
		Path:     "testing-repository/",
		From:     "test^^^^",
		To:       "test",
		Matchers: map[string]string{"simple": "(?:update)\\(.*?\\) : .*?\\n\\n.*?\\n"},
		Options: Options{
			CheckSummaryLength:  false,
			ExcludeMergeCommits: false,
			SummaryLength:       50,
		},
	}

	m, err := MatchRangeQuery(q)

	assert.NoError(t, err, "Must return no errors")
	assert.Len(t, *m, 9, "Must return two itesm")
	assert.EqualError(t, (*m)[0].MessageError, "No template match commit message", "Must contains commit message error")
	assert.NoError(t, (*m)[0].SummaryError, "Must not contains error")
}

func TestMatchRangeWithAnUnexistingCommitRange(t *testing.T) {
	err := exec.Command("../features/repo.sh").Run()
	if err != nil {
		logrus.Fatal(err)
	}

	q := RangeQuery{
		Path:     "testing-repository/",
		From:     "test~15",
		To:       "test",
		Matchers: map[string]string{"simple": "(?:update)\\(.*?\\) : .*?\\n\\n.*?\\n"},
		Options: Options{
			CheckSummaryLength:  false,
			ExcludeMergeCommits: false,
			SummaryLength:       50,
		},
	}

	m, err := MatchRangeQuery(q)

	assert.EqualError(t, err, "Reference \"test~15\" can't be found in git repository")
	assert.Len(t, *m, 0, "Must return no item")
}

func TestMatchRangeWithNoCommitsInCommitRange(t *testing.T) {
	err := exec.Command("../features/repo.sh").Run()
	if err != nil {
		logrus.Fatal(err)
	}

	q := RangeQuery{
		Path:     "testing-repository/",
		From:     "test",
		To:       "test",
		Matchers: map[string]string{"simple": "(?:update)\\(.*?\\) : .*?\\n\\n.*?\\n"},
		Options: Options{
			CheckSummaryLength:  false,
			ExcludeMergeCommits: false,
			SummaryLength:       50,
		},
	}

	m, err := MatchRangeQuery(q)

	assert.EqualError(t, err, `Can't produce a diff between test and test, check your range is correct by running "git log test..test" command`)
	assert.Len(t, *m, 0, "Must return no item")
}

func TestMatchMessageQuery(t *testing.T) {
	err := exec.Command("../features/repo.sh").Run()
	if err != nil {
		logrus.Fatal(err)
	}

	q := MessageQuery{
		Message:  "update(file) : fix",
		Matchers: map[string]string{"simple": "(?:update|feat)\\(.*?\\) : .*"},
		Options: Options{
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
		Message:  "update(file) :",
		Matchers: map[string]string{"simple": "(?:update|feat)\\(.*?\\) : .*"},
		Options: Options{
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
		Message:  "update(file) : test test test test test test test test test test test test test test",
		Matchers: map[string]string{"simple": "(?:update|feat)\\(.*?\\) : .*"},
		Options: Options{
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
		Path:     "testing-repository/",
		ID:       string(ID[:len(ID)-1]),
		Matchers: map[string]string{"simple": "(?:update|feat)\\(.*?\\) : .*?\\n\\n.*?\\n"},
		Options: Options{
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
		Path:     "testing-repository/",
		ID:       string(ID[:len(ID)-1]),
		Matchers: map[string]string{"simple": "whatever"},
		Options: Options{
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
		Path:     "testtestest/",
		ID:       string(ID[:len(ID)-1]),
		Matchers: map[string]string{"simple": "whatever"},
		Options: Options{
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
		Path:     "testing-repository/",
		ID:       "4e1243bd22c66e76c2ba9eddc1f91394e57f9f83",
		Matchers: map[string]string{"simple": "whatever"},
		Options: Options{
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
