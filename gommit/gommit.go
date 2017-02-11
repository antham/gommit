package gommit

import (
	"fmt"
	"regexp"
	"strings"

	"srcd.works/go-git.v4"
	"srcd.works/go-git.v4/plumbing/object"

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
	Options  map[string]bool
}

// RangeCommitQuery to retrieves commits and do checking
type RangeCommitQuery struct {
	Path     string
	From     string
	To       string
	Matchers map[string]string
	Options  map[string]bool
}

// MessageQuery to check only commit message
type MessageQuery struct {
	Message  string
	Matchers map[string]string
	Options  map[string]bool
}

// maxSummarySize represents maximum length of accommit summary
const maxSummarySize = 50

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
func isMergeCommit(commit *object.Commit) bool {
	return commit.NumParents() == 2
}

// IsZeroMatching check if Matching struct equals zero
func IsZeroMatching(matching *Matching) bool {
	return len(matching.Context) == 0 && matching.MessageError == nil && matching.SummaryError == nil
}

// analyzeMessage check if a message match
func analyzeMessage(message string, matchers map[string]string, options map[string]bool) *Matching {
	matching := Matching{}
	matchTemplate := false
	hasError := false

	for _, matcher := range matchers {
		t := messageMatchTemplate(message, matcher)

		if t {
			matchTemplate = true
		}
	}

	if options["check-summary-length"] && !isValidSummaryLength(message) {
		hasError = true
		matching.SummaryError = fmt.Errorf("Commit summary length is greater than 50 characters")
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

// analyzeCommit check if a commit match
func analyzeCommit(commit *object.Commit, matchers map[string]string, options map[string]bool) *Matching {
	if options["exclude-merge-commits"] && isMergeCommit(commit) {
		return &Matching{}
	}

	m := analyzeMessage(commit.Message, matchers, options)

	if IsZeroMatching(m) {
		return &Matching{}
	}

	m.Context["ID"] = commit.ID().String()

	return m
}

// analyzeCommits check if a slice of commits match
func analyzeCommits(commits *[]*object.Commit, matchers map[string]string, options map[string]bool) *[]*Matching {
	matchings := []*Matching{}

	for _, commit := range *commits {
		matching := analyzeCommit(commit, matchers, options)

		if !IsZeroMatching(matching) {
			matchings = append(matchings, matching)
		}
	}

	return &matchings
}

// MatchMessageQuery trigger regexp matching against a message
func MatchMessageQuery(query MessageQuery) (*Matching, error) {
	return analyzeMessage(query.Message, query.Matchers, query.Options), nil
}

// MatchCommitQuery trigger regexp matching against a commit
func MatchCommitQuery(query CommitQuery) (*Matching, error) {
	commit, err := fetchCommit(query.Path, query.ID)

	if err != nil {
		return &Matching{}, err
	} else if commit == nil {
		return &Matching{}, fmt.Errorf(`No commits found between with ID : "%s"`, query.ID)
	}

	return analyzeCommit(commit, query.Matchers, query.Options), nil
}

// MatchRangeCommitQuery trigger regexp matching against a range message commits
func MatchRangeCommitQuery(query RangeCommitQuery) (*[]*Matching, error) {
	commits, err := fetchCommits(query.Path, query.From, query.To)

	if err != nil {
		return &[]*Matching{}, err
	} else if len(*commits) == 0 {
		return &[]*Matching{}, fmt.Errorf(`No commits found between "%s" and "%s"`, query.From, query.To)
	}

	return analyzeCommits(commits, query.Matchers, query.Options), nil
}
