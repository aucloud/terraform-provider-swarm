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
