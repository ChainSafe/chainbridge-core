package e2e_test

import (
	"testing"

	"github.com/ChainSafe/chainbridge-core/chains/evm/e2e"
	"github.com/stretchr/testify/suite"
)

func TestRunE2ETests(t *testing.T) {
	suite.Run(t, new(e2e.IntegrationTestSuite))
}
