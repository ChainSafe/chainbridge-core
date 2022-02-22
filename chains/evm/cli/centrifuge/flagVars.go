package centrifuge

import (
	"math/big"

	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/ChainSafe/chainbridge-core/types"
	"github.com/ethereum/go-ethereum/common"
)

//flag vars
var (
	Hash       string
	Address    string
	Metadata   string
	Recipient  string
	Bridge     string
	DomainID   uint8
	ResourceID string
	Priority   string
)

//processed flag vars
var (
	StoreAddr          common.Address
	ByteHash           [32]byte
	MetadataBytes      []byte
	RecipientAddr      common.Address
	BridgeAddr         common.Address
	ResourceIdBytesArr types.ResourceID
)

// global flags
var (
	url           string
	gasPrice      *big.Int
	gasLimit      uint64
	senderKeyPair *secp256k1.Keypair
	prepare       bool
)
