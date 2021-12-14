package admin

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

var (
	validAddr   = "0xd606A00c1A39dA53EA7Bb3Ab570BBE40b156EB66"
	invalidAddr = "0xd606A00c1A39dA53EA7Bb3Ab570BBE40b156EXYZ"
)

type AdminTestSuite struct {
	suite.Suite
}

func TestAdminTestSuite(t *testing.T) {
	suite.Run(t, new(AdminTestSuite))
}

func (s *AdminTestSuite) SetupSuite() {
}
func (s *AdminTestSuite) TearDownSuite() {}

func (s *AdminTestSuite) TearDownTest() {}

func (s *AdminTestSuite) TestValidateAddAdminFlags() {
	cmd := new(cobra.Command)
	BindAddAdminFlags(cmd)

	err := cmd.Flag("admin").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("bridge").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateAddAdminFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateAddAdminFlagsInvalidAddresses() {
	cmd := new(cobra.Command)
	BindAddAdminFlags(cmd)

	err := cmd.Flag("admin").Value.Set(invalidAddr)
	s.Nil(err)
	err = cmd.Flag("bridge").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateAddAdminFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateAddRelayerFlags() {
	cmd := new(cobra.Command)
	BindAddRelayerFlags(cmd)

	err := cmd.Flag("relayer").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("bridge").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateAddRelayerFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateAddRelayerFlagsInvalidAddresses() {
	cmd := new(cobra.Command)
	BindAddRelayerFlags(cmd)

	err := cmd.Flag("relayer").Value.Set(invalidAddr)
	s.Nil(err)
	err = cmd.Flag("bridge").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateAddRelayerFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateGetThresholdFlags() {
	cmd := new(cobra.Command)
	BindGetThresholdFlags(cmd)

	err := cmd.Flag("bridge").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateGetThresholdFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateGetThresholdFlagsInvalidAddress() {
	cmd := new(cobra.Command)
	BindGetThresholdFlags(cmd)

	err := cmd.Flag("bridge").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateGetThresholdFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateIsRelayerFlags() {
	cmd := new(cobra.Command)
	BindIsRelayerFlags(cmd)

	err := cmd.Flag("bridge").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("relayer").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateIsRelayerFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateIsRelayerInvalidAddresses() {
	cmd := new(cobra.Command)
	BindIsRelayerFlags(cmd)

	err := cmd.Flag("bridge").Value.Set(invalidAddr)
	s.Nil(err)
	err = cmd.Flag("relayer").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateIsRelayerFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidatePauseFlags() {
	cmd := new(cobra.Command)
	BindPauseFlags(cmd)

	err := cmd.Flag("bridge").Value.Set(validAddr)
	s.Nil(err)

	err = ValidatePauseFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidatePauseInvalidAddress() {
	cmd := new(cobra.Command)
	BindPauseFlags(cmd)

	err := cmd.Flag("bridge").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidatePauseFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateRemoveAdminFlags() {
	cmd := new(cobra.Command)
	BindRemoveAdminFlags(cmd)

	err := cmd.Flag("bridge").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("admin").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateRemoveAdminFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateRemoveAdminInvalidAddresses() {
	cmd := new(cobra.Command)
	BindRemoveAdminFlags(cmd)

	err := cmd.Flag("bridge").Value.Set(invalidAddr)
	s.Nil(err)
	err = cmd.Flag("admin").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateRemoveAdminFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateRemoveRelayerFlags() {
	cmd := new(cobra.Command)
	BindRemoveRelayerFlags(cmd)

	err := cmd.Flag("bridge").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("relayer").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateRemoveRelayerFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateRemoveRelayerInvalidAddresses() {
	cmd := new(cobra.Command)
	BindRemoveRelayerFlags(cmd)

	err := cmd.Flag("bridge").Value.Set(invalidAddr)
	s.Nil(err)
	err = cmd.Flag("relayer").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateRemoveRelayerFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateSetDepositNonceFlags() {
	cmd := new(cobra.Command)
	BindSetDepositNonceFlags(cmd)

	err := cmd.Flag("bridge").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateSetDepositNonceFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateSetDepositNonceInvalidAddress() {
	cmd := new(cobra.Command)
	BindSetDepositNonceFlags(cmd)

	err := cmd.Flag("bridge").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateSetDepositNonceFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateSetFeeFlags() {
	cmd := new(cobra.Command)
	BindSetFeeFlags(cmd)

	err := cmd.Flag("bridge").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateSetFeeFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateSetFeeInvalidAddress() {
	cmd := new(cobra.Command)
	BindSetFeeFlags(cmd)

	err := cmd.Flag("bridge").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateSetFeeFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateSetThresholdFlags() {
	cmd := new(cobra.Command)
	BindSetThresholdFlags(cmd)

	err := cmd.Flag("bridge").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateSetThresholdFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateSetThresholdInvalidAddress() {
	cmd := new(cobra.Command)
	BindSetThresholdFlags(cmd)

	err := cmd.Flag("bridge").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateSetThresholdFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateUnpauseFlags() {
	cmd := new(cobra.Command)
	BindUnpauseFlags(cmd)

	err := cmd.Flag("bridge").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateUnpauseFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateUnpauseInvalidAddress() {
	cmd := new(cobra.Command)
	BindUnpauseFlags(cmd)

	err := cmd.Flag("bridge").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateUnpauseFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateWithdrawFlags() {
	cmd := new(cobra.Command)
	BindWithdrawFlags(cmd)

	err := cmd.Flag("bridge").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("handler").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("token-contract").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("recipient").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("amount").Value.Set("1")
	s.Nil(err)

	err = ValidateWithdrawFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateWithdrawInvalidAddresses() {
	cmd := new(cobra.Command)
	BindWithdrawFlags(cmd)

	err := cmd.Flag("bridge").Value.Set(invalidAddr)
	s.Nil(err)
	err = cmd.Flag("handler").Value.Set(invalidAddr)
	s.Nil(err)
	err = cmd.Flag("token-contract").Value.Set(invalidAddr)
	s.Nil(err)
	err = cmd.Flag("recipient").Value.Set(invalidAddr)
	s.Nil(err)
	err = cmd.Flag("amount").Value.Set("1")
	s.Nil(err)

	err = ValidateWithdrawFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateWithdrawAmountTokenConflict() {
	cmd := new(cobra.Command)
	BindWithdrawFlags(cmd)

	err := cmd.Flag("bridge").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("handler").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("token-contract").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("recipient").Value.Set(validAddr)
	s.Nil(err)
	err = cmd.Flag("amount").Value.Set("1")
	s.Nil(err)
	err = cmd.Flag("token").Value.Set("1")
	s.Nil(err)

	err = ValidateWithdrawFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}
