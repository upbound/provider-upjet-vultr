// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package network

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the network group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("vultr_dns_record", func(r *config.Resource) {
		r.References["domain"] = config.Reference{
			TerraformName: "vultr_dns_domain",
		}
	})

	p.AddResourceConfigurator("vultr_firewall_rule", func(r *config.Resource) {
		r.References["firewall_group_id"] = config.Reference{
			TerraformName: "vultr_firewall_group",
		}
	})

	p.AddResourceConfigurator("vultr_load_balancer", func(r *config.Resource) {
		r.References["vpc"] = config.Reference{
			TerraformName: "vultr_vpc",
		}
	})

	p.AddResourceConfigurator("vultr_nat_gateway", func(r *config.Resource) {
		r.References["vpc_id"] = config.Reference{
			TerraformName: "vultr_vpc",
		}
	})

	p.AddResourceConfigurator("vultr_nat_gateway_firewall_rule", func(r *config.Resource) {
		r.References["vpc_id"] = config.Reference{
			TerraformName: "vultr_vpc",
		}
		r.References["nat_gateway_id"] = config.Reference{
			TerraformName: "vultr_nat_gateway",
		}
	})

	p.AddResourceConfigurator("vultr_nat_gateway_port_forwarding_rule", func(r *config.Resource) {
		r.References["vpc_id"] = config.Reference{
			TerraformName: "vultr_vpc",
		}
		r.References["nat_gateway_id"] = config.Reference{
			TerraformName: "vultr_nat_gateway",
		}
	})

	p.AddResourceConfigurator("vultr_reverse_ipv4", func(r *config.Resource) {
		r.References["instance_id"] = config.Reference{
			TerraformName: "vultr_instance",
		}
		// The record's ip must be an address of the referenced instance;
		// resolve it from the instance's observed main_ip.
		r.References["ip"] = config.Reference{
			TerraformName: "vultr_instance",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("main_ip",true)`,
		}
	})

	p.AddResourceConfigurator("vultr_reverse_ipv6", func(r *config.Resource) {
		r.References["instance_id"] = config.Reference{
			TerraformName: "vultr_instance",
		}
		// The record's ip must be an address of the referenced instance;
		// resolve it from the instance's observed v6_main_ip.
		r.References["ip"] = config.Reference{
			TerraformName: "vultr_instance",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("v6_main_ip",true)`,
		}
	})
}
