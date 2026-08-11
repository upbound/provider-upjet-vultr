#!/usr/bin/env bash
set -aeuo pipefail

# SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
#
# SPDX-License-Identifier: CC0-1.0

# Delete the InstanceIPv4 before its Instance: once the instance is gone the
# additional IP wedges on observe and delete (the upstream Read and Delete
# surface the API error instead of clearing state).
${KUBECTL} delete instanceipv4.compute.vultr.m.upbound.io --all --all-namespaces
