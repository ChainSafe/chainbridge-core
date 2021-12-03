package centrifuge

import (
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/common"
	"math/big"
)

//flag vars
var (
	Hash    string
	Address string
)

//processed flag vars
var (
	storeAddr common.Address
	byteHash  [32]byte
)

// global flags
var (
	dstAddress    common.Address
	url           string
	gasLimit      uint64
	gasPrice      *big.Int
	senderKeyPair *secp256k1.Keypair
)
