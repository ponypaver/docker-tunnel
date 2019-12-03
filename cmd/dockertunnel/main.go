package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ponypaver/docker-tunnel/pkg/docker"
	"github.com/ponypaver/docker-tunnel/pkg/ssh"
)

const (
	defaultSSHUser        = "root"
	defaultSSHPort        = 22
	defaultPrivateKeyPath = "~/.ssh/id_rsa"
)

var (
	port           int
	user           string
	host           string
	password       string
	privateKeyPath string

	shutdownHandler      chan os.Signal
	onlyOneSignalHandler = make(chan struct{})
	shutdownSignals      = []os.Signal{os.Interrupt, syscall.SIGTERM}
)

func init() {
	flag.IntVar(&port, "P", defaultSSHPort, "ssh port")
	flag.StringVar(&host, "h", "", "remote hostname or ip")
	flag.StringVar(&user, "u", defaultSSHUser, "ssh user to login to remote host")
	flag.StringVar(&password, "p", "", "ssh password")
	flag.StringVar(&privateKeyPath, "i", defaultPrivateKeyPath, "ssh private key to use to create ssh tunnel")

	flag.Parse()
}

// https://github.com/kubernetes/apiserver/blob/master/pkg/server/signal.go
// SetupSignalHandler registered for SIGTERM and SIGINT. A stop channel is returned
// which is closed on one of these signals. If a second signal is caught, the program
// is terminated with exit code 1.
func setupSignalHandler() <-chan struct{} {
	close(onlyOneSignalHandler) // panics when called twice

	shutdownHandler = make(chan os.Signal, 2)

	stop := make(chan struct{})
	signal.Notify(shutdownHandler)
	go func() {
		<-shutdownHandler
		close(stop)
		<-shutdownHandler
		os.Exit(1) // second signal. Exit directly.
	}()

	return stop
}

func main() {
	if host == "" {
		fmt.Fprintln(os.Stderr, "host should be specified")
		flag.Usage()
		os.Exit(1)
	}

	if password == "" && privateKeyPath == "" {
		fmt.Fprintln(os.Stderr, "at least one of password and private key should be specified")
		flag.Usage()
		os.Exit(1)
	}

	client, err := ssh.NewSSHClient(user, host, password, privateKeyPath, port)
	if err != nil {
		log.Fatal(err)
	}

	tunnel := docker.NewTunnel(client, host)
	go tunnel.Start()

	stopCh := setupSignalHandler()
	<-stopCh // received either SIGHUP, or SIGTERM
	log.Printf("signal received, closing docker tunnel")

	tunnel.Close()
}
