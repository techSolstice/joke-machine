# Usage Instructions

## Building and Running the application locally

Execute `docker build -t joke-machine:1 .` within the root of the directory where this `README.md` is located.  This will create the docker image for the application.  You can add additional tags if desired to the image name, e.g. `joke-machine:1.0` or `joke-machine:latest` in addition to just `1`.

After a successful build, if you wish to test locally, run the built image via `docker run -p 5000:5000 joke-machine:1`, changing the tag if needed.  Also change the host and container ports if desired from `-p <host>:<container>`.

## Test the Kubernetes cluster

All this could be run locally via `docker compose` if desired, but for time's sake, I'll provide instructions for `kind` (or `minikube` can be used if you desire).

Execute `kind create cluster` to create the cluster with default name `kind`.
Execute `kind load docker-image joke-machine:1 --name kind` to give your new cluster access to the built docker image.
Execute `helm install joke-machine ./helm/` to deploy infrastructure locally.
Follow the instructions output in the CLI.

You should now be able to visit http://127.0.0.1:5000/health to get a 200 `healthy` response back (JSON string).

If you want to modify the K8s manifests now, run `helm upgrade joke-machine ./helm/` to apply chart changes.

## Local Development

Execute `go run joke-machine.go` as the main entrypoint and therefore gaining access to http://localhost:5000/health as well.

## Entrypoints

Go to http://localhost:5000/health to get a JSON message with healthy returned.
Go to http://localhost:5000/joke to get a templated browser version of a quote.
Go to or curl http://localhost:5000 to get a plaintext version of a quote.