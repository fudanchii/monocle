package build

import (
	"fmt"
	"path"
	"path/filepath"

	"github.com/fudanchii/monocle/errors"

	"github.com/docker/docker/client"

	uuid "github.com/satori/go.uuid"
)

func Name(buildFile string) string {
	bFileLoc, err := filepath.Abs(buildFile)
	errors.ErrCheck(err)

	name := path.Base(path.Dir(bFileLoc))
	return fmt.Sprintf("%s-%s", name, uuid.NewV5(uuid.NewV1(), buildFile))
}

func Start(buildName string, config *Build) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	runner := &RunnerCli{cli, buildName, config}

	return runner.Start()
}
