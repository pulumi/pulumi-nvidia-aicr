#!/bin/sh
# The real deploy. Gated upstream -- see the deploy stage in ../index.ts.
set -eu
. ./delivery/jobs/_example-setup.sh
pulumi env run "$AWS_ENVIRONMENT" -- pulumi up --yes
[ -n "${DELIVERY_OUTPUTS:-}" ] && pulumi stack output --json > "$DELIVERY_OUTPUTS"
