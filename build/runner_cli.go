package build

import (
	"context"
	"encoding/base64"
	"encoding/json"
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
		if err != nil && err != io.EOF {
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
	config := rc.Config.Docker.Build
	buildOpts := config.ToBuildOptions()
	ctxReader, err := config.CreateBuildContext()
	if err != nil {
		return &DockerBuildError{err}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Minute)
	defer cancel()

	fmt.Println()
	fmt.Println("----- Building image -----")
	rsp, err := rc.Cli.ImageBuild(ctx, ctxReader, buildOpts)
	err = parseBuildResponseStream(rsp.Body, err)
	err = rc.PushDockerImage(err)

	if err != nil {
		return &DockerBuildError{err}
	}

	return nil
}

type authConfig struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Email         string `json:"email",omitempty`
	Serveraddress string `json:"serveraddress"`
}

func (rc *RunnerCli) PushDockerImage(err error) error {
	var regAuth authConfig

	if err != nil {
		return err
	}

	config := rc.Config.Docker.Build.Push
	if config == nil {
		return nil
	}

	fmt.Println()
	fmt.Println("----- Pushing image -----")

	regAuth.Username = config.User
	regAuth.Password = config.Password
	regAuth.Email = config.Email
	regAuth.Serveraddress = config.Registry

	authConfigFromEnv(&regAuth)

	regJson, err := json.Marshal(regAuth)
	if err != nil {
		return err
	}

	regb64 := base64.StdEncoding.EncodeToString(regJson)
	pushOpts := types.ImagePushOptions{
		RegistryAuth: regb64,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Minute)
	defer cancel()

	image := rc.Config.Docker.Build.Tags[0]
	outr, err := rc.Cli.ImagePush(ctx, image, pushOpts)
	err = parsePushResponseStream(outr, err)

	return err
}

func authConfigFromEnv(reg *authConfig) {
	if u, ok := os.LookupEnv("DOCKER_USERNAME"); ok && u != "" {
		reg.Username = u
	}
	if p, ok := os.LookupEnv("DOCKER_PASSWORD"); ok && p != "" {
		reg.Password = p
	}
	if e, ok := os.LookupEnv("DOCKER_EMAIL"); ok && e != "" {
		reg.Email = e
	}
}

func parseBuildResponseStream(in io.Reader, err error) error {
	var (
		msg map[string]interface{} = make(map[string]interface{})
		dec *json.Decoder          = json.NewDecoder(in)
	)

	if err != nil {
		return err
	}

	for err = nil; err != io.EOF; {
		if err = dec.Decode(&msg); err != nil && err != io.EOF {
			return err
		}
		if msg["stream"] != nil {
			fmt.Print(msg["stream"].(string))
		}
		if msg["aux"] != nil {
			for k, v := range msg["aux"].(map[string]interface{}) {
				fmt.Printf("%s: %v\n", k, v)
			}
		}
	}

	return nil
}

func parsePushResponseStream(in io.Reader, err error) error {
	var (
		msg map[string]interface{} = make(map[string]interface{})
		dec *json.Decoder          = json.NewDecoder(in)
	)

	if err != nil {
		return err
	}

	for err = nil; err != io.EOF; {
		if err = dec.Decode(&msg); err != nil && err != io.EOF {
			return err
		}

		if msg["id"] != nil {
			fmt.Print(msg["id"].(string))
			fmt.Print(": ")
		}
		if msg["status"] != nil {
			fmt.Println(msg["status"].(string))
		}
		if msg["error"] != nil {
			return fmt.Errorf("push err: %s", msg["error"].(string))
		}
	}
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
