package account

import (
	"testing"

	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

var (
	validPrivateKey   = "6ec1ced059cb4a761dcee242dd17471398e863cb6f3a36cf5e570c648368803d"
	invalidPrivateKey = "0x6ec1ced059cb4a761dcee242dd17471398e863cb6f3a36cf5e570c648368803d"
)

type ImportPrivKeyTestSuite struct {
	suite.Suite
}

func TestValidateImportPrivKeyFlags(t *testing.T) {
	suite.Run(t, new(ImportPrivKeyTestSuite))
}

func (s *ImportPrivKeyTestSuite) SetupSuite() {
}
func (s *ImportPrivKeyTestSuite) TearDownSuite() {}

func (s *ImportPrivKeyTestSuite) TearDownTest() {}

func (s *ImportPrivKeyTestSuite) TestValidateImportPrivKeyFlags() {
	cmd := new(cobra.Command)
	BindImportPrivKeyFlags(cmd)

	err := cmd.Flag("private-key").Value.Set(validPrivateKey)
	s.Nil(err)

	err = ValidateImportPrivKeyFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *ImportPrivKeyTestSuite) TestValidateImportPrivKeyFlagsInvalidPrivKey() {
	cmd := new(cobra.Command)
	BindImportPrivKeyFlags(cmd)

	err := cmd.Flag("private-key").Value.Set(invalidPrivateKey)
	s.Nil(err)

	err = ValidateImportPrivKeyFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}
