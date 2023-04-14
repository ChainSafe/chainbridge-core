#!/bin/bash

# change directory to parent folder
cd ..

# source config.env file
source config.env

# run command multiple times
for i in {1..50}; do
  ./bridge-cli evm-cli erc20 deposit \
  --url $SRC_GATEWAY \
  --private-key $USER_1_PRIVATE_KEY \
  --gas-price 25000000000 \
  --amount 0.0000000000000001 \
  --domain 1 \
  --bridge $SRC_BRIDGE \
  --recipient $USER_1_ADDR \
  --resource $RESOURCE_ID \
  --decimals 18
done