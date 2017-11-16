package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

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
		buildName string
		workDir   string = "."
		rev       string = "HEAD"
		files     git.Files
		err       error
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
		var (
			wkdir, cdir string
		)

		if wkdir, err = filepath.Abs(workDir); err != nil {
			goto BailOut
		}
		if cdir, err = os.Getwd(); err != nil {
			goto BailOut
		}
		if wkdir != cdir {
			if err = os.Chdir(wkdir); err != nil {
				goto BailOut
			}
		}
		for _, buildFile := range build.CreateBuildEntries(files) {
			var manifest *build.Build
			manifest, err = build.ParseManifest(buildFile)
			if err != nil {
				break
			}

			buildName, err = build.Name(buildFile)
			if err != nil {
				break
			}

			if err = build.Start(buildName, manifest); err != nil {
				break
			}
		}
	}

BailOut:
	if err != nil {
		fmt.Printf("%s\n", err.Error())
		os.Exit(-1)
	}
}
