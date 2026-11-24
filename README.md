# Hello, World!

A demo of best DevSecOps practices with focus on automation, collaboration, security and documentation.

## Requirements

1. [Go 1.27.4 or later](https://go.dev/doc/install)
2. [git cli](https://git-scm.com/downloads)
3. [make](https://www.gnu.org/software/make/) - Optional. In case of Windows OS, it is possible to execute the build commands defined in the Makefile from the command line.

4. [Kind Kubernetes](https://kind.sigs.k8s.io/) 

## Build hello application

Clone hello repository. 

```sh
git clone https://github.com/dgeorgievski/hello.git
cd hello
```

`make build` stamps the binary with the version, short Git commit, and RFC 3339
commit date. On a `vX.Y.Z` tag the version is `X.Y.Z`; otherwise it is `dev`.
Set `TAG` to override the version for both binary and image builds.

```sh
make build
./bin/dgeorgievski-hello version
make check-build TAG=0.0.0-smoke  # Requires curl, jq, and an unused port 8080.
```

## Build hello container image

```sh
make docker-build
docker run --rm ghcr.io/dgeorgievski/hello:dev /hello version
```

The example uses an untagged checkout. A tagged checkout uses its release
version as the image tag instead of `dev`. A Docker build from a tagged Git
checkout uses the same version even when no build arguments are provided.

## Kubernetes deployment

### Kind K8S cluster

kind is a tool for running local Kubernetes clusters using Docker container “nodes”.

Create a local `hello` kind cluster.

```sh
$ kind create cluster --name hello
$ kind get clusters
hello
```

To run a local build instead of the published release, load the `dev` image:

```sh
make kind-docker-upload TAG=dev
```

The raw manifests and the Helm chart default to the published `2.2.0` image.
To use the local build with Helm, pass `--set image.tag=dev`. For the raw
manifests, apply them first, then override the deployment image:

```sh
kubectl -n hello set image deployment/hello hello=ghcr.io/dgeorgievski/hello:dev
```

Create `hello` namespace

```sh
$ kubectl create ns hello
namespace/hello created
```

Deploy `hello` app in the newly created namespace

```sh
$ kubectl create -f deploy/manifests
deployment.apps/hello created
serviceaccount/hello created
service/hello created

$ kubectl get all -n hello
NAME                     READY   STATUS    RESTARTS   AGE
kubectl get all -n hello
NAME                         READY   STATUS    RESTARTS   AGE
pod/hello-69c9cf5648-92hc9   1/1     Running   0          52s

NAME            TYPE        CLUSTER-IP     EXTERNAL-IP   PORT(S)   AGE
service/hello   ClusterIP   10.96.28.147   <none>        80/TCP    52s

NAME                    READY   UP-TO-DATE   AVAILABLE   AGE
deployment.apps/hello   1/1     1            1           53s

NAME                               DESIRED   CURRENT   READY   AGE
replicaset.apps/hello-69c9cf5648   1         1         1       52s
```


### Debugging hello Pod

```sh
kubectl -n hello \
debug -it \
--container=debug-container \
--image=alpine \
--target=hello \
hello-69c9cf5648-92hc9

Targeting container "hello". If you don't see processes from this container it may be because the container runtime doesn't support this feature.
--profile=legacy is deprecated and will be removed in the future. It is recommended to explicitly specify a profile, for example "--profile=general".
If you don't see a command prompt, try pressing enter.
/ # ps axww
PID   USER     TIME  COMMAND
    1 root      0:00 /hello server
   16 root      0:00 /bin/sh
   26 root      0:00 ps axww
/ #
/ # apk add curl
/ # apk add jq
/ # curl -s localhost:8080/hello | jq
{
  "message": "Hello, World!"
}
/ # curl -s localhost:8080/version | jq
{
  "COMMIT": "1de9255",
  "DATE": "2026-11-23T22:10:04-05:00",
  "VERSION": "2.2.0"
}
```

## Health Endpoints

The application exposes two Kubernetes health check endpoints for liveness and readiness probes:

### Liveness Probe

```sh
curl -s localhost:8080/healthz
```

Response:
```json
{"status":"ok"}
```

Returns HTTP 200 when the application process is running. Kubernetes uses this to restart a container that is in a broken state.

### Readiness Probe

```sh
curl -s localhost:8080/readyz
```

Response:
```json
{"status":"ready"}
```

Returns HTTP 200 when the application is ready to accept traffic. Kubernetes uses this to add or remove the pod from service endpoints.

## CI build and release workflows

The build workflow runs `govulncheck` and verifies the built binary's CLI and
HTTP version metadata before delegating to the shared build pipeline. The
`server` command exits nonzero on startup failure and bounds the HTTP header
read, response write, and idle connection times.

Release tags use `vX.Y.Z`; the version embedded in GoReleaser binaries and the
full container tag is `X.Y.Z`. When updating the default deployed release,
update `pkg/version/version.go`, Helm `appVersion`, the raw deployment image
tag, and the deploy workflow's default tag together.

## Git Hooks scripts

The use of git hooks has been documented in [doc/git-hooks.md](doc/git-hooks.md) file.

## Git commands

Git commands are implemented in NuShell. They are wrappers around regular git commands tha provide additional functionality.

### Git Squash

Usage TODO

## Semantic Releases

TODO
