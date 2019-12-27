# docker-tunnel

## What is this

Basically, this project implemented the following line of ssh tunnel

```
$ ssh -nNT -L /tmp/remoteHost.docker.sock:/var/run/docker.sock user@remoteHost
```

ie, establish a docker call tunnel through ssh which enables to call remote docker socket with a local stub.

## Build

default make to adapt current platform

```
$ make 
```

cross compile, to specify platform (eg: linux)

```
$ make dockertunnel-os-linux
```

## Sample docker client code

```
there is a sample client demo at /cmd/cli
```

## Sample output

```
$ make
$ build/bin/darwin/dockertunnel -h 10.0.0.10 -P 1122
2019/12/03 11:28:21 docker tunnel to 10.0.0.10:/var/run/docker.sock started at localhost:/tmp/10.0.0.10.docker.sock
2019/12/03 11:28:30 accepted a docker request from localhost:/tmp/10.0.0.10.docker.sock
2019/12/03 11:28:30 forwarding request to: 10.0.0.10:/var/run/docker.sock
2019/12/03 11:28:30 request from localhost:/tmp/10.0.0.10.docker.sock to 10.0.0.10:/var/run/docker.sock done
^C2019/12/03 11:28:32 signal received, closing docker tunnel

```

## Note

- This implementation only implemented docker unix socket tunnel
- This implementation only test on Ubuntu 16.04
- This can not work on CentOs/Rhel 7 release for known [bug](https://bugzilla.redhat.com/show_bug.cgi?id=1527565)

## License

apache 2.0
