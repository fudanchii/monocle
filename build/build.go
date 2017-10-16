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

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	cConfig := config.ToMobyContainerConfig()
	hConfig := config.ToMobyHostConfig()
	nConfig := config.ToMobyNetworkConfig()
	rs, err := cli.ContainerCreate(ctx, cConfig, hConfig, nConfig, buildName)
	errors.ErrCheck(err)

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logReader, err := cli.ContainerLogs(ctx, rs.ID, types.ContainerLogsOptions{})
	errors.ErrCheck(err)

	_, err = io.Copy(os.Stdout, logReader)
	errors.AssertFalse(err != nil && err != io.EOF, err.Error())
}
