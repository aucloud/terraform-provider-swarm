/*
	terraform-provider-swarm is a Terraform provider for the creation and management of
	Docker Swarm clusters (an alternative container orchestrator to Kubernetes and Nomad)

    Copyright (C) 2021 Sovereign Cloud Australia Pty Ltd

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as published
    by the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.
    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/

package swarm

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/aucloud/go-swarm"
)

func resourceCluster() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceClusterCreate,
		ReadContext:   resourceClusterRead,
		UpdateContext: resourceClusterUpdate,
		DeleteContext: resourceClusterDelete,
		Schema: map[string]*schema.Schema{
			"skip_manager_validation": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"nodes": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"hostname": {
							Type:     schema.TypeString,
							Required: true,
						},
						"public_address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"private_address": {
							Type:     schema.TypeString,
							Required: true,
						},
						"tags": {
							Type:     schema.TypeMap,
							Required: true,
							Elem:     &schema.Schema{Type: schema.TypeString},
						},
					},
				},
			},
			"created_at": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},

			"updated_at": {
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

	force := d.Get("skip_manager_validation").(bool)

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

	managers := vmnodes.FilterByTag(swarm.RoleTag, swarm.ManagerRole)
	if len(managers) == 0 {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "No managers found in cluster config",
			Detail:   "At least one manager must exist in the cluster config! Please check your `nodes` configuration.",
		})
		return diags
	}

	if swarmManager.Runner() == nil {
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
		if force {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Warning,
				Summary:  "Skipping Manager validation",
				Detail: fmt.Sprintf(
					"Forcing creation of %d manager cluster (unsuitable for prod, or ineffective quorum)",
					len(managers),
				),
			})
		}

		if err := swarmManager.CreateSwarm(vmnodes, force); err != nil {
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
