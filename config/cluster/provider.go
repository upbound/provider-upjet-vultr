// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package cluster

import (
	"github.com/upbound/provider-vultr/config/cluster/compute"
	"github.com/upbound/provider-vultr/config/cluster/database"
	"github.com/upbound/provider-vultr/config/cluster/iam"
	"github.com/upbound/provider-vultr/config/cluster/network"
	"github.com/upbound/provider-vultr/config/cluster/storage"
	"github.com/upbound/provider-vultr/config/cluster/vke"
)

func init() {
	ProviderConfiguration.AddConfig(compute.Configure)
	ProviderConfiguration.AddConfig(database.Configure)
	ProviderConfiguration.AddConfig(iam.Configure)
	ProviderConfiguration.AddConfig(network.Configure)
	ProviderConfiguration.AddConfig(storage.Configure)
	ProviderConfiguration.AddConfig(vke.Configure)
}
