package reference

import (
	"bytes"
	"fmt"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/core"
)

// refSolver match a SymbolicRefPathStmt against
// a repository to resolve commit refSolver and commit reference
type refSolver struct {
	stmt      *symbolicRefPathStmt
	commitRef *git.Commit
}

// newRefSolver creates a new RefSolver object
func newRefSolver(stmt *symbolicRefPathStmt, repo *git.Repository) (*refSolver, error) {
	hash, err := resolveHash(stmt.branchName, repo)

	if err != nil {
		return nil, err
	}

	commitRef, _ := repo.Commit(hash)

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
				commitRef = commit
			}
		}
	}

	return &refSolver{stmt, commitRef}, nil
}

// resolveHash give hash commit for a given string reference
func resolveHash(refCommit string, repository *git.Repository) (core.Hash, error) {
	hash := core.Hash{}

	if strings.ToLower(refCommit) == "head" {
		head, err := repository.Head()

		if err == nil {
			return head.Hash(), nil
		}
	}

	iter, err := repository.Refs()

	if err != nil {
		return core.Hash{}, err
	}

	err = iter.ForEach(func(ref *core.Reference) error {
		if ref.Name().Short() == refCommit {
			hash = ref.Hash()
		}

		return nil
	})

	if err == nil && !hash.IsZero() {
		return hash, err
	}

	hash = core.NewHash(refCommit)

	_, err = repository.Commit(hash)

	if err == nil && !hash.IsZero() {
		return hash, nil
	}

	return hash, fmt.Errorf(`Can't find reference "%s"`, refCommit)
}

// retrieveCommitPath fetch all commits between 2 references
func retrieveCommitPath(from *git.Commit, to *git.Commit) (*[]*git.Commit, error) {
	results := []*git.Commit{}
	parents := []*git.Commit{}

	err := to.Parents().ForEach(
		func(c *git.Commit) error {
			parents = append(parents, c)

			return nil
		})

	if err != nil {
		return nil, err
	}

	if to.ID() == from.ID() {
		results = append(results, to)

		return &results, nil
	}

	for i := 0; i < to.NumParents(); i++ {
		cs, err := retrieveCommitPath(from, parents[i])

		if err != nil {
			return nil, err
		}

		if len(*cs) > 0 {
			results = append(results, to)
			results = append(results, *cs...)

			return &results, nil
		}
	}

	return &results, nil
}

// parseCommitHistory return commits between two intervals
func parseCommitHistory(from *refSolver, to *refSolver) (*[]*git.Commit, error) {
	results := []*git.Commit{}

	commits, err := retrieveCommitPath(from.commitRef, to.commitRef)

	if err != nil {
		return nil, err
	}

	for i := 0; i < len(*commits)-1; i++ {
		cs, errs := parseTree((*commits)[i], (*commits)[i+1])

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

// FetchCommitInterval retrieves commit refSolver in a given interval for a repository
func FetchCommitInterval(repo *git.Repository, from string, to string) (*[]*git.Commit, error) {
	refRefSolverFrom := newParser(bytes.NewBufferString(from))
	fromStmt, err := refRefSolverFrom.parseSymbolicReferencePath()

	if err != nil {
		return &[]*git.Commit{}, err
	}

	refRefSolverTo := newParser(bytes.NewBufferString(to))
	toStmt, err := refRefSolverTo.parseSymbolicReferencePath()

	if err != nil {
		return &[]*git.Commit{}, err
	}

	fromRefSolver, err := newRefSolver(fromStmt, repo)

	if err != nil {
		return &[]*git.Commit{}, err
	}

	toRefSolver, err := newRefSolver(toStmt, repo)

	if err != nil {
		return &[]*git.Commit{}, err
	}

	return parseCommitHistory(fromRefSolver, toRefSolver)
}

// FetchCommitByID retrieves a single commit from a repository
func FetchCommitByID(repo *git.Repository, ID string) (*git.Commit, error) {
	hash := core.NewHash(ID)

	return repo.Commit(hash)
}
