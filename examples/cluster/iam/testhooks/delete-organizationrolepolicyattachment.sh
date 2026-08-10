#!/usr/bin/env bash
set -aeuo pipefail
# Delete the attachment before its parent resources: deleting a parent first
# removes the attachment on the Vultr side and the attachment then wedges on
# delete (the upstream Read does not clear state when it is gone).
${KUBECTL} delete organizationrolepolicyattachment.iam.vultr.upbound.io --all
