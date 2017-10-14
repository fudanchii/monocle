package main

import (
	"flag"
	"fmt"

	"github.com/fudanchii/monocle/build"
	"github.com/fudanchii/monocle/git"
)

var (
	buildConfig build.Build
	buildFile   = "build.yml"
)

func main() {
	workDir := "."
	flag.Parse()
	if flag.NArg() == 1 {
		workDir = flag.Arg(0)
	}

	changedFiles := git.FilesChanged(workDir, "")
	fmt.Printf("%#v\n", changedFiles)

	buildFiles := build.CreateBuildEntries(changedFiles)
	fmt.Println(buildFiles)
}
