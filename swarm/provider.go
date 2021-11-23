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

func providerConfigure(ctx context.Context, d *schema.ResourceData) (interface{}, diag.Diagnostics) {
	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	useLocal := d.Get("use_local").(bool)

	sshTimeout := d.Get("ssh_timeout").(string)
	timeout, err := time.ParseDuration(sshTimeout)
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("error parsing ssh timeout %s: %w", sshTimeout, err))
	}

	sshAddr := d.Get("ssh_addr").(string)
	sshUser := d.Get("ssh_user").(string)
	sshKey := d.Get("ssh_key").(string)

	var (
		manager  *swarm.Manager
		switcher swarm.Switcher
	)

	if useLocal {
		switcher, err = swarm.NewLocalSwitcher()
		if err != nil {
			return nil, diag.FromErr(fmt.Errorf("error creating local switcher: %w", err))
		}

		ctx := context.Background()
		if err = switcher.Switch(ctx, ""); err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to switch nodes",
				Detail: fmt.Sprintf(
					"Unable to switch to and connect to local swarm node: %s",
					err.Error(),
				),
			})
			return nil, diags
		}
	} else {
		switcher, err = swarm.NewSSHSwitcher(sshUser, sshAddr, sshKey, timeout)
		if err != nil {
			return nil, diag.FromErr(fmt.Errorf("error creating ssh switcher: %w", err))
		}

		if sshAddr != "" {
			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()
			if err = switcher.Switch(ctx, sshAddr); err != nil {
				diags = append(diags, diag.Diagnostic{
					Severity: diag.Error,
					Summary:  "Unable to switch nodes",
					Detail: fmt.Sprintf(
						"Unable to switch to and connect to remote swarm node: %s",
						err.Error(),
					),
				})
				return nil, diags
			}
		}
	}

	manager, err = swarm.NewManager(switcher, swarm.WithTimeout(timeout))
	if err != nil {
		return nil, diag.FromErr(fmt.Errorf("error creating swarm manager: %s", err))
	}

	return manager, diags
}

// Provider -
func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"use_local": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("USE_LOCAL", nil),
			},
			"ssh_timeout": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "5m",
				DefaultFunc: schema.EnvDefaultFunc("SSH_TIMEOUT", nil),
			},
			"ssh_addr": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SSH_ADDR", nil),
			},
			"ssh_user": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("SSH_USER", nil),
			},
			"ssh_key": {
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
