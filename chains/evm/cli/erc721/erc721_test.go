package erc721

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

var (
	validAddr   = "0xd606A00c1A39dA53EA7Bb3Ab570BBE40b156EB66"
	invalidAddr = "0xd606A00c1A39dA53EA7Bb3Ab570BBE40b156EXYZ"
)

type ERC721TestSuite struct {
	suite.Suite
}

func TestERC721TestSuite(t *testing.T) {
	suite.Run(t, new(ERC721TestSuite))
}

func (s *ERC721TestSuite) SetupSuite() {
}
func (s *ERC721TestSuite) TearDownSuite() {}

func (s *ERC721TestSuite) TearDownTest() {}

func (s *ERC721TestSuite) TestValidateAddMinterFlags() {
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

func (s *ERC721TestSuite) TestValidateAddMinterInvalidAddress() {
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

func (s *ERC721TestSuite) TestValidateApproveFlags() {
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

func (s *ERC721TestSuite) TestValidateApproveInvalidAddress() {
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

func (s *ERC721TestSuite) TestValidateDepositFlags() {
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

func (s *ERC721TestSuite) TestValidateDepositInvalidAddress() {
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

func (s *ERC721TestSuite) TestValidateMintFlags() {
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

func (s *ERC721TestSuite) TestValidateMintInvalidAddress() {
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

func (s *ERC721TestSuite) TestValidateOwnerFlags() {
	cmd := new(cobra.Command)
	BindOwnerFlags(cmd)

	err := cmd.Flag("contract").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateOwnerFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *ERC721TestSuite) TestValidateOwnerInvalidAddress() {
	cmd := new(cobra.Command)
	BindOwnerFlags(cmd)

	err := cmd.Flag("contract").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateOwnerFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}
