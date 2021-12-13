package account

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

var (
	validAddr   = "0xd606A00c1A39dA53EA7Bb3Ab570BBE40b156EB66"
	invalidAddr = "0xd606A00c1A39dA53EA7Bb3Ab570BBE40b156EXYZ"
)

type TransferTestSuite struct {
	suite.Suite
}

func TestValidateTransferBaseCurrencyFlags(t *testing.T) {
	suite.Run(t, new(TransferTestSuite))
}

func (s *TransferTestSuite) SetupSuite() {
}
func (s *TransferTestSuite) TearDownSuite() {}

func (s *TransferTestSuite) TearDownTest() {}

func (s *TransferTestSuite) TestValidateTransferBaseCurrencyFlags() {
	cmd := transferBaseCurrencyCmd

	cmd.Flag("recipient").Value.Set(validAddr)

	err := ValidateTransferBaseCurrencyFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *TransferTestSuite) TestValidateTransferBaseCurrencyFlagsInvalidAddress() {
	cmd := transferBaseCurrencyCmd

	cmd.Flag("recipient").Value.Set(invalidAddr)

	err := ValidateTransferBaseCurrencyFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}
