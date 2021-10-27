package admin

import (
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
	bridgeAddr  common.Address
	relayerAddr common.Address
)
