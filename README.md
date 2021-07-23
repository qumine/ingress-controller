QuMine - Ingress
---
![GitHub Release](https://img.shields.io/github/v/release/qumine/ingress-controller)
![GitHub Workflow](https://img.shields.io/github/workflow/status/qumine/ingress-controller/release)
[![GoDoc](https://godoc.org/github.com/qumine/ingress-controller?status.svg)](https://godoc.org/github.com/qumine/ingress-controller)
[![Go Report Card](https://goreportcard.com/badge/github.com/qumine/ingress-controller)](https://goreportcard.com/report/github.com/qumine/ingress-controller)
[![Quality Gate Status](https://sonarcloud.io/api/project_badges/measure?project=qumine_minecraft-server&metric=alert_status)](https://sonarcloud.io/dashboard?id=qumine_ingress-controller)

Kubernetes ingress controller for routing Minecraft connections based on the requested hostname

# Usage

## Kubernetes

*HELM Charts can be found here: [qumine/charts](https://github.com/qumine/charts)*

### Ingress

*The ingress should be run as a daemonset on all of your outwards facing nodes.*

By default the ingress should run fine without customization, but if you need to the behaviour of the ingress can be customized by setting a couple of arguments. Here is the full list of available arguments.

```
Usage:
  ingress-controller [flags]

Flags:
      --api-host string      Host for the API server to listen on (default "0.0.0.0")
      --api-port int         Port for the API server to listen on (default 8080)
  -d, --debug                Debug logging
  -h, --help                 help for ingress-controller
      --host string          Host for the API server to listen on (default "0.0.0.0")
      --kube-config string   KubeConfig path
      --port int             Port for the API server to listen on (default 25565)
      --trace                Trace logging
  -v, --version              version for ingress-controller
```

**All configuration options can also be set via environment variables** 

### Upstream Services

To enable a service to be discovered by the ingress it needs to have the ```ingress.qumine.io/hostname``` annotations.
Optionaly you can set the ```ingress.qumine.io/portname``` annotation to define which port will be used for the minecraft connection.

```yaml
apiVersion: v1
kind: Service
metadata:
  name: example
  annotations:
    ingress.qumine.io/hostname: "example"
    ingress.qumine.io/portname: "minecraft"
spec:
  ports:
  - port: 25565
    name: minecraft
  selector:
    app: example
```

## Outside of Kubernetes

If you want to run the ingress outside of kubernetes you can do so by providing the ```--kube-config``` flag or environment variable. Keep in mind tho that the routing towards the internal kubernetes services needs to be configured.

```
./ingress-controller --kube-config ~/.kube/config
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