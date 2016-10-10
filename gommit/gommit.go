package gommit

import (
	"bytes"
	"fmt"

	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
	"github.com/libgit2/git2go"
)

// CommitError represents an error when something goes wrong
type CommitError struct {
	ID           string
	Message      string
	MessageError error
	SummaryError error
}

// MAX_SUMMARY_SIZE represents maximum length of accommit summary
const MAX_SUMMARY_SIZE = 50

// fetchCommits retrieves all commits done in repository between 2 commits references
func fetchCommits(repoPath string, from string, till string) (*[]*git.Commit, error) {
	commits := []*git.Commit{}

	repo, err := git.OpenRepository(repoPath)

	if err != nil {
		return &commits, err
	}

	w, err := repo.Walk()

	if err != nil {
		return &commits, err
	}

	err = w.PushRange(from + ".." + till)

	if err != nil {
		return &commits, err
	}

	prevOid := git.Oid{}
	currOid := &git.Oid{}

	for {
		err := w.Next(currOid)

		if git.IsErrorCode(err, git.ErrEOF) {
			return &commits, nil
		}

		if currOid.Equal(&prevOid) {
			return &commits, nil
		}

		prevOid = *currOid

		c, err := repo.LookupCommit(currOid)

		if err != nil {
			return &commits, err
		}

		commits = append(commits, c)
	}
}

// messageMatchTemplate try to match a commit message against a regexp
func messageMatchTemplate(message string, template string) (bool, string) {
	r := pcre.MustCompile(template, pcre.ANCHORED)

	msgByte := []byte(message)

	g := r.Matcher(msgByte, pcre.ANCHORED).Group(0)

	return bytes.Equal(msgByte, g), string(g)
}

// isValidSummaryLength return true if size length is lower than 80 characters
func isValidSummaryLength(summary string) bool {
	return len(summary) <= MAX_SUMMARY_SIZE
}

// isMergeCommit return true if a commit is a merge commit
func isMergeCommit(commit *git.Commit) bool {
	return commit.ParentCount() == 2
}

// RunMatching trigger regexp matching against a range message commits
func RunMatching(path string, from string, till string, matchers map[string]string, options map[string]bool) (*[]CommitError, error) {
	analysis := []CommitError{}

	commits, err := fetchCommits(path, from, till)

	if err != nil {
		return &analysis, fmt.Errorf(`Interval between "%s" and "%s" can't be fetched`, from, till)
	}

	if len(*commits) == 0 {
		return &analysis, fmt.Errorf(`No commits found between "%s" and "%s"`, from, till)
	}

	for _, commit := range *commits {
		if options["exclude-merge-commits"] && isMergeCommit(commit) {
			continue
		}

		messageError := fmt.Errorf("No template match commit message")
		var summaryError error

		if options["check-summary-length"] {
			summaryError = fmt.Errorf("Commit summary length is greater than 50 characters")
		}

		for _, matcher := range matchers {
			t, _ := messageMatchTemplate(commit.Message(), matcher)

			if t {
				messageError = nil
			}
		}

		if isValidSummaryLength(commit.Summary()) {
			summaryError = nil
		}

		if messageError != nil || (summaryError != nil && options["check-summary-length"]) {
			analysis = append(analysis, CommitError{
				commit.Id().String(),
				commit.Message(),
				messageError,
				summaryError,
			})
		}
	}

	return &analysis, nil
}
