
[![Build Status](https://travis-ci.org/dmacvicar/garita.svg?branch=master)](https://travis-ci.org/dmacvicar/garita)

# Garita

Small Docker v2 registry auth server in Go.

It exists mostly as a project to learn Go, the Vagrant Docker provider, and
to understand the protocol [Portus](https://github.com/SUSE/Portus) implements.

## Features

* Authentication is only supported using htpasswd files
* Once authenticated, it provides push and pull access to the
  /$user namespace

Garita is inspired in [Portus](https://github.com/SUSE/Portus), which
is a full featured auth server and registry index.

## Running

```
garita -key path/to/server.key -htpasswd path/to/htpasswd
```

At the same time you need to configure the registry

```
auth:
  token:
    realm: http://garita.yourdomain.com/v2/token
    service: registry.yourdomain.com
    issuer: garita.yourdomain.com
    rootcertbundle: /path/to/server.crt
```

## Development Environment

The environment creates 3 containers:

* a Docker daemon (dockerd, dockerd.test.lan)
* a Registry (registry, registry.test.lan)
* garita (garita, garita.test.lan)

While the images are based on opensuse:13.2, the dockerd container requires a host kernel
with overlayfs support. (eg. openSUSE Tumbleweed or another distribution supporting
overlayfs). The dockerd container is already privileged but I don't want to mess with the loop
devices of the host.

## Running

* Compile

```
go install github.com/dmacvicar/garita
```

* Start the environment

```
vagrant up --no-parallel
```

* Everytime you rebuild

```
vagrant reload garita
```

* To see the logs

```
vagrant docker-logs -f garita
```

Run docker against the docker daemon running inside the container

```
docker -H tcp://localhost:23750 images
```

The typical testcase, pull busybox, tag it, and push it to the registry

```
docker -H tcp://localhost:23750 pull busybox
docker -H tcp://localhost:23750 tag busybox registry.test.lan/duncan/busybox
docker login registry.test.lan
docker -H tcp://localhost:23750 push registry.test.lan/duncan/busybox
```

# Bugs

The [specification](https://docs.docker.com/registry/spec/auth/token/) does not go into every detail. If I missed something please open an issue.

# Authors

* Duncan-Mac-Vicar P. <dmacvicar@suse.de>

# License

* Garita is licensed under the Apache 2.0 license.
