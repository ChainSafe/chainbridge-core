package account

import (
	"math/big"

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
	recipientAddress common.Address
	weiAmount        *big.Int
)
