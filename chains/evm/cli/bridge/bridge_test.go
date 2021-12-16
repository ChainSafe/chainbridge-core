package bridge

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

var (
	validAddr   = "0xd606A00c1A39dA53EA7Bb3Ab570BBE40b156EB66"
	invalidAddr = "0xd606A00c1A39dA53EA7Bb3Ab570BBE40b156EXYZ"
)

type BridgeTestSuite struct {
	suite.Suite
}

func TestBridgeTestSuite(t *testing.T) {
	suite.Run(t, new(BridgeTestSuite))
}

func (s *BridgeTestSuite) SetupSuite() {
}
func (s *BridgeTestSuite) TearDownSuite() {}

func (s *BridgeTestSuite) TearDownTest() {}

func (s *BridgeTestSuite) TestValidateCancelProposalFlags() {
	cmd := new(cobra.Command)
	BindCancelProposalFlags(cmd)

	err := cmd.Flag("bridge").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateCancelProposalFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *BridgeTestSuite) TestValidateCancelProposalInvalidAddress() {
	cmd := new(cobra.Command)
	BindCancelProposalFlags(cmd)

	err := cmd.Flag("bridge").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateCancelProposalFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *BridgeTestSuite) TestValidateQueryProposalFlags() {
	cmd := new(cobra.Command)
	BindQueryProposalFlags(cmd)

	err := cmd.Flag("bridge").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateQueryProposalFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *BridgeTestSuite) TestValidateQueryProposalInvalidAddress() {
	cmd := new(cobra.Command)
	BindQueryProposalFlags(cmd)

	err := cmd.Flag("bridge").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateQueryProposalFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *BridgeTestSuite) TestValidateQueryResourceFlags() {
	cmd := new(cobra.Command)
	BindQueryResourceFlags(cmd)

	err := cmd.Flag("handler").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateQueryResourceFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *BridgeTestSuite) TestValidateQueryResourceInvalidAddress() {
	cmd := new(cobra.Command)
	BindQueryResourceFlags(cmd)

	err := cmd.Flag("handler").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateQueryResourceFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *BridgeTestSuite) TestValidateRegisterGenericResourceFlags() {
	cmd := new(cobra.Command)
	BindRegisterGenericResourceFlags(cmd)

	err := cmd.Flag("handler").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("target").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("bridge").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateRegisterGenericResourceFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *BridgeTestSuite) TestValidateRegisterGenericResourceInvalidAddresses() {
	cmd := new(cobra.Command)
	BindRegisterGenericResourceFlags(cmd)

	err := cmd.Flag("handler").Value.Set(invalidAddr)
	s.Nil(err)
	err = cmd.Flag("target").Value.Set(invalidAddr)
	s.Nil(err)
	err = cmd.Flag("bridge").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateRegisterGenericResourceFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *BridgeTestSuite) TestValidateRegisterResourceFlags() {
	cmd := new(cobra.Command)
	BindRegisterResourceFlags(cmd)

	err := cmd.Flag("handler").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("target").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("bridge").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateRegisterResourceFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *BridgeTestSuite) TestValidateRegisterResourceInvalidAddresses() {
	cmd := new(cobra.Command)
	BindRegisterResourceFlags(cmd)

	err := cmd.Flag("handler").Value.Set(invalidAddr)
	s.Nil(err)
	err = cmd.Flag("target").Value.Set(invalidAddr)
	s.Nil(err)
	err = cmd.Flag("bridge").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateRegisterResourceFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *BridgeTestSuite) TestValidateSetBurnFlags() {
	cmd := new(cobra.Command)
	BindSetBurnFlags(cmd)

	err := cmd.Flag("handler").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("bridge").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("token-contract").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateSetBurnFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *BridgeTestSuite) TestValidateSetBurnInvalidAddresses() {
	cmd := new(cobra.Command)
	BindSetBurnFlags(cmd)

	err := cmd.Flag("handler").Value.Set(invalidAddr)
	s.Nil(err)
	err = cmd.Flag("bridge").Value.Set(invalidAddr)
	s.Nil(err)
	err = cmd.Flag("token-contract").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateSetBurnFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}
