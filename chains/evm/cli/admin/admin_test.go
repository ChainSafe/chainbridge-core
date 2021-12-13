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

	cmd.Flag("admin").Value.Set(validAddr)
	cmd.Flag("bridge").Value.Set(validAddr)

	err := ValidateAddAdminFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateAddAdminFlagsInvalidAddresses() {
	cmd := addAdminCmd

	// invalid addresses
	cmd.Flag("admin").Value.Set(invalidAddr)
	cmd.Flag("bridge").Value.Set(invalidAddr)

	err := ValidateAddAdminFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateAddRelayerFlags() {
	cmd := addRelayerCmd

	cmd.Flag("relayer").Value.Set(validAddr)
	cmd.Flag("bridge").Value.Set(validAddr)

	err := ValidateAddRelayerFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateAddRelayerFlagsInvalidAddresses() {
	cmd := addRelayerCmd

	// invalid addresses
	cmd.Flag("relayer").Value.Set(invalidAddr)
	cmd.Flag("bridge").Value.Set(invalidAddr)

	err := ValidateAddRelayerFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateGetThresholdFlags() {
	cmd := getThresholdCmd

	cmd.Flag("bridge").Value.Set(validAddr)

	err := ValidateGetThresholdFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateGetThresholdFlagsInvalidAddress() {
	cmd := getThresholdCmd

	cmd.Flag("bridge").Value.Set(invalidAddr)

	err := ValidateGetThresholdFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateIsRelayerFlags() {
	cmd := isRelayerCmd

	cmd.Flag("bridge").Value.Set(validAddr)
	cmd.Flag("relayer").Value.Set(validAddr)

	err := ValidateIsRelayerFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateIsRelayerInvalidAddresses() {
	cmd := isRelayerCmd

	cmd.Flag("bridge").Value.Set(invalidAddr)
	cmd.Flag("relayer").Value.Set(invalidAddr)

	err := ValidateIsRelayerFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidatePauseFlags() {
	cmd := pauseCmd

	cmd.Flag("bridge").Value.Set(validAddr)

	err := ValidatePauseFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidatePauseInvalidAddress() {
	cmd := pauseCmd

	cmd.Flag("bridge").Value.Set(invalidAddr)

	err := ValidatePauseFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateRemoveAdminFlags() {
	cmd := removeAdminCmd

	cmd.Flag("bridge").Value.Set(validAddr)
	cmd.Flag("admin").Value.Set(validAddr)

	err := ValidateRemoveAdminFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateRemoveAdminInvalidAddresses() {
	cmd := removeAdminCmd

	cmd.Flag("bridge").Value.Set(invalidAddr)
	cmd.Flag("admin").Value.Set(invalidAddr)

	err := ValidateRemoveAdminFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateRemoveRelayerFlags() {
	cmd := removeRelayerCmd

	cmd.Flag("bridge").Value.Set(validAddr)
	cmd.Flag("relayer").Value.Set(validAddr)

	err := ValidateRemoveRelayerFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateRemoveRelayerInvalidAddresses() {
	cmd := removeRelayerCmd

	cmd.Flag("bridge").Value.Set(invalidAddr)
	cmd.Flag("relayer").Value.Set(invalidAddr)

	err := ValidateRemoveRelayerFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateSetDepositNonceFlags() {
	cmd := setDepositNonceCmd

	cmd.Flag("bridge").Value.Set(validAddr)

	err := ValidateSetDepositNonceFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateSetDepositNonceInvalidAddress() {
	cmd := setDepositNonceCmd

	cmd.Flag("bridge").Value.Set(invalidAddr)

	err := ValidateSetDepositNonceFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateSetFeeFlags() {
	cmd := setFeeCmd

	cmd.Flag("bridge").Value.Set(validAddr)

	err := ValidateSetFeeFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateSetFeeInvalidAddress() {
	cmd := setFeeCmd

	cmd.Flag("bridge").Value.Set(invalidAddr)

	err := ValidateSetFeeFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateSetThresholdFlags() {
	cmd := setThresholdCmd

	cmd.Flag("bridge").Value.Set(validAddr)

	err := ValidateSetThresholdFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateSetThresholdInvalidAddress() {
	cmd := setThresholdCmd

	cmd.Flag("bridge").Value.Set(invalidAddr)

	err := ValidateSetThresholdFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateUnpauseFlags() {
	cmd := unpauseCmd

	cmd.Flag("bridge").Value.Set(validAddr)

	err := ValidateUnpauseFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateUnpauseInvalidAddress() {
	cmd := unpauseCmd

	cmd.Flag("bridge").Value.Set(invalidAddr)

	err := ValidateUnpauseFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateWithdrawFlags() {
	cmd := withdrawCmd

	cmd.Flag("bridge").Value.Set(validAddr)
	cmd.Flag("handler").Value.Set(validAddr)
	cmd.Flag("token-contract").Value.Set(validAddr)
	cmd.Flag("recipient").Value.Set(validAddr)
	cmd.Flag("amount").Value.Set("1")

	err := ValidateWithdrawFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *AdminTestSuite) TestValidateWithdrawInvalidAddresses() {
	cmd := withdrawCmd

	cmd.Flag("bridge").Value.Set(invalidAddr)
	cmd.Flag("handler").Value.Set(invalidAddr)
	cmd.Flag("token-contract").Value.Set(invalidAddr)
	cmd.Flag("recipient").Value.Set(invalidAddr)
	cmd.Flag("amount").Value.Set("1")

	err := ValidateWithdrawFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}

func (s *AdminTestSuite) TestValidateWithdrawAmountTokenConflict() {
	cmd := withdrawCmd

	cmd.Flag("bridge").Value.Set(validAddr)
	cmd.Flag("handler").Value.Set(validAddr)
	cmd.Flag("token-contract").Value.Set(validAddr)
	cmd.Flag("recipient").Value.Set(validAddr)
	cmd.Flag("amount").Value.Set("1")
	cmd.Flag("token").Value.Set("1")

	err := ValidateWithdrawFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}
