#!/bin/sh
# Sourced by preview/deploy/teardown. Builds the provider and its Node SDK from
# this checkout and points the example at them, because the pulumi-org packages
# the examples depend on are not published to npm.
set -eu
: "${EXAMPLE_DIR:?EXAMPLE_DIR is not set}"
: "${TARGET_STACK:?TARGET_STACK is not set}"

echo "=== build provider"
make provider

echo "=== build node SDK"
make build_nodejs_sdk

echo "=== install provider plugin"
VERSION=$(./bin/pulumi-resource-nvidia-aicr --version 2>/dev/null | head -1 || echo 0.0.1-dev)
pulumi plugin install resource nvidia-aicr "$VERSION" \
    --file ./bin/pulumi-resource-nvidia-aicr --reinstall

REPO_ROOT=$(pwd)
cd "$EXAMPLE_DIR"

echo "=== wire example to the locally built SDK"
# The registry copy does not exist, so let this fail and overwrite it next.
npm install --no-audit --no-fund --ignore-scripts >/dev/null 2>&1 || true
npm install --no-save --no-audit --no-fund "$REPO_ROOT/sdk/nodejs/bin"

pulumi stack select "$TARGET_STACK"
