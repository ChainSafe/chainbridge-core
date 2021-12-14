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
	cmd := &cobra.Command{}
	BindAddMinterCmdFlags(cmd)

	cmd.Flag("contract").Value.Set(validAddr)
	cmd.Flag("minter").Value.Set(validAddr)

	err := ValidateAddMinterFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *ERC721TestSuite) TestValidateAddMinterInvalidAddress() {
	cmd := addMinterCmd

	cmd.Flag("contract").Value.Set(invalidAddr)
	cmd.Flag("minter").Value.Set(invalidAddr)

	err := ValidateAddMinterFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *ERC721TestSuite) TestValidateApproveFlags() {
	cmd := approveCmd

	cmd.Flag("contract").Value.Set(validAddr)
	cmd.Flag("recipient").Value.Set(validAddr)

	err := ValidateApproveFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *ERC721TestSuite) TestValidateApproveInvalidAddress() {
	cmd := approveCmd

	cmd.Flag("contract").Value.Set(invalidAddr)
	cmd.Flag("recipient").Value.Set(invalidAddr)

	err := ValidateApproveFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *ERC721TestSuite) TestValidateDepositFlags() {
	cmd := depositCmd

	cmd.Flag("recipient").Value.Set(validAddr)
	cmd.Flag("bridge").Value.Set(validAddr)

	err := ValidateDepositFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *ERC721TestSuite) TestValidateDepositInvalidAddress() {
	cmd := depositCmd

	cmd.Flag("recipient").Value.Set(invalidAddr)
	cmd.Flag("bridge").Value.Set(invalidAddr)

	err := ValidateDepositFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *ERC721TestSuite) TestValidateMintFlags() {
	cmd := mintCmd

	cmd.Flag("contract").Value.Set(validAddr)
	cmd.Flag("recipient").Value.Set(validAddr)

	err := ValidateMintFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *ERC721TestSuite) TestValidateMintInvalidAddress() {
	cmd := mintCmd

	cmd.Flag("contract").Value.Set(invalidAddr)
	cmd.Flag("recipient").Value.Set(invalidAddr)

	err := ValidateMintFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *ERC721TestSuite) TestValidateOwnerFlags() {
	cmd := ownerCmd

	cmd.Flag("contract").Value.Set(validAddr)

	err := ValidateOwnerFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *ERC721TestSuite) TestValidateOwnerInvalidAddress() {
	cmd := ownerCmd

	cmd.Flag("contract").Value.Set(invalidAddr)

	err := ValidateOwnerFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}
