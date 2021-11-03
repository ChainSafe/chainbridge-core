package main

import (
	"testing"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/e2e/evm"
	"github.com/stretchr/testify/suite"
)

const ETHEndpoint1 = "http://localhost:8545"
const ETHEndpoint2 = "http://localhost:8547"

// Alice key is used by the relayer, Eve key is used as admin and depositter
func TestRunE2ETests(t *testing.T) {
	suite.Run(t, evm.PreSetupTestSuite(evmtransaction.NewTransaction, evmtransaction.NewTransaction, ETHEndpoint1, ETHEndpoint2, evm.EveKp))
}
