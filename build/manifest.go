package build

import (
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/fudanchii/monocle/errors"
	"github.com/fudanchii/monocle/set"
	"gopkg.in/yaml.v2"
)

const (
	BuildFile = "build.yml"
)

type DockerImageBuild struct {
	File string     `yaml:"file"`
	Repo string     `yaml:"repo"`
	Tag  string     `yaml:"tag"`
	Push DockerPush `yaml:"push"`
}

type DockerPush struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Registry string `yaml:"registry"`
}

type DockerBuild struct {
	Image string           `yaml:"image"`
	Build DockerImageBuild `yaml:"build"`
}

type Build struct {
	Docker []DockerBuild `yaml:"docker"`
}

func ParseManifest(buildFile string) *Build {
	var buildConfig Build

	content, err := ioutil.ReadFile(buildFile)
	errors.ErrCheck(err)

	err = yaml.Unmarshal([]byte(content), &buildConfig)
	errors.ErrCheck(err)

	return &buildConfig
}

func CreateBuildEntries(changes []string) []string {
	var buildEntries = set.NewSet()
	for _, entry := range changes {
		if subdir, hasManifestFile := searchBuildManifest(entry); hasManifestFile {
			buildEntries.Add(subdir)
		}
	}
	return buildEntries.ToStringSlice()
}

func searchBuildManifest(item string) (string, bool) {
	pathFragment, _ := path.Split(item)
	for _, fr := range strings.Split(pathFragment, "/") {
		bfile := path.Join(fr, BuildFile)
		if _, err := os.Stat(bfile); err == nil {
			return bfile, true
		}
	}
	return "", false
}
