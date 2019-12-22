package docker

import (
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/client"
)

func NewTunneledClient(remoteHost string) (*client.Client, error) {
	localDockerSocketFile := composeLocalDockerSocketFile(remoteHost)

	cli, err := client.NewClientWithOpts(
		client.WithHost("unix://"+localDockerSocketFile),
	)

	//cli, err := client.NewClient("unix://"+localDockerSocketFile)

	if err != nil {
		return nil, fmt.Errorf("failed to create docker client, error: %v", err)
	}

	return cli, nil
}

func composeLocalDockerSocketFile(remoteHost string) string {
	return localSocketDir + remoteHost + dockerSocketSuffix
}

func fileExist(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}
	return true
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		panic(err)
	}
}
