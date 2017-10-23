package build

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
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

func StartDocker(buildName string, config *Build) error {
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	defer cli.Close()

	runner := &RunnerCli{cli, buildName, config}

	return runner.Start()
}

func StartShell(buildName string, config *Build) error {
	var cmd bytes.Buffer

	cfg := config.Shell

	tmp, err := ioutil.TempFile("/tmp", "monocle_shell_cmd")
	if err != nil {
		return err
	}

	cmd.WriteString(CmdHeader)
	cmd.WriteString(cfg.Steps)
	tmp.Write(cmd.Bytes())

	tmp.Close()
	os.Chmod(tmp.Name(), 0755)
	defer os.Remove(tmp.Name())

	output, err := exec.Command(tmp.Name()).CombinedOutput()
	fmt.Println(string(output))
	return err
}

func Start(buildName string, config *Build) error {
	if config.Docker != nil && config.Shell != nil {
		return fmt.Errorf("cannot run both docker and shell")
	}
	switch {
	case config.Docker != nil:
		return StartDocker(buildName, config)
	case config.Shell != nil:
		return StartShell(buildName, config)
	}
	return nil
}
