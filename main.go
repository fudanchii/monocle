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
		files   git.Files
		err     error
	)

	flag.Parse()
	switch flag.NArg() {
	case 2:
		rev = flag.Arg(1)
		fallthrough
	case 1:
		workDir = flag.Arg(0)
	}

	if files, err = git.FilesChanged(workDir, rev); err == nil {
		for _, buildFile := range build.CreateBuildEntries(files) {
			var manifest *build.Build
			manifest, err = build.ParseManifest(buildFile)
			if err != nil {
				break
			}
			if err = build.Start(build.Name(buildFile), manifest); err != nil {
				break
			}
		}
	}

	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(-1)
	}
}
