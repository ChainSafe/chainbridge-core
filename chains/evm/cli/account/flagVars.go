package account

import (
	"math/big"

	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"

	"github.com/ethereum/go-ethereum/common"
)

//flag vars
var (
	PrivateKey string
	Pass       string
	Recipient  string
	Amount     string
	Decimals   uint64
)

//processed flag vars
var (
	RecipientAddress common.Address
	WeiAmount        *big.Int
)

// global flags
var (
	url           string
	gasLimit      uint64
	gasPrice      *big.Int
	senderKeyPair *secp256k1.Keypair
	prepare       bool
)
