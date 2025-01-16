# Prometheus Headscale Exporter

[![lint](https://github.com/msterhuj/headscale_exporter/actions/workflows/lint.yml/badge.svg?branch=main)](https://github.com/msterhuj/headscale_exporter/actions/workflows/lint.yml)

This is a simple exporter for headscale its aim to export data not exported by the headscale API/Metrics endpoint.

This is a work in progress and will be updated with more metrics, options, filters, etc.

## Usage

```bash
headscale-exporter -h
```

## Metrics

The exporter will expose the following metrics:

```text
# HELP headscale_api_keys Number of API keys
# TYPE headscale_api_keys gauge
headscale_api_keys 1
# HELP headscale_node Number of nodes
# TYPE headscale_node gauge
headscale_node 2
# HELP headscale_node_online Number of online nodes
# TYPE headscale_node_online gauge
headscale_node_online 1
# HELP headscale_policy_acl Number of policies
# TYPE headscale_policy_acl gauge
headscale_policy_acl 6
# HELP headscale_policy_group Number of groups
# TYPE headscale_policy_group gauge
headscale_policy_group 1
# HELP headscale_policy_group_members Number of members in each group
# TYPE headscale_policy_group_members gauge
headscale_policy_group_members{group="group:cloud-edge"} 1
# HELP headscale_policy_host Number of hosts
# TYPE headscale_policy_host gauge
headscale_policy_host 4
# HELP headscale_policy_tag_owner Number of tag owners
# TYPE headscale_policy_tag_owner gauge
headscale_policy_tag_owner 0
# HELP headscale_user Number of users
# TYPE headscale_user gauge
headscale_user 2
```

## Building

```bash
go build
```
