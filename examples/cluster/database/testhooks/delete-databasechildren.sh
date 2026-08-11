#!/usr/bin/env bash
set -aeuo pipefail

# SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
#
# SPDX-License-Identifier: CC0-1.0

# Delete cluster children before the Database: once the cluster is gone they
# wedge on observe and delete (the upstream Reads surface the API error
# instead of clearing state). Quotas go before the users they reference.
${KUBECTL} delete quota.database.vultr.upbound.io --all

${KUBECTL} delete topic.database.vultr.upbound.io --all

${KUBECTL} delete connectionpool.database.vultr.upbound.io --all

${KUBECTL} delete db.database.vultr.upbound.io --all

${KUBECTL} delete replica.database.vultr.upbound.io --all

${KUBECTL} delete user.database.vultr.upbound.io --all
