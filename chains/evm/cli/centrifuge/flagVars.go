package centrifuge

import (
	"math/big"

	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/common"
)

//flag vars
var (
	Hash    string
	Address string
)

//processed flag vars
var (
	StoreAddr common.Address
	ByteHash  [32]byte
)

// global flags
var (
	url           string
	gasPrice      *big.Int
	senderKeyPair *secp256k1.Keypair
	prepare       bool
)
