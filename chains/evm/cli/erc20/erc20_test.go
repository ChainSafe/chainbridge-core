package erc20

import (
	"testing"

	"github.com/spf13/cobra"
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
	cmd := new(cobra.Command)
	BindAddMinterFlags(cmd)

	err := cmd.Flag("contract").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("minter").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateAddMinterFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *ERC20TestSuite) TestValidateAddMinterInvalidAddress() {
	cmd := new(cobra.Command)
	BindAddMinterFlags(cmd)

	err := cmd.Flag("contract").Value.Set(invalidAddr)
	s.Nil(err)
	err = cmd.Flag("minter").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateAddMinterFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *ERC20TestSuite) TestValidateApproveFlags() {
	cmd := new(cobra.Command)
	BindApproveFlags(cmd)

	err := cmd.Flag("contract").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("recipient").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateApproveFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *ERC20TestSuite) TestValidateApproveInvalidAddress() {
	cmd := new(cobra.Command)
	BindApproveFlags(cmd)

	err := cmd.Flag("contract").Value.Set(invalidAddr)
	s.Nil(err)
	err = cmd.Flag("recipient").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateApproveFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *ERC20TestSuite) TestValidateBalanceFlags() {
	cmd := new(cobra.Command)
	BindBalanceFlags(cmd)

	err := cmd.Flag("contract").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("address").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateBalanceFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *ERC20TestSuite) TestValidateBalanceInvalidAddress() {
	cmd := new(cobra.Command)
	BindBalanceFlags(cmd)

	err := cmd.Flag("contract").Value.Set(invalidAddr)
	s.Nil(err)
	err = cmd.Flag("address").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateBalanceFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *ERC20TestSuite) TestValidateDepositFlags() {
	cmd := new(cobra.Command)
	BindDepositFlags(cmd)

	err := cmd.Flag("recipient").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("bridge").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateDepositFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *ERC20TestSuite) TestValidateDepositInvalidAddress() {
	cmd := new(cobra.Command)
	BindDepositFlags(cmd)

	err := cmd.Flag("recipient").Value.Set(invalidAddr)
	s.Nil(err)
	err = cmd.Flag("bridge").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateDepositFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *ERC20TestSuite) TestValidateGetAllowanceFlags() {
	cmd := new(cobra.Command)
	BindGetAllowanceFlags(cmd)

	err := cmd.Flag("contract").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("owner").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("spender").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateGetAllowanceFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *ERC20TestSuite) TestValidateGetAllowanceInvalidAddress() {
	cmd := new(cobra.Command)
	BindGetAllowanceFlags(cmd)

	err := cmd.Flag("contract").Value.Set(invalidAddr)
	s.Nil(err)
	err = cmd.Flag("owner").Value.Set(invalidAddr)
	s.Nil(err)
	err = cmd.Flag("spender").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateGetAllowanceFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *ERC20TestSuite) TestValidateMintFlags() {
	cmd := new(cobra.Command)
	BindMintFlags(cmd)

	err := cmd.Flag("contract").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("recipient").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateMintFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *ERC20TestSuite) TestValidateMintInvalidAddress() {
	cmd := new(cobra.Command)
	BindMintFlags(cmd)

	err := cmd.Flag("contract").Value.Set(invalidAddr)
	s.Nil(err)
	err = cmd.Flag("recipient").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateMintFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}
