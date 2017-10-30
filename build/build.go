package build

import (
	"bytes"
	"fmt"
	"io"
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

	runner := &DockerRunner{cli, buildName, config}

	return runner.Start()
}

func runShellCommand(cmdstr string) error {
	var (
		err            error
		stdout, stderr io.ReadCloser
	)

	cmd := exec.Command(cmdstr)
	if stdout, err = cmd.StdoutPipe(); err == nil {
		go io.Copy(os.Stdout, stdout)
		if stderr, err = cmd.StderrPipe(); err == nil {
			go io.Copy(os.Stderr, stderr)
		}
	}
	if err != nil {
		return err
	}

	if err = cmd.Start(); err != nil {
		return err
	}
	return cmd.Wait()
}

func StartShell(buildName string, config *Build) error {
	var (
		err error
		cmd bytes.Buffer
	)

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

	return runShellCommand(tmp.Name())
}

func Start(buildName string, config *Build) error {
	if config.Shell != nil {
		if err := StartShell(buildName, config); err != nil {
			return err
		}
	}

	if config.Docker != nil {
		if err := StartDocker(buildName, config); err != nil {
			return err
		}
	}

	return nil
}
