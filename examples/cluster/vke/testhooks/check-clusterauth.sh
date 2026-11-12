#!/usr/bin/env bash
set -aeuo pipefail

# SPDX-FileCopyrightText: 2026 The Crossplane Authors <https://crossplane.io>
#
# SPDX-License-Identifier: CC0-1.0

# Verify that the kubeconfig published by the ClusterAuth grants access to the
# VKE cluster it was fetched for.
KUBECONFIG_FILE="$(mktemp)"
trap 'rm -f "${KUBECONFIG_FILE}"' EXIT

echo "obtain kubeconfig from ClusterAuth connection secret"
${KUBECTL} -n upbound-system get secret example-vke-auth-kubeconfig -o go-template='{{ .data.kubeconfig | base64decode }}' > "${KUBECONFIG_FILE}"
echo "listing nodes"
${KUBECTL} --kubeconfig "${KUBECONFIG_FILE}" get nodes
