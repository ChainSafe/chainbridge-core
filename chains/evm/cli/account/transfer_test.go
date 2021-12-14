package account

import (
	"testing"

	"github.com/spf13/cobra"
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
	cmd := new(cobra.Command)
	BindTransferBaseCurrencyFlags(cmd)

	err := cmd.Flag("recipient").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateTransferBaseCurrencyFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *TransferTestSuite) TestValidateTransferBaseCurrencyFlagsInvalidAddress() {
	cmd := new(cobra.Command)
	BindTransferBaseCurrencyFlags(cmd)

	err := cmd.Flag("recipient").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateTransferBaseCurrencyFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}
