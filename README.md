# Prometheus Headscale Exporter

This is a simple exporter for headscale its aim to export data not exported by the headscale API/Metrics endpoint.

This is a work in progress and will be updated with more metrics, options, filters, etc.

## Usage

```bash
$ headscale-exporter -h
```

## Metrics

The exporter will expose the following metrics:

```
# HELP headscale_headscale_api_keys Number of API keys
# TYPE headscale_headscale_api_keys gauge
headscale_headscale_api_keys 1
# HELP headscale_headscale_node Number of nodes
# TYPE headscale_headscale_node gauge
headscale_headscale_node 2
# HELP headscale_headscale_node_online Number of online nodes
# TYPE headscale_headscale_node_online gauge
headscale_headscale_node_online 1
# HELP headscale_headscale_user Number of users
# TYPE headscale_headscale_user gauge
headscale_headscale_user 2
```

## Building

```bash
$ go build
```
