package erc20

import (
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
	DomainID       uint64
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
