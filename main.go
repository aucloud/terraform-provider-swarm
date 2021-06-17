package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"gitlab.mgt.aom.australiacloud.com.au/aom/terraform-provider-swarm/swarm"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return swarm.Provider()
		},
	})
}
