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

./bridge evm-cli accounts import --private-key "59c6995e998f97a5a0044966f0945389dc9e86dae88c7a8412f4603b6b78690d" --password "passwordpassword"
./bridge evm-cli accounts import --private-key "5de4111afa1a4b94908f83103eb1f1706367c2e68ca870fc3fb9a804cdab365a" --password "passwordpassword"
./bridge evm-cli accounts import --private-key "7c852118294e51e653712a81e05800f419141751be58f605c371e15141b007a6" --password "passwordpassword"
./bridge evm-cli accounts import --private-key "000000000000000000000000000000000000000000000000000000616c696365" --password "passwordpassword" # Alice
./bridge evm-cli accounts import --private-key "0000000000000000000000000000000000000000000000000000000000626f62" --password "passwordpassword" # Bob
./bridge evm-cli accounts import --private-key "00000000000000000000000000000000000000000000000000636861726c6965" --password "passwordpassword" # Charlie

mv *.key /keys

./bridge run --config /cfg/config_evm-optimism.json --keystore /keys --fresh