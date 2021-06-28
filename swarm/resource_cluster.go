package swarm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.mgt.aom.australiacloud.com.au/aom/swarm"
)

func resourceCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterCreate,
		ReadContext:   resourceClusterRead,
		UpdateContext: resourceClusterUpdate,
		DeleteContext: resourceClusterDelete,
		Schema: map[string]*schema.Schema{
			"nodes": &schema.Schema{
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"public_address": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"private_address": &schema.Schema{
							Type:     schema.TypeString,
							Required: true,
						},
						"tags": &schema.Schema{
							Type:     schema.TypeMap,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
		},
	}
}

func resourceClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	swarmManager := m.(*swarm.Manager)

	nodes := d.Get("nodes").([]interface{})
	vmnodes := make([]swarm.VMNode, len(nodes))

	for i, node := range nodes {
		node := node.(map[string]interface{})
		tags := make(map[string]string)
		for k, v := range node["tags"].(map[string]interface{}) {
			tags[k] = v.(string)
		}

		vmnodes[i] = swarm.VMNode{
			Hostname:       node["hostname"].(string),
			PublicAddress:  node["public_address"].(string),
			PrivateAddress: node["private_address"].(string),
			Tags:           tags,
		}
	}

	if err := swarmManager.CreateSwarm(vmnodes); err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to create swarm cluster",
			Detail:   fmt.Sprintf("Unable to create swarm clusterswitch to: %s", err),
		})
		return diags
	}

	node, err := swarmManager.GetInfo()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting node info: %w", err))
	}

	d.SetId(node.Swarm.Cluster.ID)

	resourceClusterRead(ctx, d, m)

	return diags
}

func resourceClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Same as dataSourceClusterRead
	return dataSourceClusterRead(ctx, d, m)
}

func resourceClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	return resourceClusterRead(ctx, d, m)
}

func resourceClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	return diags
}
