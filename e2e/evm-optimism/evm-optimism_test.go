package main

import (
	"testing"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"

	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/local"
	"github.com/ChainSafe/chainbridge-core/e2e/optimism"
	"github.com/stretchr/testify/suite"
)

const ETHEndpoint1 = "ws://localhost:8646"
const OptimismEndpoint1 = "ws://localhost:8550"
const VerifierEndpoint1 = "ws://localhost:8552"

// Funded optimism address: 0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266
const fundedOptimismPk = "ac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"

// Alice key is used by the relayer, Eve key is used as admin and depositter
func TestRunE2ETests(t *testing.T) {
	kp, err := secp256k1.NewKeypairFromString(fundedOptimismPk)
	if err != nil {
		panic(err)
	}

	suite.Run(t, optimism.SetupEVM2OptimismTestSuite(
		evmtransaction.NewTransaction,
		evmtransaction.NewTransaction,
		ETHEndpoint1,
		OptimismEndpoint1,
		VerifierEndpoint1,
		local.EveKp,
		kp))
}
