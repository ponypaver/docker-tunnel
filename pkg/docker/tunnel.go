package docker

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"

	"golang.org/x/crypto/ssh"
)

const (
	localSocketDir     = "/tmp/"
	dockerSocketSuffix = ".docker.sock"
	remoteDockerSocket = "/var/run/docker.sock"
)

type Tunnel struct {
	done                chan struct{}
	remoteHost          string
	localUnixSocketFile string
	sshClient           *ssh.Client
}

func NewTunnel(sshClient *ssh.Client, host string) *Tunnel {
	return &Tunnel{
		remoteHost:          host,
		sshClient:           sshClient,
		done:                make(chan struct{}),
		localUnixSocketFile: composeLocalDockerSocketFile(host),
	}
}

func (t *Tunnel) Start() (err error) {
	listener, err := net.Listen("unix", t.localUnixSocketFile)
	if err != nil {
		return fmt.Errorf("failed to listen local docker, error: %v", err)
	}

	log.Printf("docker tunnel to %v:%v started at localhost:%v", t.remoteHost, remoteDockerSocket, t.localUnixSocketFile)

	for {
		select {
		case <-t.done:
			return
		default:

		}

		localConn, err := listener.Accept()
		if err != nil {
			log.Printf("failed to accept local connection, error: %v", err)
			continue
		}

		log.Printf("accepted a docker request from localhost:%v", localConn.LocalAddr())

		dstConn, err := t.sshClient.Dial("unix", remoteDockerSocket)
		if err != nil {
			log.Printf("failed to dial %v:%v, error: %v", t.remoteHost, remoteDockerSocket, err)
			continue
		}

		log.Printf("forwarding request to: %v:%v", t.remoteHost, dstConn.RemoteAddr())

		go t.forward(dstConn, localConn)
	}
}

func (t *Tunnel) forward(dst, src net.Conn) {
	defer func() {
		dst.Close()
		src.Close()
		log.Printf("request from localhost:%v to %v:%v done", src.LocalAddr(), t.remoteHost, dst.RemoteAddr())
	}()

	go mustCopy(src, dst)
	mustCopy(dst, src)
}

func (t *Tunnel) Close() (err error) {
	close(t.done)

	if fileExist(t.localUnixSocketFile) {
		if err = os.Remove(t.localUnixSocketFile); err != nil {
			err = fmt.Errorf("failed to remove %v", t.localUnixSocketFile)
		}
	}

	return
}

func composeLocalDockerSocketFile(name string) string {
	return localSocketDir + name + dockerSocketSuffix
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
