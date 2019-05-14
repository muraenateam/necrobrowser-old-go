package core

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"

	"github.com/muraenateam/necrobrowser/log"
)

const (
	ChromeRemotePort    = "9222"
	ListeningIP         = "127.0.0.1"
	DockerClientVersion = "1.39"
)

type Container struct {
	Context context.Context
	Cancel  context.CancelFunc

	client      *client.Client
	Name        string
	Image       string
	LocalPort   string
	PublicPort  string
	IP          string
	Target      ChromeTarget
	MounthPaths []string
}

type ChromeTarget struct {
	Description          string `json:"description"`
	DevtoolsFrontendUrl  string `json:"devtoolsFrontendUrl"`
	ID                   string `json:"id"`
	Title                string `json:"title"`
	Type                 string `json:"type"`
	URL                  string `json:"url"`
	WebSocketDebuggerUrl string `json:"webSocketDebuggerUrl"`
}

func InitDocker(dockerImage string) error {
	// Load Docker image before starting
	return PullImage(dockerImage)
}

func NewContainer(name string, imageName string, exposedPort string, loot string) (d *Container, err error) {

	d = &Container{
		Name:       name,
		Image:      imageName,
		LocalPort:  ChromeRemotePort,
		PublicPort: exposedPort,
		IP:         ListeningIP,
	}

	// Make mount paths
	d.MounthPaths = []string{fmt.Sprintf("%s:/tmp", loot)}

	// Load the environment
	err = d.initEnvironment()
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Let's start a new container
	if err = d.Run(); err != nil {
		log.Error(err.Error())
		return
	}

	return
}

func (d *Container) initEnvironment() (err error) {

	// Environment client
	d.client, err = client.NewClientWithOpts(client.WithVersion(DockerClientVersion))
	if err != nil {
		return
	}

	d.updateContext()
	return
}

func (d *Container) updateContext() {
	// Create a new Context
	ctx := context.Background()
	// Create a new Context, with its cancellation function
	// from the original Context
	d.Context, d.Cancel = context.WithCancel(ctx)
}

func PullImage(image string) error {

	d := &Container{}
	err := d.initEnvironment()
	if err != nil {
		return err
	}

	// Pulling image
	log.Debug("Pulling image %s", image)
	out, err := d.client.ImagePull(d.Context, image, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(out)
	log.Debug(buf.String())

	d.Context.Done()
	return nil
}

// Run starts a Docker Container
func (d *Container) Run() (err error) {

	Docker := d.client

	// Load the environment
	if err = d.initEnvironment(); err != nil {
		return
	}

	// Create the container
	config := &container.Config{
		Image: d.Image,
		Cmd: []string{
			// Check here extra switches:
			// https://peter.sh/experiments/chromium-command-line-switches/
			"--headless",
			"--disable-gpu",
			"--no-default-browser-check",
			"--no-pings",
			"--no-sandbox",
			"--disable-notifications",
			"--disable-sync",
			"--disable-web-security",
			// "--user-data-dir=/tmp",
			"--remote-debugging-address=0.0.0.0",
			"--window-size=1680,2000",
			fmt.Sprintf("--remote-debugging-port=%s", d.LocalPort),
		},
		ExposedPorts: nat.PortSet{nat.Port(d.LocalPort): struct{}{}},
	}

	hostConfig := &container.HostConfig{
		PortBindings: map[nat.Port][]nat.PortBinding{
			nat.Port(d.LocalPort): {
				{HostIP: d.IP, HostPort: d.PublicPort},
			},
		},
		Binds: d.MounthPaths,
	}

	ctx := context.Background()
	containerResp, err := Docker.ContainerCreate(ctx, config, hostConfig, nil, d.Name)
	if err != nil {
		log.Debug("%v", containerResp.Warnings)
		return
	}

	// starting the container
	log.Debug("instructing the docker daemon to start (%s)[%s]", d.Name, containerResp.ID)
	if err = Docker.ContainerStart(ctx, containerResp.ID, types.ContainerStartOptions{}); err != nil {
		return
	}

	log.Debug("Container ready")
	out, err := Docker.ContainerLogs(d.Context, containerResp.ID, types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true})
	if err != nil {
		return
	}
	buf := new(bytes.Buffer)
	_, _ = buf.ReadFrom(out)
	log.Debug(buf.String())

	// fetching debugger URL
	target, err := d.getDebuggerURL()
	if err != nil {
		return err
	}

	d.Target = target
	return
}

func (d *Container) getDebuggerURL() (target ChromeTarget, err error) {

	url := fmt.Sprintf("http://%s:%s/json", d.IP, d.PublicPort)
	log.Important("Retrieving debugger URL from %s", url)

	httpClient := http.Client{
		Timeout: time.Second * 10, // Maximum of 10 secs
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	req.Close = true

	var (
		response *http.Response
		retries  int = 10
	)

	for retries > 0 {
		response, err = httpClient.Do(req)
		if err != nil {
			log.Error(err.Error())
			retries -= 1
			time.Sleep(2 * time.Second)
		} else {
			break
		}
	}

	if err != nil {
		return
	}

	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Error(err.Error())
	}

	log.Debug("data = %s\n", data)

	t := make([]ChromeTarget, 0)
	if err = json.Unmarshal(data, &t); err != nil {
		log.Error(err.Error())
		return
	}

	log.Important(t[0].WebSocketDebuggerUrl)
	return t[0], nil

}
