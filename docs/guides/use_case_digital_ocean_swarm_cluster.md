---
page_title: "Digitalocean Docker Swarm Cluster"
---
# Digitalocean Swarm Cluster with 1 manager node and 1 worker node
## Use Case Description
The Swarm Provider allows users to create a Docker Swarm cluster using Terraform. This use case shows how a single manager and node swarm cluster can be created in Digital Ocean.

## Assumptions
1. The Docker engine is installed and running on the VM instances that will swarm
1. The *swarm_cluster* schema requirements can be created for VM instances from any provider
1. The Swarm Provider uses passwordless SSH key based authentication with the target VMs

## Used directory structure
```
├── main.tf
├── outputs.tf
├── provider.tf
├── variables.tf
└── versions.tf
```

## Versions
The required version information is added to versions.tf file:
```terraform
terraform {
  required_providers {
    digitalocean = {
      source = "digitalocean/digitalocean"
      version = "~> 1.22.1"
    }
    swarm = {
      source = "aucloud/swarm"
      version = "1.0.0"
    }
  }
  required_version = ">= 0.13"
}
```

## Swarm Provider
Add the swarm provider information (e.g. provider.tf) and any configuration items (.e.g. local ssh key credentials):
```terraform
provider "digitalocean" {
  token = var.do_token
}
provider "swarm" {
  ssh_user = var.ssh_user
  ssh_key  = var.pvt_key
}
```

## Nodes schema
The *swarm_cluster* resource expects a set of nodes that provide VM instance data that can be used to perform swarming. For our 1 manager and 1 worker node use case, we can create the *swarm_cluster* resource using the dynamic **nodes** schema data as follows (in main.tf):
```terraform
resource "swarm_cluster" "cluster" {
  skip_manager_validation = true
  dynamic "nodes" {
    for_each = concat(digitalocean_droplet.manager, digitalocean_droplet.worker)
    content {
        hostname = nodes.value.name
        tags = tomap({
          "role"   = sort(nodes.value.tags)[0]
        })
        public_address  = nodes.value.ipv4_address
        private_address = nodes.value.ipv4_address_private
    }
  }
  lifecycle {
    prevent_destroy = false
  }
}
```
The tags "role" map item provides information to the swarm provider on the type of Docker node being targeted. The values 'manager' or 'worker' are allowed.
```terraform
data "digitalocean_ssh_key" "terraform" {
  name = "terraform"
}

resource "digitalocean_droplet" "manager" {
  count = 1
  image = var.image_name
  name = "dm${count.index + 1}.${var.cluster}.${var.region}"
  region = var.region
  size = var.instance_type
  private_networking = true

  tags = ["manager"]

  ssh_keys = [data.digitalocean_ssh_key.terraform.id]

  connection {
    host = self.ipv4_address
    user = var.ssh_user
    type = "ssh"
    private_key = file(var.pvt_key)
    timeout = "2m"
  }
}
...
```

## Swarm Cluster Output
The cluster creation output:
```terraform
...
cluster = {
  "created_at" = "2021-11-14T23:55:07.883081731Z"
  "id" = "0906046ued6khe6obw1nc0d5t"
  "nodes" = tolist([
    {
      "hostname" = "dm1.aucloud.SGP1"
      "private_address" = "10.104.0.3"
      "public_address" = "165.22.110.161"
      "tags" = tomap({
        "role" = "manager"
      })
    },
    {
      "hostname" = "dw1.aucloud.SGP1"
      "private_address" = "10.104.0.2"
      "public_address" = "165.22.110.63"
      "tags" = tomap({
        "role" = "worker"
      })
    },
  ])
  "skip_manager_validation" = true
  "updated_at" = tostring(null)
}
```
The Docker node list for the created cluster:
```bash
[rancher@dm1 ~]$ docker node ls
ID                            HOSTNAME            STATUS              AVAILABILITY        MANAGER STATUS      ENGINE VERSION
t9vwmbugft1c8pab37kxogzmk *   dm1.aucloud.SGP1    Ready               Active              Leader              19.03.15
pr41a4b39vek8urocsm4ha895     dw1.aucloud.SGP1    Ready               Active                                  19.03.15
```
