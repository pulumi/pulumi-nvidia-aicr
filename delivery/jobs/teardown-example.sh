#!/bin/sh
# Destroy the example stack. Safe after a half-finished deploy.
set -eu
. ./delivery/jobs/_example-setup.sh
pulumi env run "$AWS_ENVIRONMENT" -- pulumi destroy --yes
