package build

import (
	"context"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type RunnerCli struct {
	Cli    *client.Client
	Name   string
	Config *Build
}

func (rc *RunnerCli) Start() error {
	var err error

	if rc.Config.Docker.Run != nil {
		err = rc.StartDockerRun()
	}
	if err == nil && rc.Config.Docker.Build != nil {
		err = rc.StartDockerBuild()
	}
	return err
}

func (rc *RunnerCli) StartDockerRun() error {
	config := rc.Config.Docker.Run
	cleanup, err := config.DockerCmdToShellScript()
	defer cleanup()
	if err != nil {
		return &DockerRunError{err}
	}

	cConfig, hConfig, nConfig, err := config.ToDockerClientConfig()
	if err != nil {
		return &DockerRunError{err}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	rs, err := rc.Cli.ContainerCreate(ctx, cConfig, hConfig, nConfig, rc.Name)
	if err == nil {
		err = rc.Cli.ContainerStart(ctx, rs.ID, types.ContainerStartOptions{})
	}
	if err != nil {
		return &DockerRunError{err}
	}

	ctx, cancel = context.WithTimeout(context.Background(), 1*time.Hour)
	defer cancel()

	scodeChan, errChan := rc.Cli.ContainerWait(ctx, rs.ID, container.WaitConditionNextExit)

	if reader, err := rc.Cli.ContainerLogs(ctx, rs.ID, types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
	}); err == nil {
		_, err = io.Copy(os.Stdout, reader)
		if !(err == nil || err == io.EOF) {
			return &DockerRunError{err}
		}
	}

	select {
	case rs := <-scodeChan:
		if rs.StatusCode == 0 {
			return nil
		}
		return &DockerRunError{fmt.Errorf("container exit code: %d\n", rs.StatusCode)}
	case err = <-errChan:
		return &DockerRunError{err}
	}
}

func (rc *RunnerCli) StartDockerBuild() error {
	return nil
}

type DockerRunError struct {
	err error
}

func (dr *DockerRunError) Error() string {
	return fmt.Sprintf("docker run error, cause: %s", dr.err.Error())
}

type DockerBuildError struct {
	err error
}

func (db *DockerBuildError) Error() string {
	return fmt.Sprintf("docker build error, cause: %s", db.err.Error())
}
