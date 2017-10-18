package build

type DockerImageBuild struct {
	File string     `yaml:"file"`
	Repo string     `yaml:"repo"`
	Tag  string     `yaml:"tag"`
	Push DockerPush `yaml:"push"`
}

type DockerRunBuild struct {
	Image    string           `yaml:"image"`
	Steps    string           `yaml:"steps"`
	Workdir  string           `yaml:"workdir"`
	Env      []string         `yaml:"env"`
	Volumes  []string         `yaml:"volumes"`
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
	Run   *DockerRunBuild   `yaml:"run"`
	Build *DockerImageBuild `yaml:"build"`
}

type Build struct {
	Docker DockerBuild `yaml:"docker"`
}
