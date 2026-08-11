#!/usr/bin/env bash
set -aeuo pipefail

# SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
#
# SPDX-License-Identifier: CC0-1.0

# Delete cluster children before the Database: once the cluster is gone they
# wedge on observe and delete (the upstream Reads surface the API error
# instead of clearing state). Quotas go before the users they reference.
${KUBECTL} delete quota.database.vultr.m.upbound.io --all --all-namespaces

${KUBECTL} delete topic.database.vultr.m.upbound.io --all --all-namespaces

${KUBECTL} delete connectionpool.database.vultr.m.upbound.io --all --all-namespaces

${KUBECTL} delete db.database.vultr.m.upbound.io --all --all-namespaces

${KUBECTL} delete replica.database.vultr.m.upbound.io --all --all-namespaces

${KUBECTL} delete user.database.vultr.m.upbound.io --all --all-namespaces
