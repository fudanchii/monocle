package build

import (
	"io/ioutil"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/fudanchii/monocle/errors"
)

const CmdHeader = `#! /usr/bin/env sh
set -o errexit
set -o pipefail

`

type tmpFileCleanFunc func()

func (b *DockerRunBuild) DockerCmdToShellScript() tmpFileCleanFunc {
	var cmd bytes.Buffer
	tmp, err := ioutil.TempFile("/tmp", "monocle_steps_cmd")
	errors.ErrCheck(err)

	buffer.WriteString(CmdHeader)
	buffer.WriteString(b.Steps)
	tmp.Write(buffer.String())

	b.Steps = tmp.Name()
	append(b.Env, b.Steps)

	tmp.Close()
	os.Chmod(b.Steps, 0755)

	return func() {
		os.Remove(tmp.Name)
	}
}

func (b *DockerRunBuild) ToDockerContainerConfig() *container.Config {
	cfg := &container.Config{}
	cfg.Image = b.Image
	cfg.WorkingDir = b.Workdir
	cfg.Tty = true
	cfg.Cmd = b.Steps
	cfg.Env = make([]string, len(b.Env))
	copy(cfg.Env, b.Env)
	return cfg
}

func (b *DockerRunBuild) ToDockerHostConfig() *container.HostConfig {
	hcfg := &container.HostConfig{}
	hcfg.AutoRemove = true
	for _, v := range b.Volumes {
		mnt := splitMounts(v)
		hcfg.Mounts = append(hcfg.Mounts, mnt)
	}
	return hcfg
}

func (b *DockerRunBuild) ToDockerNetworkConfig() *network.NetworkingConfig {
	return &network.NetworkingConfig{}
}

func (b *DockerRunBuild) ToDockerClientConfig() (*container.Config, *container.HostConfig, *network.NetworkingConfig) {
	return b.ToDockerContainerConfig(), b.ToDockerHostConfig(), b.ToDockerNetworkingConfig()
}

func splitMounts(v string) mount.Mount {
	mnt := mount.Mount{}
	vols := strings.Split(v, ":")
	switch len(vols) {
	case 1:
		mnt.Source = vols[0]
		mnt.Target = vols[0]
	case 2:
		mnt.Source = vols[0]
		mnt.Target = vols[1]
	default:
		mnt.Source = vols[0]
		mnt.Source = vols[1]
		if vols[2] == "ro" {
			mnt.ReadOnly = true
		}
	}
	return mnt
}
