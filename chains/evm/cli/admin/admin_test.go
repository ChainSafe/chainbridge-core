package admin

import (
	"testing"

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
	cmd := addAdminCmd

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
	cmd := addAdminCmd

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
	cmd := addRelayerCmd

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
	cmd := addRelayerCmd

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
	cmd := getThresholdCmd

	err := cmd.Flag("bridge").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateGetThresholdFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateGetThresholdFlagsInvalidAddress() {
	cmd := getThresholdCmd

	err := cmd.Flag("bridge").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateGetThresholdFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateIsRelayerFlags() {
	cmd := isRelayerCmd

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
	cmd := isRelayerCmd

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
	cmd := pauseCmd

	err := cmd.Flag("bridge").Value.Set(validAddr)
	s.Nil(err)

	err = ValidatePauseFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidatePauseInvalidAddress() {
	cmd := pauseCmd

	err := cmd.Flag("bridge").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidatePauseFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateRemoveAdminFlags() {
	cmd := removeAdminCmd

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
	cmd := removeAdminCmd

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
	cmd := removeRelayerCmd

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
	cmd := removeRelayerCmd

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
	cmd := setDepositNonceCmd

	err := cmd.Flag("bridge").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateSetDepositNonceFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateSetDepositNonceInvalidAddress() {
	cmd := setDepositNonceCmd

	err := cmd.Flag("bridge").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateSetDepositNonceFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateSetFeeFlags() {
	cmd := setFeeCmd

	err := cmd.Flag("bridge").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateSetFeeFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateSetFeeInvalidAddress() {
	cmd := setFeeCmd

	err := cmd.Flag("bridge").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateSetFeeFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateSetThresholdFlags() {
	cmd := setThresholdCmd

	err := cmd.Flag("bridge").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateSetThresholdFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateSetThresholdInvalidAddress() {
	cmd := setThresholdCmd

	err := cmd.Flag("bridge").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateSetThresholdFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateUnpauseFlags() {
	cmd := unpauseCmd

	err := cmd.Flag("bridge").Value.Set(validAddr)
	s.Nil(err)

	err = ValidateUnpauseFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateUnpauseInvalidAddress() {
	cmd := unpauseCmd

	err := cmd.Flag("bridge").Value.Set(invalidAddr)
	s.Nil(err)

	err = ValidateUnpauseFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateWithdrawFlags() {
	cmd := withdrawCmd

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
	cmd := withdrawCmd

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
	cmd := withdrawCmd

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
