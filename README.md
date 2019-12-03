## docker-tunnel

### what is this

Basically, this project implemented the following line of ssh tunnel

```
ssh -nNT -L /tmp/remoteHost.docker.sock:/var/run/docker.sock user@remoteHost
```

It's main purpose is to call remote docker socket with a local stub.

### build

default make to adapt current platform

```
$ make 
```

cross compile, to specify platform (eg: linux)

```
make dockertunnel-os-linux
```

### sample docker client code

```
func NewTunneledClient(remoteHost string) (*client.Client, error) {
	localDockerSocketFile := docker.ComposeLocalDockerSocketFile(remoteHost)
	opt := client.WithHost(localDockerSocketFile)

	cli, err := client.NewClientWithOpts(opt)
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client, error: %v", err)
	}

	return cli, nil
}
```

### sample output

```
$ make
$ build/bin/darwin/dockertunnel -h 10.0.0.10 -P 1122
2019/12/03 11:28:21 docker tunnel to 10.0.0.10:/var/run/docker.sock started at localhost:/tmp/10.0.0.10.docker.sock
2019/12/03 11:28:30 accepted a docker request from localhost:/tmp/10.0.0.10.docker.sock
2019/12/03 11:28:30 forwarding request to: 10.0.0.10:/var/run/docker.sock
2019/12/03 11:28:30 request from localhost:/tmp/10.0.0.10.docker.sock to 10.0.0.10:/var/run/docker.sock done
^C2019/12/03 11:28:32 signal received, closing docker tunnel

```
