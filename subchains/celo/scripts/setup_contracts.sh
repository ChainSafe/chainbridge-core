#!/usr/bin/env bash
# Copyright 2020 ChainSafe Systems
# SPDX-License-Identifier: LGPL-3.0-only

CONTRACTS_REPO="https://github.com/ChainSafe/chainbridge-celo-solidity/"
CONTRACTS_BRANCH="main"
CONTRACTS_COMMIT="70c97947f6b05397c2887fe8213a3d49d3d06de7"
CONTRACTS_DIR="./solidity"
DEST_DIR="./bindings"

set -eux

case $TARGET in
	"build")
		git clone -b $CONTRACTS_BRANCH $CONTRACTS_REPO $CONTRACTS_DIR
    pushd $CONTRACTS_DIR
    git checkout $CONTRACTS_COMMIT

    make install-deps
    make bindings

    popd

    mkdir $DEST_DIR
    cp -r $CONTRACTS_DIR/build/bindings/go/* $DEST_DIR
		;;

	"cli-only")
		git clone -b $CONTRACTS_BRANCH $CONTRACTS_REPO $CONTRACTS_DIR
    pushd $CONTRACTS_DIR
    git checkout $CONTRACTS_COMMIT
		;;

esac
