terraform {
    required_providers {
        swarm = {
            source = "aucloud/swarm"
            version = "1.0.0"
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

data "swarm_nodes" "local_nodes" {
}

data "swarm_cluster" "local_cluster" {
}

output "nodes" {
  value = data.swarm_nodes.local_nodes.all
}

output "cluster" {
  value = data.swarm_cluster.local_cluster
}
