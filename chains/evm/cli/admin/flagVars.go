package admin

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

//flag vars
var (
	Admin            string
	Relayer          string
	DepositNonce     uint64
	DomainID         uint8
	Fee              string
	RelayerThreshold uint64
	Amount           string
	TokenID          string
	Handler          string
	Token            string
	Decimals         uint64
	Recipient        string
	Bridge           string
)

//processed flag vars
var (
	bridgeAddr    common.Address
	handlerAddr   common.Address
	relayerAddr   common.Address
	recipientAddr common.Address
	tokenAddr     common.Address
	realAmount    *big.Int
)
