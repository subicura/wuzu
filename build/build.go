package build

import (
	"bytes"
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/fsouza/go-dockerclient"
	"github.com/subicura/wuzu/config"
	"os"
)

var client *docker.Client

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Build() {
	log.Infoln("build start")

	if _, err := os.Stat(config.DefaultFilename); os.IsNotExist(err) {
		log.Errorln(".wuzu.yml file is required.")
		os.Exit(1)
	}

	wuzuConfig, err := config.Parse(config.DefaultFilename)
	checkErr(err)

	endpoint := os.Getenv("DOCKER_HOST")
	if endpoint == "" {
		endpoint = "unix:///var/run/docker.sock"
	}
	client, err = docker.NewClient(endpoint)
	checkErr(err)

	log.Infoln(fmt.Sprintf("pull image: %s", wuzuConfig.Build.From))
	err = pullImage(&wuzuConfig.Build)
	checkErr(err)

	log.Infoln(fmt.Sprintf("run build command: %s", wuzuConfig.Build.Run))
	err = run(&wuzuConfig.Build)

	log.Infoln("build complete")
}

func pullImage(c *config.BuildConfig) error {
	opts := docker.PullImageOptions{Repository: c.From}
	return client.PullImage(opts, docker.AuthConfiguration{})
}

func run(c *config.BuildConfig) error {
	conf := &docker.Config{
		Image:      c.From,
		WorkingDir: c.Dest,
		Cmd:        []string{"/bin/sh", "-c", c.Run},
	}
	hostConf := &docker.HostConfig{
		Binds: []string{c.Src + ":" + c.Dest},
	}
	containerConf := docker.CreateContainerOptions{
		Config:     conf,
		HostConfig: hostConf,
	}

	// create container
	container, err := client.CreateContainer(containerConf)
	checkErr(err)

	defer func() {
		// remove container
		err := client.RemoveContainer(docker.RemoveContainerOptions{
			ID: container.ID,
		})
		checkErr(err)
	}()

	// start container
	err = client.StartContainer(container.ID, hostConf)
	checkErr(err)

	// attach container
	var buf bytes.Buffer
	err = client.AttachToContainer(docker.AttachToContainerOptions{
		Container:    container.ID,
		OutputStream: &buf,
		ErrorStream:  &buf,
		Stream:       true,
		Logs:         true,
		Stdout:       true,
		Stderr:       true,
	})
	checkErr(err)
	log.Info(buf.String())

	return nil
}
