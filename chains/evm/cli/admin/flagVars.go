package admin

import (
	"math/big"

	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"

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
	BridgeAddr    common.Address
	HandlerAddr   common.Address
	RelayerAddr   common.Address
	RecipientAddr common.Address
	TokenAddr     common.Address
	RealAmount    *big.Int
)

// global flags
var (
	url           string
	gasLimit      uint64
	gasPrice      *big.Int
	senderKeyPair *secp256k1.Keypair
	prepare       bool
)
