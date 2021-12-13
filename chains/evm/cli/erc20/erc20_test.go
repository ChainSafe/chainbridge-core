package erc20

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

var (
	validAddr   = "0xd606A00c1A39dA53EA7Bb3Ab570BBE40b156EB66"
	invalidAddr = "0xd606A00c1A39dA53EA7Bb3Ab570BBE40b156EXYZ"
)

type ERC20TestSuite struct {
	suite.Suite
}

func TestERC20TestSuite(t *testing.T) {
	suite.Run(t, new(ERC20TestSuite))
}

func (s *ERC20TestSuite) SetupSuite() {
}
func (s *ERC20TestSuite) TearDownSuite() {}

func (s *ERC20TestSuite) TearDownTest() {}

func (s *ERC20TestSuite) TestValidateAddMinterFlags() {
	cmd := addMinterCmd

	cmd.Flag("contract").Value.Set(validAddr)
	cmd.Flag("minter").Value.Set(validAddr)

	err := ValidateAddMinterFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *ERC20TestSuite) TestValidateAddMinterInvalidAddress() {
	cmd := addMinterCmd

	cmd.Flag("contract").Value.Set(invalidAddr)
	cmd.Flag("minter").Value.Set(invalidAddr)

	err := ValidateAddMinterFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *ERC20TestSuite) TestValidateApproveFlags() {
	cmd := approveCmd

	cmd.Flag("contract").Value.Set(validAddr)
	cmd.Flag("recipient").Value.Set(validAddr)

	err := ValidateApproveFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *ERC20TestSuite) TestValidateApproveInvalidAddress() {
	cmd := approveCmd

	cmd.Flag("contract").Value.Set(invalidAddr)
	cmd.Flag("recipient").Value.Set(invalidAddr)

	err := ValidateApproveFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *ERC20TestSuite) TestValidateBalanceFlags() {
	cmd := balanceCmd

	cmd.Flag("contract").Value.Set(validAddr)
	cmd.Flag("address").Value.Set(validAddr)

	err := ValidateBalanceFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *ERC20TestSuite) TestValidateBalanceInvalidAddress() {
	cmd := balanceCmd

	cmd.Flag("contract").Value.Set(invalidAddr)
	cmd.Flag("address").Value.Set(invalidAddr)

	err := ValidateBalanceFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *ERC20TestSuite) TestValidateDepositFlags() {
	cmd := depositCmd

	cmd.Flag("recipient").Value.Set(validAddr)
	cmd.Flag("bridge").Value.Set(validAddr)

	err := ValidateDepositFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *ERC20TestSuite) TestValidateDepositInvalidAddress() {
	cmd := depositCmd

	cmd.Flag("recipient").Value.Set(invalidAddr)
	cmd.Flag("bridge").Value.Set(invalidAddr)

	err := ValidateDepositFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *ERC20TestSuite) TestValidateGetAllowanceFlags() {
	cmd := getAllowanceCmd

	cmd.Flag("contract").Value.Set(validAddr)
	cmd.Flag("owner").Value.Set(validAddr)
	cmd.Flag("spender").Value.Set(validAddr)

	err := ValidateGetAllowanceFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *ERC20TestSuite) TestValidateGetAllowanceInvalidAddress() {
	cmd := getAllowanceCmd

	cmd.Flag("contract").Value.Set(invalidAddr)
	cmd.Flag("owner").Value.Set(invalidAddr)
	cmd.Flag("spender").Value.Set(invalidAddr)

	err := ValidateGetAllowanceFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *ERC20TestSuite) TestValidateMintFlags() {
	cmd := mintCmd

	cmd.Flag("contract").Value.Set(validAddr)
	cmd.Flag("recipient").Value.Set(validAddr)

	err := ValidateMintFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *ERC20TestSuite) TestValidateMintInvalidAddress() {
	cmd := mintCmd

	cmd.Flag("contract").Value.Set(invalidAddr)
	cmd.Flag("recipient").Value.Set(invalidAddr)

	err := ValidateMintFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}
