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
			"all": &schema.Schema{
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"name": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"labels": &schema.Schema{
							Type:     schema.TypeList,
							Computed: true,
							Elem: &schema.Schema{
								Type: schema.TypeString,
							},
						},
						"os": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_type": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"os_version": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"kernel_version": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"server_version": &schema.Schema{
							Type:     schema.TypeString,
							Computed: true,
						},
						"cpus": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"manager": &schema.Schema{
							Type:     schema.TypeBool,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
