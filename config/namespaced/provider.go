// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package namespaced

import (
	"github.com/upbound/provider-vultr/config/namespaced/compute"
	"github.com/upbound/provider-vultr/config/namespaced/database"
	"github.com/upbound/provider-vultr/config/namespaced/iam"
	"github.com/upbound/provider-vultr/config/namespaced/network"
	"github.com/upbound/provider-vultr/config/namespaced/storage"
	"github.com/upbound/provider-vultr/config/namespaced/vke"
)

func init() {
	ProviderConfiguration.AddConfig(compute.Configure)
	ProviderConfiguration.AddConfig(database.Configure)
	ProviderConfiguration.AddConfig(iam.Configure)
	ProviderConfiguration.AddConfig(network.Configure)
	ProviderConfiguration.AddConfig(storage.Configure)
	ProviderConfiguration.AddConfig(vke.Configure)
}
