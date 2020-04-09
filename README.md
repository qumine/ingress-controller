QuMine - Ingress
---
![GitHub Release](https://img.shields.io/github/v/release/quhive/qumine-ingress)
![GitHub Workflow](https://img.shields.io/github/workflow/status/quhive/qumine-ingress/release)
[![GoDoc](https://godoc.org/github.com/quhive/QuMine-Ingress?status.svg)](https://godoc.org/github.com/quhive/qumine-Ingress)
[![Go Report Card](https://goreportcard.com/badge/github.com/quhive/QuMine-Ingress)](https://goreportcard.com/report/github.com/quhive/qumine-Ingress)

Kubernetes ingress controller for minecraft servers

# Status

- [x] Basic routing of minecraft clients based on hostnames
- [x] Basic Logging
- [x] Metrics (Prometheus)
- [ ] Healthchecks for Ingress
- [ ] Healthchecks for Upstreams
- [ ] Custom error Handling for if the upstream is not available (Status response?)
- [ ] Frontend to quickly see all upstreams

# Usage


## Kubernetes

*HELM Charts can be found here: [quhive/charts](https://github.com/quhive/charts)*

### Ingress

*The ingress should be run as a daemonset on all of your outwards facing nodes.*

By default the ingress should run fine without customization, but if you need to the behaviour of the ingress can be customized by setting a couple of arguments. Here is the full list of available arguments

```
  -api-host string
        Address the rest api will listen on (default "0.0.0.0")
  -api-port int
        Port the rest api will listen on (default 8080)
  -debug
        Enable debugging log level
  -help
        Show this page
  -host string
        Address the server will listen on (default "0.0.0.0")
  -kube-config string
        Path of the kube config file to use
  -port int
        Port the server will listen on (default 25565)
  -version
        Show the current version
```

### Upstream Services

To enable a service to be discovered by the ingress it needs to have the ```qumine-ingress.quhive.io/hostname``` annotations.
Optionaly you can set the ```qumine-ingress.quhive.io/portname``` annotation to define which port will be used for the minecraft connection.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: example
  annotations:
    qumine-ingress.quhive.io/hostname: "example"
    qumine-ingress.quhive.io/portname: "minecraft"
spec:
  ports:
  - port: 25565
    name: minecraft
  selector:
    app: example
```

## Outside of Kubernetes

```
Will follow in the near future
```

# Development

## Perfrom a Snapshot release locally

```
docker run -it --rm \
  -v ${PWD}:/build -w /build \
  -v /var/run/docker.sock:/var/run/docker.sock \
  goreleaser/goreleaser \
  release --snapshot --rm-dist
```