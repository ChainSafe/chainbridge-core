package bridge

import (
	"testing"

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
	cmd := cancelProposalCmd

	cmd.Flag("bridge").Value.Set(validAddr)

	err := ValidateCancelProposalFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *BridgeTestSuite) TestValidateCancelProposalInvalidAddress() {
	cmd := cancelProposalCmd

	// invalid addresses
	cmd.Flag("bridge").Value.Set(invalidAddr)

	err := ValidateCancelProposalFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *BridgeTestSuite) TestValidateQueryProposalFlags() {
	cmd := queryProposalCmd

	cmd.Flag("bridge").Value.Set(validAddr)

	err := ValidateQueryProposalFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *BridgeTestSuite) TestValidateQueryProposalInvalidAddress() {
	cmd := queryProposalCmd

	// invalid addresses
	cmd.Flag("bridge").Value.Set(invalidAddr)

	err := ValidateQueryProposalFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *BridgeTestSuite) TestValidateQueryResourceFlags() {
	cmd := queryResourceCmd

	cmd.Flag("handler").Value.Set(validAddr)

	err := ValidateQueryResourceFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *BridgeTestSuite) TestValidateQueryResourceInvalidAddress() {
	cmd := queryResourceCmd

	// invalid addresses
	cmd.Flag("handler").Value.Set(invalidAddr)

	err := ValidateQueryResourceFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *BridgeTestSuite) TestValidateRegisterGenericResourceFlags() {
	cmd := registerGenericResourceCmd

	cmd.Flag("handler").Value.Set(validAddr)
	cmd.Flag("target").Value.Set(validAddr)
	cmd.Flag("bridge").Value.Set(validAddr)

	err := ValidateRegisterGenericResourceFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *BridgeTestSuite) TestValidateRegisterGenericResourceInvalidAddresses() {
	cmd := registerGenericResourceCmd

	// invalid addresses
	cmd.Flag("handler").Value.Set(invalidAddr)
	cmd.Flag("target").Value.Set(invalidAddr)
	cmd.Flag("bridge").Value.Set(invalidAddr)

	err := ValidateRegisterGenericResourceFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *BridgeTestSuite) TestValidateRegisterResourceFlags() {
	cmd := registerResourceCmd

	cmd.Flag("handler").Value.Set(validAddr)
	cmd.Flag("target").Value.Set(validAddr)
	cmd.Flag("bridge").Value.Set(validAddr)

	err := ValidateRegisterResourceFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *BridgeTestSuite) TestValidateRegisterResourceInvalidAddresses() {
	cmd := registerResourceCmd

	// invalid addresses
	cmd.Flag("handler").Value.Set(invalidAddr)
	cmd.Flag("target").Value.Set(invalidAddr)
	cmd.Flag("bridge").Value.Set(invalidAddr)

	err := ValidateRegisterResourceFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *BridgeTestSuite) TestValidateSetBurnFlags() {
	cmd := setBurnCmd

	cmd.Flag("handler").Value.Set(validAddr)
	cmd.Flag("bridge").Value.Set(validAddr)
	cmd.Flag("token-contract").Value.Set(validAddr)

	err := ValidateSetBurnFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *BridgeTestSuite) TestValidateSetBurnInvalidAddresses() {
	cmd := setBurnCmd

	// invalid addresses
	cmd.Flag("handler").Value.Set(invalidAddr)
	cmd.Flag("bridge").Value.Set(invalidAddr)
	cmd.Flag("token-contract").Value.Set(invalidAddr)

	err := ValidateSetBurnFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}
