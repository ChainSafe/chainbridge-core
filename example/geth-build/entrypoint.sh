#!/usr/bin/env sh
# Copyright 2020 ChainSafe Systems
# SPDX-License-Identifier: LGPL-3.0-only

# Exit on failure
set -ex

geth init /root/genesis.json
rm -f /root/.ethereum/keystore/*

# If accounts are not set, set all accounts.
if [ -z $ACCOUNTS ]
then
  ACCOUNTS="0xf4314cb9046bece6aa54bb9533155434d0c76909,0xff93B45308FD417dF303D6515aB04D9e89a750Ca,0x8e0a907331554AF72563Bd8D43051C2E64Be5d35,0x24962717f8fA5BA3b931bACaF9ac03924EB475a0,0x148FfB2074A9e59eD58142822b3eB3fcBffb0cd7,0x4CEEf6139f00F9F4535Ad19640Ff7A0137708485"
fi

# Copy requested accounts to celo keystore.
ACCOUNTS_TRIMMED=${ACCOUNTS//0x/}
ACCOUNTS_PATTERN=${ACCOUNTS_TRIMMED//,/|}
find /root/keystore | grep -iE ${ACCOUNTS_PATTERN} | xargs -i cp {} /root/.ethereum/keystore/

# Identify the docker container external IP.
IP=$(ip -4 -o address | \
  grep -Eo -m 1 'eth0\s+inet\s+[0-9]{1,3}[.][0-9]{1,3}[.][0-9]{1,3}[.][0-9]{1,3}' | \
  grep -Eo '[0-9]{1,3}[.][0-9]{1,3}[.][0-9]{1,3}[.][0-9]{1,3}')

if [ ! -z $BOOTNODE ]
then
  BOOTNODE="--bootnodes ${BOOTNODE}"
else
  BOOTNODE=" "
fi
if [ ! -z $NODEKEY ]; then NODEKEY="--nodekeyhex ${NODEKEY}"; fi
if [ -z $NETWORKID ]; then NETWORKID="5"; fi

if [ ! -z $MINE ];
then
#  MINE="--mine --miner.etherbase=0x0000000000000000000000000000000000000000 --miner.threads=1"
    MINE="--mine"
else
  MINE=" "
fi

exec geth \
  --unlock ${ACCOUNTS} \
  --password /root/password.txt \
  --ws \
  --ws.port 8546 \
  --ws.origins="*" \
  --ws.addr 0.0.0.0 \
  --http \
  --http.addr 0.0.0.0\
  --http.port 8545 \
  --http.corsdomain="*" \
  --http.vhosts="*" \
  --nat=extip:${IP} \
  --networkid ${NETWORKID} \
  --allow-insecure-unlock ${MINE} ${BOOTNODE} ${NODEKEY}

#  HTTP based JSON-RPC API options:
#
#--http Enable the HTTP-RPC server
#--http.addr HTTP-RPC server listening interface (default: localhost)
#--http.port HTTP-RPC server listening port (default: 8545)
#--http.api API's offered over the HTTP-RPC interface (default: eth,net,web3)
#--http.corsdomain Comma separated list of domains from which to accept cross origin requests (browser enforced)
#--ws Enable the WS-RPC server
#--ws.addr WS-RPC server listening interface (default: localhost)
#--ws.port WS-RPC server listening port (default: 8546)
#--ws.api API's offered over the WS-RPC interface (default: eth,net,web3)
#--ws.origins Origins from which to accept websockets requests
#--ipcdisable Disable the IPC-RPC server
#--ipcapi API's offered over the IPC-RPC interface (default: admin,debug,eth,miner,net,personal,shh,txpool,web3)
#--ipcpath Filename for IPC socket/pipe within the datadir (explicit paths escape it)
