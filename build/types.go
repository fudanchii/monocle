package build

import (
	"github.com/docker/docker/api/types"
)

type DockerImageBuild struct {
	File  string                      `yaml:"file"`
	Root  string                      `yaml:"root"`
	Tags  []string                    `yaml:"tags"`
	Auths map[string]types.AuthConfig `yam:"auths"`
	Push  *DockerAuthConfig           `yaml:"push"`
}

type DockerRunBuild struct {
	Image    string            `yaml:"image"`
	Steps    string            `yaml:"steps"`
	Workdir  string            `yaml:"workdir"`
	Env      []string          `yaml:"env"`
	Volumes  []string          `yaml:"volumes"`
	Services []DockerServices  `yaml:"services"`
	Auth     *DockerAuthConfig `yaml:"auth"`
}

type DockerServices struct {
	Image string `yaml:"image"`
	Name  string `yaml:"name"`
}

type DockerAuthConfig struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Email    string `yaml:"email"`
	Registry string `yaml:"registry"`
}

type DockerBuild struct {
	Run   *DockerRunBuild   `yaml:"run"`
	Build *DockerImageBuild `yaml:"build"`
}

type Build struct {
	Docker *DockerBuild `yaml:"docker"`
	Shell  *ShellBuild  `yaml:"shell"`
}

type ShellBuild struct {
	Steps string `yaml:"steps"`
}
