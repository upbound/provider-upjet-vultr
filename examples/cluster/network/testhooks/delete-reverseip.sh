#!/usr/bin/env bash
set -aeuo pipefail

# SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
#
# SPDX-License-Identifier: CC0-1.0

# Delete reverse DNS records before their Instance: once the instance is
# gone the records wedge on observe and delete (the upstream Read and
# Delete surface the API error instead of clearing state).
${KUBECTL} delete reverseipv4.network.vultr.upbound.io --all

${KUBECTL} delete reverseipv6.network.vultr.upbound.io --all
