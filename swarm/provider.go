package swarm

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"gitlab.mgt.aom.australiacloud.com.au/aom/swarm"
)

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	userLocal := d.Get("use_local").(bool)

	sshAddr := d.Get("ssh_addr").(string)
	sshUser := d.Get("ssh_user").(string)
	sshKey := d.Get("ssh_key").(string)

	var (
		err      error
		manager  *swarm.Manager
		switcher swarm.Switcher
	)

	if userLocal {
		switcher, err = swarm.NewLocalSwitcher()
		if err != nil {
			return nil, diag.FromErr(fmt.Errorf("error creating local switcher: %w", err))
		}

		if err = switcher.Switch(""); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to switch nodes",
				Detail:   "Unable to switch to and connect to local swarm node",
			})
			return nil, diags
		}
	} else {
		switcher, err = swarm.NewSSHSwitcher(sshUser, sshAddr, sshKey)
		if err != nil {
			return nil, diag.FromErr(fmt.Errorf("error creating ssh switcher: %w", err))
		}

		if err = switcher.Switch(sshAddr); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to switch nodes",
				Detail:   "Unable to switch to and connect to remote swarm node",
			})
			return nil, diags
		}
	}

	manager, err = swarm.NewManager(switcher)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("error creating swarm manager: %s", err))
	}

	return manager, diags
}

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"use_local": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("USE_LOCAL", nil),
			},
			"ssh_addr": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SSH_ADDR", nil),
			},
			"ssh_user": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SSH_USER", nil),
			},
			"ssh_key": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("SSH_KEY", nil),
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"swarm_cluster": resourceCluster(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"swarm_cluster": dataSourceCluster(),
			"swarm_nodes":   dataSourceNodes(),
		},
		ConfigureContextFunc: providerConfigure,
	}
}
