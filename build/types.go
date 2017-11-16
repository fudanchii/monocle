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
	Image    string              `yaml:"image"`
	Steps    string              `yaml:"steps"`
	Workdir  string              `yaml:"workdir"`
	Env      []string            `yaml:"env"`
	Volumes  []map[string]string `yaml:"volumes"`
	Services []DockerServices    `yaml:"services"`
	Auth     *DockerAuthConfig   `yaml:"auth"`
}

type DockerServices struct {
	Image     string                  `yaml:"image"`
	Name      string                  `yaml:"name"`
	Command   string                  `yaml:"command"`
	Env       []string                `yaml:"env"`
	Auth      *DockerAuthConfig       `yaml:"auth"`
	Ports     []string                `yaml:"port"`
	Volumes   []map[string]string     `yaml:"volumes"`
	Readiness *DockerServiceReadiness `yaml:"readiness"`
	DependsOn []string                `yaml:"depends_on"`
}

type DockerServiceReadiness struct {
	HTTP           *DockerServiceHttpReady `yaml:"http"`
	Cmd            *DockerServiceCmdReady  `yaml:"cmd"`
	TCP            *DockerServiceTcpReady  `yaml:"tcp"`
	UDP            *DockerServiceUdpReady  `yaml:"udp"`
	InitialWait    int                     `yaml:"initial_wait"`
	ErrorThreshold int                     `yaml:"error_threshold"`
	OkThreshold    int                     `yaml:"ok_threshold"`
}

type DockerServiceHttpReady struct {
	Endpoint       string `yaml:"endpoint"`
	ExpectedStatus []int  `yaml:"expected_status"`
}

type DockerServiceCmdReady struct {
	Command string `yaml:"command"`
}

type DockerServiceTcpReady struct {
	Port             int    `yaml:"port"`
	Send             string `yaml:"send"`
	ExpectedResponse string `yaml:"ecpected_response"`
}

type DockerServiceUdpReady DockerServiceTcpReady

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
