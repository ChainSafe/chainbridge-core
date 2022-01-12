#!/bin/bash

set -e

RETRIES=${RETRIES:-40}

# waits for l2geth to be up
curl --fail \
    --show-error \
    --silent \
    --retry-connrefused \
    --retry $RETRIES \
    --retry-delay 1 \
    --output /dev/null \
    $L2_NODE_WEB3_URL

./bridge run --config /cfg/config_evm-optimism.json --keystore /keys --fresh