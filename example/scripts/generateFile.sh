#!/bin/bash

source ./config.env

for i in {1..10}; do
    RELAYER_ADDR_VAR="RELAYER_${i}_ADDR"
    RELAYER_PRIVATE_KEY_VAR="RELAYER_${i}_PRIVATE_KEY"
    CONFIG_FILE="config_evm-evm_$i.json"
    cat > $CONFIG_FILE << EOF
{
  "relayer": {
    "opentelemetryCollectorURL": "http://otel-collector:4318"
  },
  "chains": [
    {
      "id": 0,
      "from": "${!RELAYER_ADDR_VAR}",
      "name": "sepolia",
      "type": "evm",
      "endpoint": "https://multi-orbital-vineyard.ethereum-sepolia.discover.quiknode.pro/134970db120a7d265235909912719b6d5c7d96ac/",
      "bridge": "0x9141eBa24846567491e1FB535c02085b33306573",
      "erc20Handler": "0xC5F1dE147460bF7DC4E3d16D5BD489341000dB8D",
      "erc721Handler": "",
      "genericHandler": "",
      "gasLimit": 9000000,
      "maxGasPrice": 20000000000,
      "blockConfirmations": 2,
      "blockInterval": 2,
      "key": "${!RELAYER_PRIVATE_KEY_VAR}"
    },
    {
      "id": 1,
      "from": "${!RELAYER_ADDR_VAR}",
      "name": "fantom-testnet",
      "type": "evm",
      "endpoint": "https://endpoints.omniatech.io/v1/fantom/testnet/public",
      "bridge": "0xEbFBd3Cb62E4AE2b55c31b563964aBf0ab710c84",
      "erc20Handler": "0x4f15f11d4689646012722CB556315705259866B4",
      "erc721Handler": "",
      "genericHandler": "",
      "gasLimit": 9000000,
      "maxGasPrice": 20000000000,
      "blockConfirmations": 2,
      "blockInterval": 2,
      "key": "${!RELAYER_PRIVATE_KEY_VAR}"
    }
  ]
}
EOF
done
