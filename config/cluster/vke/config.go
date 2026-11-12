// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package vke

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the vke group
func Configure(p *config.Provider) {
	// ClusterAuth is a hand-written resource: the Terraform provider no
	// longer exposes the kubeconfig on vultr_kubernetes, so a dedicated
	// controller fetches it from the Vultr API and publishes it.
	p.BasePackages.ControllerMap["vke/clusterauth"] = "vke"

	p.AddResourceConfigurator("vultr_kubernetes", func(r *config.Resource) {
		r.UseAsync = true
		r.References["vpc_id"] = config.Reference{
			TerraformName: "vultr_vpc",
		}
		// The autoscaler changes node_quantity at runtime; late-initializing
		// it into spec would make the provider fight the autoscaler.
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"node_pools.node_quantity"},
		}
	})

	p.AddResourceConfigurator("vultr_kubernetes_node_pools", func(r *config.Resource) {
		r.UseAsync = true
		r.References["cluster_id"] = config.Reference{
			TerraformName: "vultr_kubernetes",
		}
		// The autoscaler changes node_quantity at runtime; late-initializing
		// it into spec would make the provider fight the autoscaler.
		r.LateInitializer = config.LateInitializer{
			IgnoredFields: []string{"node_quantity"},
		}
	})
}
