package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/ponypaver/docker-tunnel/pkg/docker"
)

var (
	host string
)

func init() {
	flag.StringVar(&host, "h", "", "remote hostname or ip")
	flag.Parse()
}

func main() {
	if host == "" {
		fmt.Fprintln(os.Stderr, "host should be specified")
		flag.Usage()
		os.Exit(1)
	}

	cli, err := docker.NewTunneledClient(host)

	if err != nil {
		panic(err)
	}

	containers, err := cli.ContainerList(context.Background(), types.ContainerListOptions{})
	if err != nil {
		panic(err)
	}
	for _, container := range containers {
		fmt.Printf("%s %s\n", container.ID[:10], container.Image)
	}
}