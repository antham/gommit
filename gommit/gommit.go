package gommit

import (
	"bytes"
	"fmt"

	"github.com/glenn-brown/golang-pkg-pcre/src/pkg/pcre"
	"github.com/libgit2/git2go"
)

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

// RunMatching trigger regexp matching against a range message commits
func RunMatching(path string, from string, till string, matchers map[string]string) (*[]map[string]string, error) {
	analysis := []map[string]string{}

	commits, err := fetchCommits(path, from, till)

	if err != nil {
		return &analysis, fmt.Errorf(`Interval between "%s" and "%s" can't be fetched`, from, till)
	}

	if len(*commits) == 0 {
		return &analysis, fmt.Errorf(`No commits found between "%s" and "%s"`, from, till)
	}

	for _, commit := range *commits {
		var ok bool

		for _, matcher := range matchers {
			t, _ := messageMatchTemplate(commit.Message(), matcher)

			if t {
				ok = true
			}
		}

		if !ok {
			analysis = append(analysis, map[string]string{
				"id":      commit.Id().String(),
				"message": commit.Message(),
			})
		}
	}

	return &analysis, nil
}
