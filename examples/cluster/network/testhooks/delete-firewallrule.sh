#!/usr/bin/env bash
set -aeuo pipefail

# SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
#
# SPDX-License-Identifier: CC0-1.0

# Delete the FirewallRule before the FirewallGroup: deleting the group first
# cascade-removes the rule on the Vultr side and the rule then wedges on
# observe (the upstream Read does not recognize the resulting 404).
${KUBECTL} delete firewallrule.network.vultr.upbound.io --all
