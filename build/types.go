package build

import (
	"github.com/docker/docker/api/types"
)

type DockerImageBuild struct {
	File  string                      `yaml:"file"`
	Root  string                      `yaml:"root"`
	Tags  []string                    `yaml:"tags" templatable:""`
	Auths map[string]types.AuthConfig `yaml:"auths"`
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
	Image string            `yaml:"image"`
	Name  string            `yaml:"name"`
	Auth  *DockerAuthConfig `yaml:"auth"`
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
	*Variables `yaml:"variables"`
	Docker     *DockerBuild `yaml:"docker"`
	Shell      *ShellBuild  `yaml:"shell"`
}

type ShellBuild struct {
	Steps string `yaml:"steps"`
}

type Variables struct {
	Eval map[string]string `yaml:"eval"`
	Env  map[string]string `yaml:"env"`
}
