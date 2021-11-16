package bridge

import (
	"github.com/ChainSafe/chainbridge-core/types"
	"github.com/ethereum/go-ethereum/common"
)

//flag vars
var (
	Bridge          string
	DataHash        string
	DomainID        uint64
	DepositNonce    uint64
	Handler         string
	ResourceID      string
	Target          string
	Deposit         string
	DepositerOffset uint64
	Execute         string
	Hash            bool
	TokenContract   string
)

//processed flag vars
var (
	bridgeAddr         common.Address
	resourceIdBytesArr types.ResourceID
	handlerAddr        common.Address
	targetContractAddr common.Address
	tokenContractAddr  common.Address
	depositSigBytes    [4]byte
	executeSigBytes    [4]byte
)
