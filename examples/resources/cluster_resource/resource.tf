resource "swarm_cluster" "local_cluster" {
  nodes {
    hostname        = "localhost"
    public_address  = "127.0.0.1"
    private_address = "127.0.0.1"
    tags = {
      role = "manager"
    }
  }
}