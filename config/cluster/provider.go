// SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: CC0-1.0

package cluster

import (
	"github.com/upbound/provider-vultr/config/cluster/vke"
)

func init() {
	ProviderConfiguration.AddConfig(vke.Configure)
}
