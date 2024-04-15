# Chord
An implementation of the Chord protocol (Stoica et al.) using Go and gRPC, along with an example application for storing key-value pairs in memory.

- `chord` is the core package containing the Chord logic
- `dht` is the example application which consumes `chord` and exposes a basic service for storing arbitrary bytes in key-value pairs
- `protos` contains the protobuf definitions for both `chord` and `dht`, these can be generated using `make`.
- `test_bench` contains the scripts and configuration files for evaluating the application.

## Installation

You must have an installation of [Go](https://go.dev/doc/install) (1.21+) and [protoc](https://grpc.io/docs/protoc-installation/).

1. Run `make` to generate the Go source code from the protobufs in `protos`
2. Run `go install`

# Usage

### Example Usage 
Start an initial node
```bash
chord_dht -port 8080
```
In another terminal, start another node and use the existing node as a bootstrap.
```bash
chord_dht -bootstrap 127.0.0.1:8080 -port 8081
```
The nodes should stabilize and acknowledge each other as successors.

The accompanying tools in `python_library` can be used to demonstrate setting and retrieving keys, e.g.:

```bash
python set_key.py 127.0.0.1 test test123
python get_key.py 127.0.0.1 test
```

When stopping a process with a SIGTERM (CTRL+C), the node will transfer the keys to its immediate successor. The node will also continuously transfer away any keys that don't belong to it.


## External Addresses
If you want the node to be contacted on an address other than `127.0.0.1`, the application must be informed by setting the `-address` flag. This is important as the node's Chord identifier is based on this.

For example:
`chord_dht -address 10.24.0.1`
or: `chord_dht -address $(hostname -i)`

Be careful not to mix up the use of `127.0.0.1` and `localhost`, as these will result in different IDs.

## Local Test Bench
To run networks on a local setup, it's easiest to use the `docker-compose.yaml`, which builds Docker images based on the local source code and bootstraps 10 nodes. 

`docker-compose up`

There's also a test container which mounts the source code from `python_library` which allows you to run Python scripts in the context of the Docker network, use the `./chord_python` executable in place of `python`.

```bash
./chord_python set_key.py peer2 test hello
```



## Test Bench

To run larger experiments, the `test_bench` folder contains some tools to run networks on Kubernetes.

To run a test scenario, first apply the `lead_deployment.yaml` and then update `peer_deployment.yaml` with the lead pod's IP address so that the other nodes can use it to bootstrap, then apply `peer_deployment.yaml`.

You can change the number of replicas to your liking. Each peer pod requests 32MiB of memory, but some will go over this in practice.

### Instrumentation

The `chord_dht` program exports some Prometheus metrics on `:2112`, apply the `pod_monitor.yaml` to tell Prometheus to collect them. See `scripts/setup_kube_prometheus.sh` to set up a Prometheus and Grafana installation.