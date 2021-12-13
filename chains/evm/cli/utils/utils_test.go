package utils

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

var (
	validAddr     = "0xd606A00c1A39dA53EA7Bb3Ab570BBE40b156EB66"
	invalidAddr   = "0xd606A00c1A39dA53EA7Bb3Ab570BBE40b156EXYZ"
	validTxHash   = "0x455096e686c929229577767350d5c9151c609c2ba3e50a447e7092018d7f2dac"
	invalidTxHash = "455096e686c929229577767350d5c9151c609c2ba3e50a447e7092018d7f2dac"
)

type UtilsTestSuite struct {
	suite.Suite
}

func TestUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}

func (s *UtilsTestSuite) SetupSuite() {
}
func (s *UtilsTestSuite) TearDownSuite() {}

func (s *UtilsTestSuite) TearDownTest() {}

func (s *UtilsTestSuite) TestValidateSimulateFlags() {
	cmd := simulateCmd

	cmd.Flag("from").Value.Set(validAddr)
	cmd.Flag("tx-hash").Value.Set(validTxHash)

	err := ValidateSimulateFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *UtilsTestSuite) TestValidateSimulateInvalidAddress() {
	cmd := simulateCmd

	cmd.Flag("from").Value.Set(invalidAddr)
	cmd.Flag("tx-hash").Value.Set(validTxHash)

	err := ValidateSimulateFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *UtilsTestSuite) TestValidateSimulateInvalidTxHash() {
	cmd := simulateCmd

	cmd.Flag("from").Value.Set(validAddr)
	cmd.Flag("tx-hash").Value.Set(invalidTxHash)

	err := ValidateSimulateFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}
