#!/usr/bin/env bash
set -aeuo pipefail

# SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
#
# SPDX-License-Identifier: CC0-1.0

# Delete the ObjectStorageBucket before its ObjectStorage: deleting the
# subscription first removes the bucket on the Vultr side and the bucket then
# wedges on observe (the upstream Read errors instead of clearing state).
${KUBECTL} delete objectstoragebucket.storage.vultr.m.upbound.io --all --all-namespaces
