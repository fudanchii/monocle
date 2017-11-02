package build

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/fudanchii/monocle/git"
	"github.com/fudanchii/monocle/set"
	"gopkg.in/yaml.v2"
)

const (
	// BuildFile is a constant for the default manifest file
	BuildFile = "build.yml"
)

// ParseManifest parse given build.yml file and will return
// pointer to Build object and error if any.
// ParseManifest will also extrapolate any string property with `templatable`
// tag, if Build.Variables is populated.
func ParseManifest(buildFile string) (*Build, error) {
	var buildConfig Build

	content, err := ioutil.ReadFile(buildFile)
	if err == nil {
		err = yaml.Unmarshal([]byte(content), &buildConfig)
	}

	return evalVars(&buildConfig, err)
}

// CreateBuildEntries will iterate over git changes entries and create
// a set of build manifest related to that change, build manifest were found
// by calling searchBuildManifest to check if there is BuildFile exist in that folder.
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
	var pathFragment = item
	for {
		pathFragment = path.Dir(pathFragment)
		bfile := path.Join(workdir, pathFragment, BuildFile)
		if _, err := os.Stat(bfile); err == nil {
			return bfile, true
		}
		if pathFragment == "." {
			break
		}
	}
	return "", false
}
