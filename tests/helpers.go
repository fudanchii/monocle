package tests

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

type cleanFunc func()

func noop() {}

const (
	gitUser  = "user.name=monocle"
	gitEmail = "user.email=monocle@monocletest.com"
)

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

type fileRep struct {
	name    string
	content []byte
}

func createCommit(dir string, files []fileRep, cmsg string, err error) error {
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

	for i := 0; i < len(files); i++ {
		if err = ioutil.WriteFile(files[i].name, files[i].content, 0644); err != nil {
			return err
		}
	}

	if output, err := exec.Command("git", "add", ".").CombinedOutput(); err != nil {
		return fmt.Errorf("err: %s\nerr: %s", output, err.Error())
	}

	if output, err := exec.Command("git", "-c", gitUser, "-c", gitEmail, "commit", "-am", cmsg).CombinedOutput(); err != nil {
		return fmt.Errorf("err: %s\nerr: %s", output, err.Error())
	}

	return nil
}

func seedSimpleCommit(dir string, err error) error {
	f := []fileRep{fileRep{"a", []byte{'a', 'b', 'c', '\n'}}}
	return createCommit(dir, f, "First commit!", err)
}

func seedAnotherCommit(dir string, err error) error {
	f := []fileRep{
		fileRep{"b", []byte{'d', 'e', 'f', '\n'}},
		fileRep{"c", []byte{'g', 'h', 'i', '\n'}},
	}
	return createCommit(dir, f, "Second commit!", err)
}
