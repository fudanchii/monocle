package git

import (
	"path"

	"github.com/libgit2/git2go"
)

type Object struct {
	repo *git.Repository
	err  error
}

func (o *Object) OpenRepo(dir string) {
	if o.err == nil {
		o.repo, o.err = git.OpenRepository(dir)
	}
}

func (o *Object) LookupCommit(rev string) *git.Commit {
	var (
		obj *git.Object
		cc  *git.Commit
	)

	if o.err != nil {
		return nil
	}

	if obj, o.err = o.repo.RevparseSingle(rev); o.err == nil {
		cc, o.err = o.repo.LookupCommit(obj.Id())
	}
	return cc
}

func (o *Object) WalkCommitTree(c *git.Commit, fn git.TreeWalkCallback) {
	var tree *git.Tree

	if o.err != nil {
		return
	}

	tree, o.err = c.Tree()
	if o.err == nil {
		o.err = tree.Walk(fn)
	}
}

func (o *Object) Diff(pc *git.Commit, cc *git.Commit) *git.Diff {
	var (
		ptree    *git.Tree
		ctree    *git.Tree
		diffOpts git.DiffOptions
		diff     *git.Diff
	)

	if o.err != nil {
		return diff
	}

	// Ignore error here since it's theoretically wont error in any circumstance
	// (1 > 0 &&  1 <= 1)
	diffOpts, _ = git.DefaultDiffOptions()

	if ctree, o.err = cc.Tree(); o.err != nil {
		return diff
	}

	if ptree, o.err = pc.Tree(); o.err == nil {
		diff, o.err = o.repo.DiffTreeToTree(ptree, ctree, &diffOpts)
	}

	return diff
}

func (o *Object) GetDiffDeltaEntries(diff *git.Diff) []string {
	var (
		ndelta  int
		entries []string
	)

	if o.err != nil {
		return entries
	}

	if ndelta, o.err = diff.NumDeltas(); o.err == nil {
		for i := 0; i < ndelta; i++ {
			var dd git.DiffDelta
			if dd, o.err = diff.GetDelta(i); o.err == nil {
				entries = append(entries, dd.NewFile.Path)
			}
		}
	}
	return entries

}

type Files struct {
	WorkDir string
	Entries []string
}

func FilesChanged(dir string, rev string) (Files, error) {
	var (
		result Files   = Files{WorkDir: dir}
		gitObj *Object = &Object{repo: nil, err: nil}
		cc     *git.Commit
	)

	if rev == "" {
		rev = "HEAD"
	}

	gitObj.OpenRepo(dir)
	if cc = gitObj.LookupCommit(rev); gitObj.err != nil {
		return result, gitObj.err
	}

	pc := cc.Parent(0)
	if pc == nil { // First commit
		gitObj.WalkCommitTree(cc, func(s string, entry *git.TreeEntry) int {
			result.Entries = append(result.Entries, path.Join(s, entry.Name))
			return 0
		})
		return result, gitObj.err
	}

	result.Entries = gitObj.GetDiffDeltaEntries(gitObj.Diff(pc, cc))

	return result, gitObj.err
}
