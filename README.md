# Terraform Provider for Docker Swarm

[![Go](https://github.com/aucloud/terraform-provider-swarm/actions/workflows/go.yml/badge.svg)](https://github.com/aucloud/terraform-provider-swarm/actions/workflows/go.yml)

`terraform-provider-swarm` is a [Terraform](https://terraform.io) provider for the creation and management of [Docker Swarm](https://docs.docker.com/engine/swarm/) clusters (_an alternative container orchestrator to Kubernetes and Nomad_)

## Quick Start

Run the following command to build the provider

```console
make install
```

### Sample Configuration

`main.tf`:
```terraform
terraform {
    required_providers {
        swarm = {
            source = "aucloud/swarm"
            version = "~> 1.2"
        }
    }
}

provider "swarm" {
  use_local = true
}

resource "swarm_cluster" "local_cluster" {
  nodes {
    hostname = "localhost"
    public_address = "127.0.0.1"
    private_address = "127.0.0.1"
    tags = {
      role = "manager"
    }
  }
}

output "cluster" {
  value = data.swarm_cluster.local_cluster
}
```

Initialize and apply the Terraform:
```console
terraform init
terraform apply
```

For a full example see [examples/main.tf](/examples/main.tf)

## License

`terraform-provider-swarm` is licensed under the terms of the [AGPLv3](/LICENSE)
