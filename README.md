# Terraform Provider for Docker Swarm

[![Go](https://github.com/aucloud/terraform-provider-swarm/actions/workflows/go.yml/badge.svg)](https://github.com/aucloud/terraform-provider-swarm/actions/workflows/go.yml)

`terraform-provider-swarm` is a [Terraform](https://terraform.io) provider for the creation and management of [Docker Swarm](https://docs.docker.com/engine/swarm/) clusters (_an alternative container orchestrator to Kubernetes and Nomad_).

Run the following command to build the provider

```shell
go build -o terraform-provider-swarm
```

## Test sample configuration

First, build and install the provider.

```shell
make install
```

Then, run the following command to initialize the workspace and apply the sample configuration.

```shell
terraform init && terraform apply
```
