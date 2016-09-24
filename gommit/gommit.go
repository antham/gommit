package gommit

import (
	"github.com/libgit2/git2go"
)

// FetchCommits retrieves all commits done in repository between 2 commits references
func FetchCommits(repoPath string, from string, till string) (*[]*git.Commit, error) {
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
