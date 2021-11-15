# Swarm Provider Use Case - 3 node manager swarm
## Use Case Description
The Swarm Provider allows users to create a Docker Swarm cluster using Terraform. This use case shows how three VM instances can be added as managers in a swarm cluster.

## Assumptions
1. The Docker engine is installed and running on the VM instances that will swarm
1. The *swarm_cluster* schema requirements can be created for VM instances from any provider
1. The Swarm Provider uses passwordless SSH key based authentication with the target VMs

## Used directory structure
```
├── cloudinit.tpl
├── data.tf
├── main.tf
├── provider.tf
├── swarm.tfvars
├── variables.tf
└── vm.tf
```

## Swarm Provider
Add the swarm provider information (e.g. provider.tf) and any configuration items (.e.g. local ssh key credentials)
```terraform
terraform {
  required_version = ">= 0.13"
  required_providers {
    vsphere = {
      source  = "hashicorp/vsphere"
      version = "1.21.0"
    }
    swarm = {
      source = "aucloud/swarm"
      version = "1.0.0"
    }
  }
}

... omitted ...

provider "swarm" {
  ssh_user = var.terraform_ssh_user
  ssh_key  = var.terraform_ssh_key
}
```
## Nodes schema
The *swarm_cluster* resource expects a set of nodes that provide VM instance data that can be used to perform swarming. For our 3 manager node use case, we can create the *swarm_cluster* resource using the dynamic **nodes** schema data as follows (in main.tf):
```terraform
resource "swarm_cluster" "cluster" {
  skip_manager_validation = true
  dynamic "nodes" {
    for_each = vsphere_virtual_machine.vm[*]
    content {
        hostname = nodes.value.name
        tags = tomap({
          "role"   = lookup(nodes.value.custom_attributes, data.vsphere_custom_attribute.role.id, "manager"),
          "labels" = lookup(nodes.value.custom_attributes, data.vsphere_custom_attribute.swarm_labels.id, "")
        })
        public_address  = nodes.value.guest_ip_addresses[0]
        private_address = nodes.value.guest_ip_addresses[0] # [2] for private address
    }
  }
  lifecycle {
    prevent_destroy = false
  }
}
```
The tags "role" map item provides information to the swarm provider on the type of Docker node being targeted. The values 'manager' or 'worker' are allowed.

## Swarm Cluster Output
The manager nodes after creation:
```terraform
Apply complete! Resources: 4 added, 0 changed, 0 destroyed.

Outputs:

cluster = {
  "created_at" = "2021-11-11T05:45:55.42101054Z"
  "id" = "um5m2mo3nmi10c8kyh733bafv"
  "nodes" = tolist([
    {
      "hostname" = "swarming-dm1"
      "private_address" = "10.9.9.48"
      "public_address" = "10.9.9.48"
      "tags" = tomap({
        "labels" = ""
        "role" = "manager"
      })
    },
    {
      "hostname" = "swarming-dm2"
      "private_address" = "10.9.9.46"
      "public_address" = "10.9.9.46"
      "tags" = tomap({
        "labels" = ""
        "role" = "manager"
      })
    },
    {
      "hostname" = "swarming-dm3"
      "private_address" = "10.9.9.56"
      "public_address" = "10.9.9.56"
      "tags" = tomap({
        "labels" = ""
        "role" = "manager"
      })
    },
  ])
  "skip_manager_validation" = true
  "updated_at" = tostring(null)
}
```
The Docker node list for the created swarm:
```
$ [~] ssh root@10.9.9.48 "docker node list"
Warning: Permanently added '10.9.9.48' (ECDSA) to the list of known hosts.
ID                            HOSTNAME       STATUS    AVAILABILITY   MANAGER STATUS   ENGINE VERSION
6pplrzfcz1xj43rxx3q582ch9 *   dm1-swarming   Ready     Active         Reachable        20.10.8
xjce4vxyjq280g1riy0lwivix     dm2-swarming   Ready     Active         Reachable        20.10.8
9ostpze3vksrsj1st6lxu1780     dm3-swarming   Ready     Active         Leader           20.10.8
```
