---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "swarm Provider"
subcategory: ""
description: |-
  
---

# swarm Provider



## Example Usage

```terraform
provider "swarm" {
  use_local               = true
  skip_manager_validation = true

  // ssh_user = "terraform
  // ssh_key = "$HOME/.ssh/terraform_rsa""
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Optional

- **ssh_addr** (String)
- **ssh_key** (String, Sensitive)
- **ssh_timeout** (String)
- **ssh_user** (String)
- **use_local** (Boolean)
