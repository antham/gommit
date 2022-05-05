package reference

import (
	"fmt"
	"os"
	"os/exec"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/sirupsen/logrus"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
)

var repo *git.Repository
var gitRepositoryPath = "testing-repository"

func setup() {
	err := exec.Command("../features/repo.sh").Run()

	if err != nil {
		logrus.Fatal(err)
	}

	path, err := os.Getwd()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	repo, err = git.PlainOpen(path + "/" + gitRepositoryPath)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getCommitFromRef(ref string) *object.Commit {
	cmd := exec.Command("git", "rev-parse", ref)
	cmd.Dir = gitRepositoryPath

	ID, err := cmd.Output()
	ID = ID[:len(ID)-1]

	if err != nil {
		logrus.WithField("ID", string(ID)).Fatal(err)
	}

	c, err := repo.CommitObject(plumbing.NewHash(string(ID)))

	if err != nil {
		logrus.WithField("ID", ID).Fatal(err)
	}

	return c
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func TestResolveRef(t *testing.T) {
	type g struct {
		ref string
		f   func(*object.Commit, error)
	}

	tests := []g{
		{
			"HEAD",
			func(o *object.Commit, err error) {
				assert.NoError(t, err)
				assert.True(t, o.ID().String() == getCommitFromRef("HEAD").ID().String(), "Must resolve HEAD reference")
			},
		},
		{
			"test1",
			func(o *object.Commit, err error) {
				assert.NoError(t, err)
				assert.True(t, o.ID().String() == getCommitFromRef("test1").ID().String(), "Must resolve branch reference")
			},
		},
		{
			"test1",
			func(o *object.Commit, err error) {
				assert.NoError(t, err)
				assert.True(t, o.ID().String() == getCommitFromRef("test1").ID().String(), "Must resolve commit id")
			},
		},
		{
			"whatever",
			func(o *object.Commit, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, `Reference "whatever" can't be found in git repository`)
			},
		},
	}

	for _, test := range tests {
		test.f(resolveRef(test.ref, repo))
	}
}

func TestResolveRefWithErrors(t *testing.T) {
	type g struct {
		ref  string
		repo *git.Repository
		f    func(*object.Commit, error)
	}

	tests := []g{
		{
			"whatever",
			repo,
			func(o *object.Commit, err error) {
				assert.Error(t, err)
				assert.EqualError(t, err, `Reference "whatever" can't be found in git repository`)
			},
		},
	}

	for _, test := range tests {
		test.f(resolveRef(test.ref, test.repo))
	}
}

func TestFetchCommitInterval(t *testing.T) {
	type g struct {
		toRef   string
		fromRef string
		f       func(*[]*object.Commit, error)
	}
	tests := []g{
		{
			"HEAD",
			"test",
			func(cs *[]*object.Commit, err error) {
				assert.Error(t, err)
				assert.Regexp(t, `Can't produce a diff between .*? and .*?, check your range is correct by running "git log .*?\.\..*?" command`, err.Error())
			},
		},
		{
			"HEAD~1",
			"HEAD~3",
			func(cs *[]*object.Commit, err error) {
				assert.Error(t, err)
				assert.Regexp(t, `Can't produce a diff between .*? and .*?, check your range is correct by running "git log .*?\.\..*?" command`, err.Error())
			},
		},
		{
			"HEAD~3",
			"test~2^2",
			func(cs *[]*object.Commit, err error) {
				assert.NoError(t, err)
				assert.Len(t, *cs, 5)

				commitTests := []string{
					"Merge branch 'test2' into test1\n",
					"feat(file6) : new file 6\n\ncreate a new file 6\n",
					"feat(file5) : new file 5\n\ncreate a new file 5\n",
					"feat(file4) : new file 4\n\ncreate a new file 4\n",
					"feat(file3) : new file 3\n\ncreate a new file 3\n",
				}

				for i, c := range *cs {
					assert.Equal(t, commitTests[i], c.Message)
				}
			},
		},
		{
			"HEAD~4",
			"test~2^2^2",
			func(cs *[]*object.Commit, err error) {
				assert.NoError(t, err)
				assert.Len(t, *cs, 5)

				commitTests := []string{
					"feat(file6) : new file 6\n\ncreate a new file 6\n",
					"feat(file5) : new file 5\n\ncreate a new file 5\n",
					"feat(file4) : new file 4\n\ncreate a new file 4\n",
					"feat(file3) : new file 3\n\ncreate a new file 3\n",
					"feat(file2) : new file 2\n\ncreate a new file 2\n",
				}

				for i, c := range *cs {
					assert.Equal(t, commitTests[i], c.Message)
				}
			},
		},
		{
			getCommitFromRef("HEAD~4").ID().String(),
			getCommitFromRef("test~2^2^2").ID().String(),
			func(cs *[]*object.Commit, err error) {
				assert.NoError(t, err)
				assert.Len(t, *cs, 5)

				commitTests := []string{
					"feat(file6) : new file 6\n\ncreate a new file 6\n",
					"feat(file5) : new file 5\n\ncreate a new file 5\n",
					"feat(file4) : new file 4\n\ncreate a new file 4\n",
					"feat(file3) : new file 3\n\ncreate a new file 3\n",
					"feat(file2) : new file 2\n\ncreate a new file 2\n",
				}

				for i, c := range *cs {
					assert.Equal(t, commitTests[i], c.Message)
				}
			},
		},
		{
			"whatever",
			"HEAD~1",
			func(cs *[]*object.Commit, err error) {
				assert.EqualError(t, err, `Reference "whatever" can't be found in git repository`)
			},
		},
		{
			"HEAD~1",
			"whatever",
			func(cs *[]*object.Commit, err error) {
				assert.EqualError(t, err, `Reference "whatever" can't be found in git repository`)
			},
		},
		{
			"HEAD",
			"HEAD",
			func(cs *[]*object.Commit, err error) {
				assert.EqualError(t, err, `Can't produce a diff between HEAD and HEAD, check your range is correct by running "git log HEAD..HEAD" command`)
			},
		},
	}
	for _, test := range tests {
		test.f(FetchCommitInterval(repo, test.toRef, test.fromRef))
	}
}

func TestShallowCloneProducesNoErrors(t *testing.T) {
	repositoryPath := "shallow-repository-test"
	cmd := exec.Command("rm", "-rf", repositoryPath)
	_, err := cmd.Output()
	assert.NoError(t, err)

	cmd = exec.Command("git", "clone", "--depth", "2", "https://github.com/octocat/Spoon-Knife.git", repositoryPath)
	_, err = cmd.Output()
	assert.NoError(t, err)

	path, err := os.Getwd()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	repo, err := git.PlainOpen(path + "/" + repositoryPath)
	if err != nil {
		logrus.Fatal(err)
	}

	commits, err := FetchCommitInterval(repo, "HEAD~1", "HEAD")
	assert.NoError(t, err)
	assert.Len(t, *commits, 1, "Must fetch commits in shallow clone")
}

func TestFetchCommitByID(t *testing.T) {
	cmd := exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = gitRepositoryPath

	ID, err := cmd.Output()

	if err != nil {
		logrus.Fatal(err)
	}

	commit, err := FetchCommitByID(repo, string(ID[:len(ID)-1]))

	assert.NoError(t, err, "Must return no errors")
	assert.Equal(t, "feat(file8) : new file 8\n\ncreate a new file 8\n", commit.Message, "Must return commit linked to this id")
}

func TestFetchCommitByIDWithAWrongCommitID(t *testing.T) {
	_, err := FetchCommitByID(repo, "whatever")

	assert.Error(t, err, "Must return an error")
}
