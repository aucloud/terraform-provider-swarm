package swarm

import (
	"context"
	"fmt"
	"time"

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
			"force_single_manager_cluster": &schema.Schema{
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
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
			"created_at": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"updated_at": &schema.Schema{
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
	}
}

func resourceClusterCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	swarmManager := m.(*swarm.Manager)

	forceSingleManagerCluster := d.Get("force_single_manager_cluster").(bool)

	nodes := d.Get("nodes").([]interface{})
	vmnodes := make(swarm.VMNodes, len(nodes))

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

	if swarmManager.Runner() == nil {
		managers := vmnodes.FilterByTag(swarm.RoleTag, swarm.ManagerRole)

		if err := swarmManager.SwitchNode(managers[0].PublicAddress); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to switch to first manager node",
				Detail: fmt.Sprintf(
					"Error switching to first manager node %s via %s: %s",
					managers[0].Hostname, managers[0].PublicAddress, err.Error(),
				),
			})
			return diags
		}
	}

	node, err := swarmManager.GetInfo()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to retrieve node information",
			Detail: fmt.Sprintf(
				"Error getting node info from %s: %s",
				swarmManager.Switcher().String(), err.Error(),
			),
		})
		return diags
	}

	if node.Swarm.Cluster.ID == "" {
		if forceSingleManagerCluster {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Forcing Single Manager Cluster",
				Detail:   "Forcing single manager clsuter (unsuitable for prod)",
			})
		}

		if err := swarmManager.CreateSwarm(vmnodes, forceSingleManagerCluster); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to create swarm cluster",
				Detail:   fmt.Sprintf("Unable to create swarm cluster: %s", err),
			})
			return diags
		}

		node, err = swarmManager.GetInfo()
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to retrieve cluster information",
				Detail: fmt.Sprintf(
					"Error getting node info from %s: %s",
					swarmManager.Switcher().String(), err.Error(),
				),
			})
			return diags
		}

		if err := d.Set("created_at", time.Now().Format(time.RFC3339)); err != nil {
			return diag.FromErr(err)
		}
	}

	d.SetId(node.Swarm.Cluster.ID)

	resourceClusterRead(ctx, d, m)

	return diags
}

func resourceClusterRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	swarmManager := m.(*swarm.Manager)

	nodes := d.Get("nodes").([]interface{})
	vmnodes := make(swarm.VMNodes, len(nodes))

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

	if swarmManager.Runner() == nil {
		managers := vmnodes.FilterByTag(swarm.RoleTag, swarm.ManagerRole)

		if err := swarmManager.SwitchNode(managers[0].PublicAddress); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Unable to switch to first manager node",
				Detail: fmt.Sprintf(
					"Error switching to first manager node %s via %s: %s",
					managers[0].Hostname, managers[0].PublicAddress, err.Error(),
				),
			})

			// TODO: Really need to see if we can figure our a more reliable
			//       way to identity whether the underlying machines on which
			//       the swarm cluster was formed are **really** gone or not
			//       instead of basically assuming here based on an error.
			d.SetId("")
			return diags
		}
	}

	node, err := swarmManager.GetInfo()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to retrieve node information",
			Detail: fmt.Sprintf(
				"Error getting node info from %s: %s",
				swarmManager.Switcher().String(), err.Error(),
			),
		})
		return diags
	}

	d.SetId(node.Swarm.Cluster.ID)

	dataSourceClusterRead(ctx, d, m)

	return diags
}

func resourceClusterUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	swarmManager := m.(*swarm.Manager)

	if d.HasChange("nodes") {
		nodes := d.Get("nodes").([]interface{})
		vmnodes := make(swarm.VMNodes, len(nodes))

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

		if swarmManager.Runner() == nil {
			managers := vmnodes.FilterByTag(swarm.RoleTag, swarm.ManagerRole)

			if err := swarmManager.SwitchNode(managers[0].PublicAddress); err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Unable to switch to first manager node",
					Detail: fmt.Sprintf(
						"Error switching to first manager node %s via %s: %s",
						managers[0].Hostname, managers[0].PublicAddress, err.Error(),
					),
				})
				return diags
			}
		}

		if err := swarmManager.UpdateSwarm(vmnodes); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to update swarm cluster",
				Detail:   fmt.Sprintf("Unable to update swarm cluster: %s", err),
			})
			return diags
		}

		if err := d.Set("updated_at", time.Now().Format(time.RFC3339)); err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceClusterRead(ctx, d, m)
}

func resourceClusterDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	diags = append(diags, diag.Diagnostic{
		Severity: diag.Warning,
		Summary:  "Swarm Cluster Destruction NOT IMPLEMENTED",
		Detail:   "Swarm Cluster Destruction is not yet implemented",
	})

	return diags
}
