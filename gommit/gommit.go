package gommit

import (
	"fmt"
	"regexp"
	"strings"

	"gopkg.in/src-d/go-git.v4"

	"github.com/antham/gommit/reference"
)

// CommitError represents an error when something goes wrong
type CommitError struct {
	ID           string
	Message      string
	MessageError error
	SummaryError error
}

// Query to retrieves commits and do checking
type Query struct {
	Path     string
	From     string
	To       string
	Matchers map[string]string
	Options  map[string]bool
}

// maxSummarySize represents maximum length of accommit summary
const maxSummarySize = 50

// fetchCommits retrieves all commits done in repository between 2 commits references
func fetchCommits(repoPath string, from string, to string) (*[]*git.Commit, error) {
	repo, err := git.NewFilesystemRepository(repoPath + "/.git/")

	if err != nil {
		return nil, err
	}

	return reference.FetchCommitInterval(repo, from, to)
}

// messageMatchTemplate try to match a commit message against a regexp
func messageMatchTemplate(message string, template string) bool {
	r := regexp.MustCompile(template)

	g := r.FindStringSubmatch(message)

	return len(g) > 0 && g[0] == message
}

// isValidSummaryLength return true if size length is lower than 80 characters
func isValidSummaryLength(message string) bool {
	chunks := strings.Split(message, "\n")

	return len(chunks) == 0 || len(chunks[0]) <= maxSummarySize
}

// isMergeCommit return true if a commit is a merge commit
func isMergeCommit(commit *git.Commit) bool {
	return commit.NumParents() == 2
}

// analyzeCommit check if a commit fit a query
func analyzeCommit(commit *git.Commit, query *Query) *CommitError {
	if query.Options["exclude-merge-commits"] && isMergeCommit(commit) {
		return nil
	}

	messageError := fmt.Errorf("No template match commit message")
	var summaryError error

	if query.Options["check-summary-length"] {
		summaryError = fmt.Errorf("Commit summary length is greater than 50 characters")
	}

	for _, matcher := range query.Matchers {
		t := messageMatchTemplate(commit.Message, matcher)

		if t {
			messageError = nil
		}
	}

	if isValidSummaryLength(commit.Message) {
		summaryError = nil
	}

	if messageError != nil || (summaryError != nil && query.Options["check-summary-length"]) {
		return &CommitError{
			commit.ID().String(),
			commit.Message,
			messageError,
			summaryError,
		}
	}

	return nil
}

// RunMatching trigger regexp matching against a range message commits
func RunMatching(query Query) (*[]CommitError, error) {
	analysis := []CommitError{}

	commits, err := fetchCommits(query.Path, query.From, query.To)

	if err != nil {
		return &analysis, err
	}

	if len(*commits) == 0 {
		return &analysis, fmt.Errorf(`No commits found between "%s" and "%s"`, query.From, query.To)
	}

	for _, commit := range *commits {
		commitError := analyzeCommit(commit, &query)
		if commitError != nil {
			analysis = append(analysis, *commitError)
		}
	}

	return &analysis, nil
}
