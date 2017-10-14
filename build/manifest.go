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

type DockerImageBuild struct {
	File string     `yaml:"file"`
	Repo string     `yaml:"repo"`
	Tag  string     `yaml:"tag"`
	Push DockerPush `yaml:"push"`
}

type DockerRunBuild struct {
	Image    string           `yaml:"image"`
	Steps    string           `yaml:"steps"`
	Services []DockerServices `yaml:"services"`
}

type DockerServices struct {
	Image string `yaml:"image"`
	Name  string `yaml:"name"`
}

type DockerPush struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Registry string `yaml:"registry"`
}

type DockerBuild struct {
	Run   DockerRunBuild   `yaml:"run"`
	Build DockerImageBuild `yaml:"build"`
}

type Build struct {
	Docker []DockerBuild `yaml:"docker"`
	Envs   []string      `yaml:"envs"`
}

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
		if subdir, hasManifestFile := searchBuildManifest(changes.WorkDir, entry); hasManifestFile {
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
