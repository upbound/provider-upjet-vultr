#!/usr/bin/env bash
set -aeuo pipefail

# SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
#
# SPDX-License-Identifier: CC0-1.0

# Delete NAT gateway rules before their NATGateway and VPC: once the parents
# are gone the rules wedge on observe and delete (the upstream Read and
# Delete surface the API error instead of clearing state). Firewall rules go
# first — a firewall rule requires its port's forwarding rule to exist.
${KUBECTL} delete natgatewayfirewallrule.network.vultr.m.upbound.io --all --all-namespaces

# Give the API a moment to propagate the firewall rule deletion — a port
# forwarding rule delete fired too early is rejected, yet the rule can
# still vanish server-side afterwards and the CR then wedges on observe.
sleep 10
${KUBECTL} delete natgatewayportforwardingrule.network.vultr.m.upbound.io --all --all-namespaces
