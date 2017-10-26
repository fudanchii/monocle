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

const (
	reqTimeout     = 30 * time.Second
	pullTimeout    = 30 * time.Minute
	pushTimeout    = 30 * time.Minute
	runWaitTimeout = 1 * time.Hour
	buildTimeout   = 30 * time.Minute
)

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
	var (
		pullOpts types.ImagePullOptions
		err      error
		rsReader io.ReadCloser
	)

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

	retryCount := 0
CreateContainer:
	ctx, cancel := context.WithTimeout(context.Background(), reqTimeout)
	defer cancel()

	rs, err := rc.Cli.ContainerCreate(ctx, cConfig, hConfig, nConfig, rc.Name)
	if err == nil {
		err = rc.Cli.ContainerStart(ctx, rs.ID, types.ContainerStartOptions{})
	}

	if err != nil {
		if client.IsErrNotFound(err) { // docker wont tell us what is not found, blindly believe it's the image
			if retryCount > 0 {
				goto BailOut
			}

			ctx, cancel := context.WithTimeout(context.Background(), pullTimeout)
			defer cancel()

			if pullOpts, err = rc.imagePullOptions(config.Auth); err == nil {
				fmt.Printf("Pulling image: %s\n", config.Image)
				rsReader, err = rc.Cli.ImagePull(ctx, config.Image, pullOpts)
				if err != nil {
					goto BailOut
				}
				if err = parsePullPushResponseStream(rsReader, err); err != nil {
					goto BailOut
				}
				retryCount = retryCount + 1
				goto CreateContainer
			}
		}
	BailOut:
		return &DockerRunError{err}
	}

	ctx, cancel = context.WithTimeout(context.Background(), runWaitTimeout)
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

	ctx, cancel := context.WithTimeout(context.Background(), buildTimeout)
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

func (rc *RunnerCli) PushDockerImage(err error) error {
	if err != nil {
		return err
	}

	config := rc.Config.Docker.Build.Push
	if config == nil {
		return nil
	}

	fmt.Println()
	fmt.Println("----- Pushing image -----")

	pullOpts, err := rc.imagePullOptions(config)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), pushTimeout)
	defer cancel()

	image := rc.Config.Docker.Build.Tags[0]
	outr, err := rc.Cli.ImagePush(ctx, image, types.ImagePushOptions(pullOpts))
	err = parsePullPushResponseStream(outr, err)

	return err
}

type authConfig struct {
	Username      string `json:"username"`
	Password      string `json:"password"`
	Email         string `json:"email",omitempty`
	Serveraddress string `json:"serveraddress"`
}

func (rc *RunnerCli) imagePullOptions(config *DockerAuthConfig) (types.ImagePullOptions, error) {
	var (
		regAuth authConfig
		result  types.ImagePullOptions
	)

	if config == nil {
		return result, nil
	}

	regAuth.Username = config.User
	regAuth.Password = config.Password
	regAuth.Email = config.Email
	regAuth.Serveraddress = config.Registry

	authConfigFromEnv(&regAuth)

	regJson, err := json.Marshal(regAuth)
	if err != nil {
		return result, err
	}

	regb64 := base64.StdEncoding.EncodeToString(regJson)
	result.RegistryAuth = regb64
	return result, nil
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

func parsePullPushResponseStream(in io.Reader, err error) error {
	var (
		msg               map[string]interface{} = make(map[string]interface{})
		dec               *json.Decoder          = json.NewDecoder(in)
		currentID, nextID string
	)

	if err != nil {
		return err
	}

	for err = nil; err != io.EOF; {
		if err = dec.Decode(&msg); err != nil && err != io.EOF {
			return err
		}

		if msg["id"] != nil {
			nextID, _ = msg["id"].(string)
			if nextID == currentID {
				fmt.Print("\r\033[K")
			} else {
				currentID = nextID
				fmt.Println()
			}
			fmt.Print(nextID)
			fmt.Print(": ")
		}
		if msg["status"] != nil {
			fmt.Printf("%s ", msg["status"].(string))
		}
		if msg["progress"] != nil {
			fmt.Print(msg["progress"].(string))
		}
		if msg["error"] != nil {
			fmt.Println()
			return fmt.Errorf("push err: %s", msg["error"].(string))
		}
	}
	fmt.Println()
	return nil
}

type DockerRunError struct {
	err error
}

func (dr *DockerRunError) Error() string {
	return fmt.Sprintf("docker run error, caused by: [%T] %s", dr.err, dr.err.Error())
}

type DockerBuildError struct {
	err error
}

func (db *DockerBuildError) Error() string {
	return fmt.Sprintf("docker build error, caused by: [%T] %s", db.err, db.err.Error())
}
