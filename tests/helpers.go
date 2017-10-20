package tests

import (
	"fmt"
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

	if output, err := exec.Command("git", "init", dir).CombinedOutput(); err != nil {
		return "", noop, fmt.Errorf("err: %s\nerr: %s", output, err.Error())
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

	if output, err := exec.Command("git", "add", ".").CombinedOutput(); err != nil {
		return fmt.Errorf("err: %s\nerr: %s", output, err.Error())
	}

	if output, err := exec.Command("git", "commit", "-am", "First commit!").CombinedOutput(); err != nil {
		return fmt.Errorf("err: %s\nerr: %s", output, err.Error())
	}

	return nil
}
