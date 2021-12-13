package account

import (
	"testing"

	"github.com/stretchr/testify/suite"
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
	cmd := importPrivKeyCmd

	cmd.Flag("private-key").Value.Set(
		"6ec1ced059cb4a761dcee242dd17471398e863cb6f3a36cf5e570c648368803d",
	)

	err := ValidateImportPrivKeyFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *ImportPrivKeyTestSuite) TestValidateImportPrivKeyFlagsInvalidPrivKey() {
	cmd := importPrivKeyCmd

	cmd.Flag("private-key").Value.Set("0x6ec1ced059cb4a761dcee242dd17471398e863cb6f3a36cf5e570c648368803d") // invalid private key

	err := ValidateImportPrivKeyFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}
