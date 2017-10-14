package git

import (
	"github.com/fudanchii/monocle/errors"
	"github.com/libgit2/git2go"
)

func FilesChanged(dir string, rev string) []string {
	var result []string

	if rev == "" {
		rev = "HEAD"
	}

	repo, err := git.OpenRepository(dir)
	errors.ErrCheck(err)

	obj, err := repo.RevparseSingle(rev)
	errors.ErrCheck(err)

	commit, err := repo.LookupCommit(obj.Id())
	errors.ErrCheck(err)

	parent := commit.Parent(0)
	errors.NilCheck(parent)

	ctree, err := commit.Tree()
	errors.ErrCheck(err)

	ptree, err := parent.Tree()
	errors.ErrCheck(err)

	diffOptions, err := git.DefaultDiffOptions()
	errors.ErrCheck(err)

	diff, err := repo.DiffTreeToTree(ptree, ctree, &diffOptions)
	errors.ErrCheck(err)

	ndelta, err := diff.NumDeltas()
	errors.ErrCheck(err)

	for i := 0; i < ndelta; i++ {
		dd, err := diff.GetDelta(i)
		errors.ErrCheck(err)
		result = append(result, dd.NewFile.Path)
	}

	return result
}
