// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package iam

import (
	"time"

	"github.com/crossplane/upjet/v2/pkg/config"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// Configure configures the iam group
func Configure(p *config.Provider) {
	p.AddResourceConfigurator("vultr_oidc_issuer", func(r *config.Resource) {
		// source is one of "external" or "vke"; source_id carries the VKE
		// cluster ID when source is "vke".
		r.References["source_id"] = config.Reference{TerraformName: "vultr_kubernetes"}
	})

	p.AddResourceConfigurator("vultr_oidc_token", func(r *config.Resource) {
		// Upstream leaves the issued tokens unmarked; keep them out of
		// status and publish them as connection details instead.
		r.TerraformResource.Schema["access_token"].Sensitive = true
		r.TerraformResource.Schema["id_token"].Sensitive = true
		r.TerraformResource.Schema["refresh_token"].Sensitive = true
	})

	p.AddResourceConfigurator("vultr_organization_policy_group_attachment", func(r *config.Resource) {
		r.References["policy_id"] = config.Reference{TerraformName: "vultr_organization_policy"}
		r.References["group_id"] = config.Reference{TerraformName: "vultr_organization_group"}
	})

	p.AddResourceConfigurator("vultr_organization_policy_user_attachment", func(r *config.Resource) {
		r.References["policy_id"] = config.Reference{TerraformName: "vultr_organization_policy"}
		r.References["user_id"] = config.Reference{TerraformName: "vultr_user"}
	})

	p.AddResourceConfigurator("vultr_organization_role_group_attachment", func(r *config.Resource) {
		r.References["role_id"] = config.Reference{TerraformName: "vultr_organization_role"}
		r.References["group_id"] = config.Reference{TerraformName: "vultr_organization_group"}
	})

	p.AddResourceConfigurator("vultr_organization_role_policy_attachment", func(r *config.Resource) {
		r.References["role_id"] = config.Reference{TerraformName: "vultr_organization_role"}
		r.References["policy_id"] = config.Reference{TerraformName: "vultr_organization_policy"}
	})

	p.AddResourceConfigurator("vultr_organization_role_session", func(r *config.Resource) {
		// Upstream leaves the session token unmarked; keep it out of
		// status and publish it as a connection detail instead.
		r.TerraformResource.Schema["token"].Sensitive = true
		r.References["role_id"] = config.Reference{TerraformName: "vultr_organization_role"}
		r.References["user_id"] = config.Reference{TerraformName: "vultr_user"}
	})

	p.AddResourceConfigurator("vultr_organization_role_trust", func(r *config.Resource) {
		r.References["role"] = config.Reference{TerraformName: "vultr_organization_role"}
		r.References["group"] = config.Reference{TerraformName: "vultr_organization_group"}
		r.References["user"] = config.Reference{TerraformName: "vultr_user"}
		// The API rewrites the time-of-day of date_expires depending on
		// the caller, so only the calendar date round-trips stably;
		// comparing the full timestamp leaves a permanent diff
		r.TerraformResource.Schema["date_expires"].DiffSuppressFunc = func(_, old, new string, _ *schema.ResourceData) bool {
			oldT, oldErr := time.Parse(time.RFC3339, old)
			newT, newErr := time.Parse(time.RFC3339, new)
			if oldErr != nil || newErr != nil {
				return false
			}
			oy, om, od := oldT.Date()
			ny, nm, nd := newT.Date()
			return oy == ny && om == nm && od == nd
		}
	})
}
