package swarm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/aucloud/go-swarm"
)

func dataSourceClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	swarmManager := m.(*swarm.Manager)

	node, err := swarmManager.GetInfo()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting node info: %w", err))
	}

	var remoteManagers []map[string]interface{}

	for _, remoteManager := range node.Swarm.RemoteManagers {
		remoteManagers = append(remoteManagers, map[string]interface{}{
			"id":   remoteManager.NodeID,
			"addr": remoteManager.Addr,
		})
	}

	if err := d.Set("created_at", node.Swarm.Cluster.CreatedAt); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("nodes", node.Swarm.Nodes); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("managers", node.Swarm.Managers); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("remote_managers", remoteManagers); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(node.Swarm.Cluster.ID)

	return diags
}

func dataSourceCluster() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceClusterRead,
		Schema: map[string]*schema.Schema{
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Computed: true,
			},
			"nodes": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"managers": &schema.Schema{
				Type:     schema.TypeInt,
				Computed: true,
			},
			"remote_managers": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"addr": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
