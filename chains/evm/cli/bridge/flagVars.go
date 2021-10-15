package bridge

import (
	"github.com/ethereum/go-ethereum/common"
)

//flag vars
var (
	Bridge        string
	DataHash      string
	DomainID      uint64
	DepositNonce  uint64
	Handler       string
	ResourceID    string
	Target        string
	Deposit       string
	Execute       string
	Hash          bool
	TokenContract string
)

//processed flag vars
var (
	bridgeAddr         common.Address
	resourceIdBytesArr [32]byte
	handlerAddr        common.Address
	targetContractAddr common.Address
)
