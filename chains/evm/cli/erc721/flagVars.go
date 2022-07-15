package erc721

import (
	"math/big"

	"github.com/ChainSafe/sygma-core/crypto/secp256k1"
	"github.com/ChainSafe/sygma-core/types"
	"github.com/ethereum/go-ethereum/common"
)

// flag vars
var (
	Erc721Address  string
	Dst            string
	Token          string
	Metadata       string
	Recipient      string
	Bridge         string
	DestionationID string
	ResourceID     string
	Minter         string
	Priority       string
)

// processed flag vars
var (
	Erc721Addr    common.Address
	DstAddress    common.Address
	TokenId       *big.Int
	RecipientAddr common.Address
	BridgeAddr    common.Address
	DestinationID int
	ResourceId    types.ResourceID
	MinterAddr    common.Address
)

// global flags
var (
	url           string
	gasLimit      uint64
	gasPrice      *big.Int
	senderKeyPair *secp256k1.Keypair
	prepare       bool
	err           error
)
