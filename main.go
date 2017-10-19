package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fudanchii/monocle/build"
	"github.com/fudanchii/monocle/git"
)

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage:\n\t%s <repo_dir> [git_rev]\n", os.Args[0])
		flag.PrintDefaults()
	}
}

func main() {
	var (
		workDir string = "."
		rev     string = "HEAD"
	)

	flag.Parse()
	switch flag.NArg() {
	case 2:
		rev = flag.Arg(1)
		fallthrough
	case 1:
		workDir = flag.Arg(0)
	}

	for _, buildFile := range build.CreateBuildEntries(git.FilesChanged(workDir, rev)) {
		build.Start(build.Name(buildFile), build.ParseManifest(buildFile))
	}
}
