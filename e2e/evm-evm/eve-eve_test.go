package main

import (
	"testing"

	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/local"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/e2e/evm"
	"github.com/stretchr/testify/suite"
)

// Alice key is used by the relayer, Eve key is used as admin and depositter
func TestRunE2ETests(t *testing.T) {
	suite.Run(t, evm.SetupEVM2EVMTestSuite(evmtransaction.NewTransaction, evmtransaction.NewTransaction, ETHEndpoint1, ETHEndpoint2, local.EveKp))
}
