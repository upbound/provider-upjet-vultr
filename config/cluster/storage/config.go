// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package storage

import (
	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the storage group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("vultr_object_storage_bucket", func(r *config.Resource) {
		r.References["object_storage_id"] = config.Reference{
			TerraformName: "vultr_object_storage",
		}
	})
}
