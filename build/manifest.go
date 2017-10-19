package build

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/fudanchii/monocle/errors"
	"github.com/fudanchii/monocle/git"
	"github.com/fudanchii/monocle/set"
	"gopkg.in/yaml.v2"
)

const (
	BuildFile = "build.yml"
)

func ParseManifest(buildFile string) *Build {
	var buildConfig Build

	content, err := ioutil.ReadFile(buildFile)
	errors.ErrCheck(err)

	err = yaml.Unmarshal([]byte(content), &buildConfig)
	errors.ErrCheck(err)

	return &buildConfig
}

func CreateBuildEntries(changes git.Files) []string {
	var buildEntries = set.NewSet()
	for _, entry := range changes.Entries {
		if subdir, ok := searchBuildManifest(changes.WorkDir, entry); ok {
			buildEntries.Add(subdir)
		}
	}
	return buildEntries.ToStringSlice()
}

func searchBuildManifest(workdir, item string) (string, bool) {
	var pathFragment string
	for _, fr := range strings.Split(path.Dir(item), "/") {
		pathFragment = path.Join(pathFragment, fr)
		bfile := path.Join(workdir, pathFragment, BuildFile)
		if _, err := os.Stat(bfile); err == nil {
			return bfile, true
		}
	}
	return "", false
}
