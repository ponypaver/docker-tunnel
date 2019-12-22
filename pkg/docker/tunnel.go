package docker

import (
	"fmt"
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
		err = fmt.Errorf("failed to listen local docker, error: %v", err)
		return
	}

	log.Printf("docker tunnel to %v:%v started at localhost:%v", t.remoteHost, remoteDockerSocket, t.localUnixSocketFile)

	for {
		select {
		case <-t.done:
			return
		default:

		}

		localConn, acceptErr := listener.Accept()
		if err != nil {
			log.Printf("failed to accept local connection, error: %v", acceptErr)
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
		src.Close()
		dst.Close()
		log.Printf("request from localhost:%v to %v:%v done", src.LocalAddr(), t.remoteHost, dst.RemoteAddr())
	}()

	go mustCopy(src, dst)
	go mustCopy(dst, src)

	<-t.done
	log.Println("forward done")
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
