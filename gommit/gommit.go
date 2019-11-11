package gommit

import (
	"fmt"
	"strings"

	"github.com/dlclark/regexp2"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"

	"github.com/antham/gommit/reference"
)

// Matching represents an error when something goes wrong
type Matching struct {
	Context      map[string]string
	MessageError error
	SummaryError error
}

// CommitQuery to retrieves a commit and do checking
type CommitQuery struct {
	Path     string
	ID       string
	Matchers map[string]string
	Options  Options
}

// RangeQuery to retrieves commits and do checking
type RangeQuery struct {
	Path     string
	From     string
	To       string
	Matchers map[string]string
	Options  Options
}

// MessageQuery to check only commit message
type MessageQuery struct {
	Message  string
	Matchers map[string]string
	Options  Options
}

// Options represents options picked from configuration
type Options struct {
	CheckSummaryLength  bool
	ExcludeMergeCommits bool
	SummaryLength       int
}

// fetchCommits retrieves all commits in repository between 2 commits references
func fetchCommits(repoPath string, from string, to string) (*[]*object.Commit, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	return reference.FetchCommitInterval(repo, from, to)
}

// fetchCommit retrieve a single commit in repository from its ID
func fetchCommit(repoPath string, ID string) (*object.Commit, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, err
	}

	return reference.FetchCommitByID(repo, ID)
}

// messageMatchTemplate tries to match a commit message against a regexp
func messageMatchTemplate(message string, template string) bool {
	r := regexp2.MustCompile(template, 0)

	b, err := r.MatchString(message)

	return err == nil && b
}

// isValidSummaryLength returns true if message size length is lower than characters given by
// summaryLength
func isValidSummaryLength(summaryLength int, message string) bool {
	chunks := strings.Split(message, "\n")

	return len(chunks) == 0 || len(chunks[0]) <= summaryLength
}

// isMergeCommit returns true if a commit is a merge commit
func isMergeCommit(commit *object.Commit) bool {
	return commit.NumParents() == 2
}

// IsZeroMatching checks if Matching struct equals zero
func IsZeroMatching(matching *Matching) bool {
	return len(matching.Context) == 0 && matching.MessageError == nil && matching.SummaryError == nil
}

// analyzeMessage checks if a message match expectations
func analyzeMessage(message string, matchers map[string]string, options Options) *Matching {
	matching := Matching{}
	matchTemplate := false
	hasError := false

	for _, matcher := range matchers {
		t := messageMatchTemplate(message, matcher)

		if t {
			matchTemplate = true
		}
	}

	if options.CheckSummaryLength && !isValidSummaryLength(options.SummaryLength, message) {
		hasError = true
		matching.SummaryError = fmt.Errorf("Commit summary length is greater than %d characters", options.SummaryLength)
	}

	if !matchTemplate {
		hasError = true
		matching.MessageError = fmt.Errorf("No template match commit message")
	}

	if hasError {
		matching.Context = map[string]string{"message": message}
	}

	return &matching
}

// analyzeCommit checks if a commit message match expectations
func analyzeCommit(commit *object.Commit, matchers map[string]string, options Options) *Matching {
	if options.ExcludeMergeCommits && isMergeCommit(commit) {
		return &Matching{}
	}

	m := analyzeMessage(commit.Message, matchers, options)

	if IsZeroMatching(m) {
		return &Matching{}
	}

	m.Context["ID"] = commit.ID().String()

	return m
}

// analyzeCommits checks if a slice of commits message match expectations
func analyzeCommits(commits *[]*object.Commit, matchers map[string]string, options Options) *[]*Matching {
	matchings := []*Matching{}

	for _, commit := range *commits {
		matching := analyzeCommit(commit, matchers, options)

		if !IsZeroMatching(matching) {
			matchings = append(matchings, matching)
		}
	}

	return &matchings
}

// MatchMessageQuery triggers regexp matching against a message
func MatchMessageQuery(query MessageQuery) (*Matching, error) {
	return analyzeMessage(query.Message, query.Matchers, query.Options), nil
}

// MatchCommitQuery triggers regexp matching against a commit
func MatchCommitQuery(query CommitQuery) (*Matching, error) {
	commit, err := fetchCommit(query.Path, query.ID)

	if err != nil {
		return &Matching{}, err
	}

	return analyzeCommit(commit, query.Matchers, query.Options), nil
}

// MatchRangeQuery triggers regexp matching against a range of commit messages
func MatchRangeQuery(query RangeQuery) (*[]*Matching, error) {
	commits, err := fetchCommits(query.Path, query.From, query.To)

	if err != nil {
		return &[]*Matching{}, err
	}

	return analyzeCommits(commits, query.Matchers, query.Options), nil
}
