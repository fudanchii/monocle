package build

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/fudanchii/monocle/errors"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"

	uuid "github.com/satori/go.uuid"
)

func Name(buildFile string) string {
	bFileLoc, err := filepath.Abs(buildFile)
	errors.ErrCheck(err)

	name := path.Base(path.Dir(bFileLoc))
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
	errors.ErrCheck(err)

	err = cli.ContainerStart(ctx, rs.ID, types.ContainerStartOptions{})
	errors.ErrCheck(err)

	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	scodeChan, errChan := cli.ContainerWait(ctx, rs.ID, container.WaitConditionNextExit)

	reader, err := cli.ContainerLogs(ctx, rs.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	})
	errors.ErrCheck(err)

	_, err = io.Copy(os.Stdout, reader)
	if !(err == nil || err == io.EOF) {
		fmt.Println("err: ", err.Error())
		os.Exit(-1)
	}

	select {
	case _ = <-scodeChan:
		return
	case errCode := <-errChan:
		errors.ErrCheck(errCode)
	}
}

func startDockerBuild(cli *client.Client, buildName string, config *DockerImageBuild) {
}
