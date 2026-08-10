// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package vke

import (
	"encoding/base64"

	"github.com/crossplane/crossplane-runtime/v2/pkg/errors"
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the vke group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("vultr_kubernetes", func(r *config.Resource) {
		r.UseAsync = true
		// The kube_config attribute is base64-encoded; publish the decoded
		// kubeconfig so consumers can mount it directly.
		r.Sensitive.AdditionalConnectionDetailsFn = func(attr map[string]any) (map[string][]byte, error) {
			conn := map[string][]byte{}
			if kc, ok := attr["kube_config"].(string); ok && kc != "" {
				decoded, err := base64.StdEncoding.DecodeString(kc)
				if err != nil {
					return nil, errors.Wrap(err, "cannot base64-decode the kube_config attribute")
				}
				conn["kubeconfig"] = decoded
			}
			return conn, nil
		}
	})

	p.AddResourceConfigurator("vultr_kubernetes_node_pools", func(r *config.Resource) {
		r.UseAsync = true
		r.References["cluster_id"] = config.Reference{
			TerraformName: "vultr_kubernetes",
		}
	})
}
