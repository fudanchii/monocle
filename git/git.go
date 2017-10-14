package git

import (
	"path"

	"github.com/fudanchii/monocle/errors"
	"github.com/libgit2/git2go"
)

type Files struct {
	WorkDir string
	Entries []string
}

func FilesChanged(dir string, rev string) Files {
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

	ctree, err := commit.Tree()
	errors.ErrCheck(err)

	parent := commit.Parent(0)

	if parent == nil {
		err = ctree.Walk(func(s string, entry *git.TreeEntry) int {
			result = append(result, path.Join(s, entry.Name))
			return 0
		})
		errors.ErrCheck(err)
	} else {
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
	}

	return Files{
		WorkDir: dir,
		Entries: result,
	}
}
