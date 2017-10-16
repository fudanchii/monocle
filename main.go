package main

import (
	"flag"
	"fmt"

	"github.com/fudanchii/monocle/build"
	"github.com/fudanchii/monocle/git"
)

func main() {
	workDir := "."
	flag.Parse()
	if flag.NArg() == 1 {
		workDir = flag.Arg(0)
	}

	buildFiles := build.CreateBuildEntries(git.FilesChanged(workDir, ""))
	for _, buildFile := range buildFiles {
		buildConfig := build.ParseManifest(buildFile)
		fmt.Printf("%#v\n", buildConfig)
		build.Start(buildConfig)
	}
}
