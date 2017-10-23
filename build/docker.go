package build

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/mount"
	"github.com/docker/docker/api/types/network"
)

const CmdHeader = `#!/bin/sh
set -o errexit

`

type tmpFileCleanFunc func()

func noop() {}

func (b *DockerRunBuild) DockerCmdToShellScript() (tmpFileCleanFunc, error) {
	var cmd bytes.Buffer

	tmp, err := ioutil.TempFile("/tmp", "monocle_steps_cmd")
	if err != nil {
		return noop, err
	}

	cmd.WriteString(CmdHeader)
	cmd.WriteString(b.Steps)
	tmp.Write(cmd.Bytes())

	b.Steps = tmp.Name()
	b.Volumes = append(b.Volumes, b.Steps)

	tmp.Close()
	os.Chmod(b.Steps, 0755)

	return func() { os.Remove(tmp.Name()) }, err
}

func (b *DockerRunBuild) ToDockerContainerConfig() *container.Config {
	cfg := &container.Config{}
	cfg.Image = b.Image
	cfg.WorkingDir = b.Workdir
	cfg.Tty = true
	cfg.Cmd = []string{b.Steps}
	cfg.Env = make([]string, len(b.Env))
	copy(cfg.Env, b.Env)
	return cfg
}

func (b *DockerRunBuild) ToDockerHostConfig() (*container.HostConfig, error) {
	var err error

	hcfg := &container.HostConfig{}
	hcfg.AutoRemove = true
	for _, v := range b.Volumes {
		if mnt, err := splitMounts(v); err == nil {
			hcfg.Mounts = append(hcfg.Mounts, mnt)
		}
	}
	return hcfg, err
}

func (b *DockerRunBuild) ToDockerNetworkingConfig() *network.NetworkingConfig {
	return &network.NetworkingConfig{}
}

func (b *DockerRunBuild) ToDockerClientConfig() (*container.Config, *container.HostConfig, *network.NetworkingConfig, error) {
	hostConfig, err := b.ToDockerHostConfig()
	return b.ToDockerContainerConfig(), hostConfig, b.ToDockerNetworkingConfig(), err
}

func splitMounts(v string) (mount.Mount, error) {
	mnt := mount.Mount{}
	vols := strings.Split(v, ":")

	nv, err := filepath.Abs(vols[0])
	if err != nil {
		return mnt, err
	}

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
	return mnt, err
}

func (b *DockerImageBuild) ToBuildOptions() types.ImageBuildOptions {
	opts := types.ImageBuildOptions{}
	opts.Tags = b.Tags
	opts.Dockerfile = b.File
	opts.AuthConfigs = b.Auths
	return opts
}

func (b *DockerImageBuild) CreateBuildContext() (io.Reader, error) {
	var rawBuff = bytes.NewBuffer([]byte{})

	owd, err := os.Getwd()
	if err != nil {
		return rawBuff, err
	}

	if strings.HasPrefix(b.Root, "/") {
		return rawBuff, fmt.Errorf("root cannot reside out of current working dir: %s", owd)
	}

	if b.Root != "." && b.Root != "" {
		if err = os.Chdir(b.Root); err != nil {
			return rawBuff, err
		}
		defer os.Chdir(owd)
	}

	gz := gzip.NewWriter(rawBuff)
	tw := tar.NewWriter(gz)

	err = filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		link := ""
		if info.Mode()&os.ModeSymlink != 0 {
			if link, err = os.Readlink(path); err != nil {
				return err
			}
		}

		fheader, err := tar.FileInfoHeader(info, link)
		if err != nil {
			return err
		}

		fheader.Name = path

		if err = tw.WriteHeader(fheader); err != nil {
			return err
		}

		if !info.IsDir() {
			if f, err := os.Open(path); err == nil {
				_, err = io.Copy(tw, f)
				f.Close()
			}
		}

		return err
	})

	err = tw.Close()
	err = gz.Close()

	return rawBuff, err
}
