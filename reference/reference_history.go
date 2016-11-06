package reference

import (
	"bytes"
	"fmt"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/core"
)

// history match a SymbolicRefPathStmt against
// a repository to resolve commit history and commit reference
type history struct {
	stmt          *symbolicRefPathStmt
	commitRef     *git.Commit
	commitHistory []*git.Commit
}

// newHistory creates a new History object
func newHistory(stmt *symbolicRefPathStmt, repo *git.Repository) (*history, error) {
	commitHistory := []*git.Commit{}
	iterRef, err := repo.Refs()

	if err != nil {
		return nil, err
	}

	hash, err := resolveHash(stmt.branchName, iterRef)

	if err != nil {
		return nil, err
	}

	commitRef, _ := repo.Commit(hash)
	commitHistory = append(commitHistory, commitRef)

	for _, path := range stmt.refPath {
		parents := commitRef.Parents()
		parentCount := commitRef.NumParents()

		if parentCount == 0 {
			return nil, fmt.Errorf(`Can't find reference`)
		}

		for i := 1; i <= parentCount; i++ {
			commit, err := parents.Next()

			if err != nil {
				return nil, fmt.Errorf("Can't find parent")
			}

			if path == i {
				commitHistory = append(commitHistory, commit)

				commitRef = commit
			}
		}
	}

	return &history{stmt, commitRef, commitHistory}, nil
}

// resolveHash give hash commit for a given string reference
func resolveHash(branchName string, iter core.ReferenceIter) (core.Hash, error) {
	hash := core.Hash{}

	err := iter.ForEach(func(ref *core.Reference) error {
		if ref.Name().Short() == branchName {
			hash = ref.Hash()
		}

		return nil
	})

	if err != nil {
		return hash, err
	}

	if hash.IsZero() {
		return hash, fmt.Errorf(`Can't find reference "%s"`, branchName)
	}

	return hash, nil
}

// parseHistoryInterval return commits between two intervals
func parseHistoryInterval(from *history, to *history) (*[]*git.Commit, error) {
	results := []*git.Commit{}
	commits := []*git.Commit{}

	for i, c := range from.commitHistory {
		if c.ID() == to.commitRef.ID() {
			commits = from.commitHistory[i:]

			break
		}
	}

	for i := 0; i < len(commits)-1; i++ {
		cs, errs := parseTree(commits[i], commits[i+1])

		if len(errs) > 0 {
			return nil, fmt.Errorf("An error occured when retrieving commits between %s and %s", from.commitRef.ID(), to.commitRef.ID())
		}

		results = append(results, cs...)
	}

	return &results, nil
}

// parseTree recursively parse a given tree to extract commits till boundary is reached
func parseTree(commit *git.Commit, bound *git.Commit) ([]*git.Commit, []error) {
	commits := []*git.Commit{}
	errors := []error{}

	if commit.ID() == bound.ID() || commit.NumParents() == 0 {
		return commits, errors
	}

	commits = append(commits, commit)

	parents := []*git.Commit{}

	err := commit.Parents().ForEach(
		func(c *git.Commit) error {
			parents = append(parents, c)

			return nil
		})

	if err != nil {
		errors = append(errors, err)
		return commits, errors
	}

	if len(parents) == 2 {
		cs, errs := parseTree(parents[1], bound)
		errors = append(errors, errs...)
		commits = append(commits, cs...)
	}

	if len(parents) == 1 {
		cs, errs := parseTree(parents[0], bound)
		errors = append(errors, errs...)
		commits = append(commits, cs...)
	}

	return commits, errors
}

// FetchCommitInterval retrieves commit history in a given interval for a repository
func FetchCommitInterval(repo *git.Repository, from string, to string) (*[]*git.Commit, error) {

	refHistoryFrom := newParser(bytes.NewBufferString(from))
	fromStmt, err := refHistoryFrom.parseSymbolicReferencePath()

	if err != nil {
		return &[]*git.Commit{}, err
	}

	refHistoryTo := newParser(bytes.NewBufferString(to))
	toStmt, err := refHistoryTo.parseSymbolicReferencePath()

	if err != nil {
		return &[]*git.Commit{}, err
	}

	fromHistory, err := newHistory(fromStmt, repo)

	if err != nil {
		return &[]*git.Commit{}, err
	}

	toHistory, err := newHistory(toStmt, repo)

	if err != nil {
		return &[]*git.Commit{}, err
	}

	return parseHistoryInterval(fromHistory, toHistory)
}
