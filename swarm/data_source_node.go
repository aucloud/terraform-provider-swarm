package swarm

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceNode() *schema.Resource {
	return &schema.Resource{
		Schema: map[string]*schema.Schema{
			"node": &schema.Schema{
				Type:     schema.TypeMap,
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
								Type:     schema.TypeString,
								Computed: true,
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
						"cpus": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
						"memory": &schema.Schema{
							Type:     schema.TypeInt,
							Computed: true,
						},
					},
				},
			},
		},
	}
}
