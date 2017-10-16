package build

import (
	"context"
	"time"

	"github.com/fudanchii/monocle/errors"

	"github.com/moby/moby/api/types"
	"github.com/moby/moby/client"
)

func Start(buildName string, config *Build) {
	cli, err := client.NewEnvClient()
	errors.ErrCheck(err)

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	cConfig := config.ToMobyContainerConfig()
	hConfig := config.ToMobyHostConfig()
	nConfig := config.ToMobyNetworkConfig()
	rs, err := cli.ContainerCreate(ctx, cConfig, hConfig, nConfig, buildName)
	errors.ErrCheck(err)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logReader, err := cli.ContainerLogs(ctx, rs.ID, types.ContainerLogsOptions{})
	errors.ErrCheck(err)

	_, err := io.Copy(os.Stdout, logReader)
	errors.AssertFalse(err != nil && err != io.EOF, err.Error)
}
