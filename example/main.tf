terraform {
    required_providers {
        swarm = {
            source = "australiacloud.com.au/aom/swarm"
            version = "0.1"
        }
    }
}

provider "swarm" {
}

data "swarm_nodes" "local_nodes" {
}

output "nodes" {
  value = data.swarm_nodes.local_nodes.all
}
