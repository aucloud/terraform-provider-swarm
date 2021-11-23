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
	"strconv"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/aucloud/go-swarm"
)

func dataSourceNodesRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	swarmManager := m.(*swarm.Manager)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	nodes := make([]map[string]interface{}, 0)

	node, err := swarmManager.GetInfo()
	if err != nil {
		return diag.FromErr(fmt.Errorf("error getting node info: %w", err))
	}
	nodes = append(nodes, map[string]interface{}{
		"id":             node.ID,
		"name":           node.Name,
		"labels":         node.Labels,
		"cpus":           node.NCPU,
		"memory":         node.MemTotal,
		"os":             node.OperatingSystem,
		"os_type":        node.OSType,
		"os_version":     node.OSVersion,
		"kernel_version": node.KernelVersion,
		"server_version": node.ServerVersion,
		"manager":        node.IsManager(),
	})

	if err := d.Set("all", nodes); err != nil {
		return diag.FromErr(err)
	}

	// always run
	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func dataSourceNodes() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceNodesRead,
		Schema: map[string]*schema.Schema{
			"all": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"labels": {
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"os": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_type": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"kernel_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_version": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpus": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory": {
							Type:     schema.TypeInt,
							Computed: true,
						},
						"manager": {
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
