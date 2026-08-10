// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package database

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	"github.com/crossplane/upjet/v2/pkg/config"
)

// Configure configures the database group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("vultr_database", func(r *config.Resource) {
		r.References["vpc_id"] = config.Reference{
			TerraformName: "vultr_vpc",
		}
		// Upstream leaves every credential unmarked; keep them out of
		// status and publish them as connection details instead.
		r.TerraformResource.Schema["password"].Sensitive = true
		r.TerraformResource.Schema["access_key"].Sensitive = true
		r.TerraformResource.Schema["access_cert"].Sensitive = true
		r.TerraformResource.Schema["ferretdb_credentials"].Sensitive = true
		r.TerraformResource.Schema["read_replicas"].Elem.(*schema.Resource).Schema["password"].Sensitive = true
		r.TerraformResource.Schema["read_replicas"].Elem.(*schema.Resource).Schema["ferretdb_credentials"].Sensitive = true
	})

	p.AddResourceConfigurator("vultr_database_connection_pool", func(r *config.Resource) {
		r.References["database_id"] = config.Reference{
			TerraformName: "vultr_database",
		}
		r.References["database"] = config.Reference{
			TerraformName: "vultr_database_db",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("name",false)`,
		}
		r.References["username"] = config.Reference{
			TerraformName: "vultr_database_user",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("username",false)`,
		}
	})

	p.AddResourceConfigurator("vultr_database_connector", func(r *config.Resource) {
		r.References["database_id"] = config.Reference{
			TerraformName: "vultr_database",
		}
	})

	p.AddResourceConfigurator("vultr_database_db", func(r *config.Resource) {
		r.References["database_id"] = config.Reference{
			TerraformName: "vultr_database",
		}
	})

	p.AddResourceConfigurator("vultr_database_quota", func(r *config.Resource) {
		r.References["database_id"] = config.Reference{
			TerraformName: "vultr_database",
		}
		r.References["user"] = config.Reference{
			TerraformName: "vultr_database_user",
			Extractor:     `github.com/crossplane/upjet/v2/pkg/resource.ExtractParamPath("username",false)`,
		}
	})

	p.AddResourceConfigurator("vultr_database_replica", func(r *config.Resource) {
		r.References["database_id"] = config.Reference{
			TerraformName: "vultr_database",
		}
		// Same credential handling as the parent cluster.
		r.TerraformResource.Schema["password"].Sensitive = true
		r.TerraformResource.Schema["ferretdb_credentials"].Sensitive = true
	})

	p.AddResourceConfigurator("vultr_database_topic", func(r *config.Resource) {
		r.References["database_id"] = config.Reference{
			TerraformName: "vultr_database",
		}
	})

	p.AddResourceConfigurator("vultr_database_user", func(r *config.Resource) {
		r.References["database_id"] = config.Reference{
			TerraformName: "vultr_database",
		}
		// Upstream leaves the credentials unmarked; password stays settable
		// through a secret ref, the rest go to connection details.
		r.TerraformResource.Schema["password"].Sensitive = true
		r.TerraformResource.Schema["access_key"].Sensitive = true
		r.TerraformResource.Schema["access_cert"].Sensitive = true
	})
}
