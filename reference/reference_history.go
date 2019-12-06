package reference

import (
	"fmt"
	"io"
	"regexp"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

// node is a tree node in commit tree
type node struct {
	value  *object.Commit
	parent *node
}

// errNoDiffBetweenReferences is triggered when we can't
// produce any diff between 2 references
type errNoDiffBetweenReferences struct {
	from string
	to   string
}

func (e errNoDiffBetweenReferences) Error() string {
	return fmt.Sprintf(`Can't produce a diff between %s and %s, check your range is correct by running "git log %[1]s..%[2]s" command`, e.from, e.to)
}

// errReferenceNotFound is triggered when reference can't be
// found in git repository
type errReferenceNotFound struct {
	ref string
}

func (e errReferenceNotFound) Error() string {
	return fmt.Sprintf(`Reference "%s" can't be found in git repository`, e.ref)
}

// errBrowsingTree is triggered when something wrong occurred during commit analysis process
var errBrowsingTree = fmt.Errorf("An issue occurred during tree analysis")

// FetchCommitInterval retrieves commit refSolver in a given interval for a repository
func FetchCommitInterval(repo *git.Repository, from string, to string) (*[]*object.Commit, error) {
	fromCommit, err := resolveRef(from, repo)
	if err != nil {
		return nil, err
	}

	toCommit, err := resolveRef(to, repo)
	if err != nil {
		return nil, err
	}

	var ok bool

	exclusionList, err := buildOriginCommitList(fromCommit)
	if err != nil {
		return nil, err
	}

	if _, ok = exclusionList[toCommit.ID().String()]; ok {
		return nil, errNoDiffBetweenReferences{from, to}
	}

	commits, err := findDiffCommits(toCommit, &exclusionList)
	if err != nil {
		return nil, err
	}

	if len(*commits) == 0 {
		return nil, errNoDiffBetweenReferences{from, to}
	}

	return commits, nil
}

// FetchCommitByID retrieves a single commit from a repository
func FetchCommitByID(repo *git.Repository, ID string) (*object.Commit, error) {
	hash := plumbing.NewHash(ID)

	return repo.CommitObject(hash)
}

// resolveRef gives hash commit for a given string reference
func resolveRef(refCommit string, repository *git.Repository) (*object.Commit, error) {
	hash, err := repository.ResolveRevision(plumbing.Revision(refCommit))

	if (err != nil || hash.IsZero()) && regexp.MustCompile("[0-9a-f]{40}").MatchString(refCommit) {
		i, cErr := repository.CommitObjects()

		if cErr != nil {
			return nil, errReferenceNotFound{refCommit}
		}

		var c *object.Commit

		cErr = i.ForEach(func(o *object.Commit) error {
			if o.ID().String() == refCommit {
				c = o

				return io.EOF
			}

			return nil
		})

		if cErr != nil && cErr != io.EOF {
			return nil, errReferenceNotFound{refCommit}
		}

		return c, nil
	}

	if err == nil && !hash.IsZero() {
		return repository.CommitObject(*hash)
	}

	return &object.Commit{}, errReferenceNotFound{refCommit}
}

// buildOriginCommitList browses git tree from a given commit
// till root commit using kind of breadth first search algorithm
// and grab commit ID to a map with ID as key
func buildOriginCommitList(commit *object.Commit) (map[string]bool, error) {
	queue := append([]*object.Commit{}, commit)
	seen := map[string]bool{commit.ID().String(): true}

	for len(queue) > 0 {
		current := queue[0]
		queue = append([]*object.Commit{}, queue[1:]...)

		err := current.Parents().ForEach(
			func(c *object.Commit) error {
				if _, ok := seen[c.ID().String()]; !ok {
					seen[c.ID().String()] = true
					queue = append(queue, c)
				}

				return nil
			})

		if err != nil && err.Error() != plumbing.ErrObjectNotFound.Error() {
			return seen, errBrowsingTree
		}
	}

	return seen, nil
}

// findDiffCommits extracts commits that are no part of a given commit list
// using kind of depth first search algorithm to keep commits ordered
func findDiffCommits(commit *object.Commit, exclusionList *map[string]bool) (*[]*object.Commit, error) {
	commits := []*object.Commit{}
	queue := append([]*node{}, &node{value: commit})
	seen := map[string]bool{commit.ID().String(): true}
	var current *node

	for len(queue) > 0 {
		current = queue[0]
		queue = append([]*node{}, queue[1:]...)

		if _, ok := (*exclusionList)[current.value.ID().String()]; !ok {
			commits = append(commits, current.value)
		}

		err := current.value.Parents().ForEach(
			func(c *object.Commit) error {
				if _, ok := seen[c.ID().String()]; !ok {
					seen[c.ID().String()] = true
					n := &node{value: c, parent: current}
					queue = append([]*node{n}, queue...)
				}

				return nil
			})

		if err != nil && err.Error() != plumbing.ErrObjectNotFound.Error() {
			return &commits, errBrowsingTree
		}
	}

	return &commits, nil
}
