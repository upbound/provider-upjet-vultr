// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package compute

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the compute group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("vultr_instance", func(r *config.Resource) {
		r.References["firewall_group_id"] = config.Reference{
			TerraformName: "vultr_firewall_group",
		}
		r.References["iso_id"] = config.Reference{
			TerraformName: "vultr_iso_private",
		}
		r.References["reserved_ip_id"] = config.Reference{
			TerraformName: "vultr_reserved_ip",
		}
		r.References["script_id"] = config.Reference{
			TerraformName: "vultr_startup_script",
		}
		r.References["snapshot_id"] = config.Reference{
			TerraformName: "vultr_snapshot",
		}
		r.References["ssh_key_ids"] = config.Reference{
			TerraformName: "vultr_ssh_key",
		}
		r.References["vpc_ids"] = config.Reference{
			TerraformName: "vultr_vpc",
		}
	})

	p.AddResourceConfigurator("vultr_instance_ipv4", func(r *config.Resource) {
		r.References["instance_id"] = config.Reference{
			TerraformName: "vultr_instance",
		}
	})

	p.AddResourceConfigurator("vultr_bare_metal_server", func(r *config.Resource) {
		r.References["script_id"] = config.Reference{
			TerraformName: "vultr_startup_script",
		}
		r.References["snapshot_id"] = config.Reference{
			TerraformName: "vultr_snapshot",
		}
		r.References["ssh_key_ids"] = config.Reference{
			TerraformName: "vultr_ssh_key",
		}
		r.References["vpc_id"] = config.Reference{
			TerraformName: "vultr_vpc",
		}
	})

	p.AddResourceConfigurator("vultr_inference", func(r *config.Resource) {
		// Upstream leaves the subscription's API key unmarked; keep it out
		// of status and publish it as a connection detail instead.
		r.TerraformResource.Schema["api_key"].Sensitive = true
	})

	p.AddResourceConfigurator("vultr_snapshot", func(r *config.Resource) {
		r.References["instance_id"] = config.Reference{
			TerraformName: "vultr_instance",
		}
	})
}
