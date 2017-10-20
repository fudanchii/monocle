package tests

import (
	"io/ioutil"
	"os"
	"os/exec"
)

type cleanFunc func()

func noop() {}

func createFixtureRepo() (string, cleanFunc, error) {
	dir, err := ioutil.TempDir("/tmp", "monocle_test_")
	if err != nil {
		return "", noop, err
	}

	if err = exec.Command("git", "init", dir).Run(); err != nil {
		return "", noop, err
	}

	return dir, func() {
		exec.Command("rm", "-rf", dir).Run()
	}, nil
}

func seedSimpleCommit(dir string, err error) error {
	if err != nil {
		return err
	}

	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	defer os.Chdir(cwd)

	if err = os.Chdir(dir); err != nil {
		return err
	}

	if err = ioutil.WriteFile("a", []byte{'a', 'b', 'c', '\n'}, 0644); err != nil {
		return err
	}

	if err = exec.Command("git", "add", ".").Run(); err != nil {
		return err
	}

	if err = exec.Command("git", "commit", "-am", "First commit!").Run(); err != nil {
		return err
	}

	return nil
}
