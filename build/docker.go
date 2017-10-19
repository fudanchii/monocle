package build

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
	"github.com/fudanchii/monocle/errors"
)

const CmdHeader = `#!/bin/sh
set -o errexit

`

type tmpFileCleanFunc func()

func (b *DockerRunBuild) DockerCmdToShellScript() tmpFileCleanFunc {
	var cmd bytes.Buffer
	tmp, err := ioutil.TempFile("/tmp", "monocle_steps_cmd")
	errors.ErrCheck(err)

	cmd.WriteString(CmdHeader)
	cmd.WriteString(b.Steps)
	tmp.Write(cmd.Bytes())

	b.Steps = tmp.Name()
	b.Volumes = append(b.Volumes, b.Steps)

	tmp.Close()
	os.Chmod(b.Steps, 0755)

	return func() {
		os.Remove(tmp.Name())
	}
}

func (b *DockerRunBuild) ToDockerContainerConfig() *container.Config {
	cfg := &container.Config{}
	cfg.Image = b.Image
	cfg.WorkingDir = b.Workdir
	cfg.Tty = true
	cfg.Cmd = []string{b.Steps}
	cfg.Env = make([]string, len(b.Env))
	copy(cfg.Env, b.Env)
	fmt.Println(cfg.Env)
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

func (b *DockerRunBuild) ToDockerNetworkingConfig() *network.NetworkingConfig {
	return &network.NetworkingConfig{}
}

func (b *DockerRunBuild) ToDockerClientConfig() (*container.Config, *container.HostConfig, *network.NetworkingConfig) {
	return b.ToDockerContainerConfig(), b.ToDockerHostConfig(), b.ToDockerNetworkingConfig()
}

func splitMounts(v string) mount.Mount {
	mnt := mount.Mount{}
	vols := strings.Split(v, ":")

	nv, err := filepath.Abs(vols[0])
	errors.ErrCheck(err)

	vols[0] = nv
	switch len(vols) {
	case 1:
		mnt.Source = vols[0]
		mnt.Target = vols[0]
	case 2:
		mnt.Source = vols[0]
		mnt.Target = vols[1]
	default:
		mnt.Source = vols[0]
		mnt.Target = vols[1]
		if vols[2] == "ro" {
			mnt.ReadOnly = true
		}
	}
	mnt.Type = mount.TypeBind
	return mnt
}
