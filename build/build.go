package build

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"time"

	"github.com/fudanchii/monocle/errors"

	"github.com/docker/docker/api/types"
	"github.com/moby/moby/client"
	"github.com/satori/go.uuid"
)

func Name(buildFile string) string {
	name := path.Base(path.Dir(buildFile))
	return fmt.Sprintf("%s-%s", name, uuid.NewV5(uuid.NewV1(), buildFile))
}

func Start(buildName string, config *Build) {
	cli, err := client.NewEnvClient()
	errors.ErrCheck(err)

	if config.Docker.Run != nil {
		startDockerRun(cli, buildName, config.Docker.Run)
	}
	if config.Docker.Build != nil {
		startDockerBuild(cli, buildName, config.Docker.Build)
	}
}

func startDockerRun(cli *client.Client, buildName string, config *DockerRunBuild) {
	cleanup := config.DockerCmdToShellScript()
	cConfig, hConfig, nConfig := config.ToDockerClientConfig()
	defer cleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	rs, err := cli.ContainerCreate(ctx, cConfig, hConfig, nConfig, buildName)
	defer cleanUpContainer(rs.ID)
	errors.ErrCheck(err)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	reader, err := cli.ContainerLogs(ctx, rs.ID, types.ContainerLogsOptions{})
	errors.ErrCheck(err)

	_, err = io.Copy(os.Stdout, reader)
	errors.Assert(err == nil || err == io.EOF, err.Error())
}

func startDockerBuild(cli *client.Client, buildName string, config *DockerImageBuild) {
}
