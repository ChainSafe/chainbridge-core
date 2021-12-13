package erc20

import (
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/types"
	"github.com/ethereum/go-ethereum/common"
)

//flag vars
var (
	Amount         string
	Decimals       uint64
	DstAddress     string
	Erc20Address   string
	Recipient      string
	Bridge         string
	DomainID       uint8
	ResourceID     string
	AccountAddress string
	OwnerAddress   string
	SpenderAddress string
	Minter         string
)

//processed flag vars
var (
	recipientAddress   common.Address
	realAmount         *big.Int
	erc20Addr          common.Address
	minterAddr         common.Address
	bridgeAddr         common.Address
	resourceIdBytesArr types.ResourceID
)

// global flags
var (
	dstAddress    common.Address
	url           string
	gasLimit      uint64
	gasPrice      *big.Int
	senderKeyPair *secp256k1.Keypair
)
