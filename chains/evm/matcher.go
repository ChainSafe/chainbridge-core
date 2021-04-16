package evm

import (
	"github.com/ChainSafe/chainbridgev2/bindings/eth/bindings/Bridge"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
)

type EVMResourceIDMatcher struct {
	bridgeContractAddress string
	chainReader           ChainClient
}

func (l *EVMResourceIDMatcher) MatchResourceIDToHandlerAddress(rID [32]byte) (string, error) {
	b, err := Bridge.NewBridgeCaller(common.HexToAddress(l.bridgeContractAddress), l.chainReader)
	if err != nil {
		return "", err
	}
	addr, err := b.ResourceIDToHandlerAddress(&bind.CallOpts{}, rID)
	if err != nil {
		return "", err
	}
	return addr.String(), nil
}
