#!/bin/sh
# Plan the example against the freshly built provider. Creates nothing.
set -eu
. ./delivery/jobs/_example-setup.sh
pulumi env run "$AWS_ENVIRONMENT" -- pulumi preview --diff
