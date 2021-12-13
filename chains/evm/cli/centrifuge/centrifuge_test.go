package centrifuge

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

var (
	validAddr   = "0xd606A00c1A39dA53EA7Bb3Ab570BBE40b156EB66"
	invalidAddr = "0xd606A00c1A39dA53EA7Bb3Ab570BBE40b156EXYZ"
)

type CentrifugeTestSuite struct {
	suite.Suite
}

func TestCentrifugeTestSuite(t *testing.T) {
	suite.Run(t, new(CentrifugeTestSuite))
}

func (s *CentrifugeTestSuite) SetupSuite() {
}
func (s *CentrifugeTestSuite) TearDownSuite() {}

func (s *CentrifugeTestSuite) TearDownTest() {}

func (s *CentrifugeTestSuite) TestValidateGetHashFlags() {
	cmd := getHashCmd

	cmd.Flag("address").Value.Set(validAddr)

	err := ValidateGetHashFlags(
		cmd,
		[]string{},
	)
	s.Nil(err)
}

func (s *CentrifugeTestSuite) TestValidateGetHashInvalidAddress() {
	cmd := getHashCmd

	cmd.Flag("address").Value.Set(invalidAddr)

	err := ValidateGetHashFlags(
		cmd,
		[]string{},
	)
	s.NotNil(err)
}
