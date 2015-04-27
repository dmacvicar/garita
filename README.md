
# Garita

Small Docker registry auth server in Go

# Environment

The environment creates 3 containers:

* a Docker daemon (dockerd, dockerd.test.lan)
* a Registry (registry, registry.test.lan)
* garita (garita, garita.test.lan)

# Running

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
