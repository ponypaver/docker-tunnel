package ssh

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"time"

	"golang.org/x/crypto/ssh"
	homedir "github.com/mitchellh/go-homedir"
)

const (
	defaultTimeout = time.Second * 60
)

func newClientConfig(user, passWd, keyPath string) (*ssh.ClientConfig, error) {
	var authMethods []ssh.AuthMethod

	if keyPath[:2] == "~/" {
		home, err := homedir.Dir()
		if err != nil {
			return nil, err
		}
		keyPath = filepath.Join(home, keyPath[2:])
	}

	keyData, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, fmt.Errorf("read private key error: %v", err)
	}

	privateKey, err := ssh.ParsePrivateKey(keyData)
	if err != nil {
		return nil, fmt.Errorf("parse private key error: %v", err)
	}

	authMethods = append(authMethods, ssh.Password(passWd), ssh.PublicKeys(privateKey))

	return &ssh.ClientConfig{
		User:            user,
		Auth:            authMethods,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         defaultTimeout,
	}, nil

}

func NewSSHClient(user, host, passWd, keyPath string, port int) (*ssh.Client, error) {
	config, err := newClientConfig(user, passWd, keyPath)
	if err != nil {
		return nil, fmt.Errorf("create ssh client config error: %v", err)
	}

	client, err := ssh.Dial("tcp", fmt.Sprintf("%v:%v", host, port), config)
	if err != nil {
		return nil, fmt.Errorf("dial %v:%v error: %v", host, port, err)
	}

	return client, nil
}
