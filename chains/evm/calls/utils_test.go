package calls

import (
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/stretchr/testify/suite"

	"math/big"

	"testing"
)

type UtilsTestSuite struct {
	suite.Suite
}

func TestRunUtilsTestSuite(t *testing.T) {
	suite.Run(t, new(UtilsTestSuite))
}

func (s *UtilsTestSuite) SetupSuite()    {}
func (s *UtilsTestSuite) TearDownSuite() {}
func (s *UtilsTestSuite) SetupTest()     {}
func (s *UtilsTestSuite) TearDownTest()  {}

func (s *UtilsTestSuite) TestToCallArg() {
	kp, err := secp256k1.GenerateKeypair()

	s.Nil(err)
	address := common.HexToAddress(kp.Address())

	msg := ethereum.CallMsg{
		From:     common.Address{},
		To:       &address,
		Value:    big.NewInt(1),
		Gas:      uint64(21000),
		GasPrice: big.NewInt(3000),
		Data:     []byte("test"),
	}
	got := ToCallArg(msg)
	want := map[string]interface{}{
		"from":     msg.From,
		"to":       msg.To,
		"value":    (*hexutil.Big)(msg.Value),
		"gas":      hexutil.Uint64(msg.Gas),
		"gasPrice": (*hexutil.Big)(msg.GasPrice),
		"data":     hexutil.Bytes(msg.Data),
	}
	s.Equal(want, got)
}

func (s *UtilsTestSuite) TestToCallArgWithEmptyMessage() {
	msg := ethereum.CallMsg{}
	got := ToCallArg(msg)
	want := map[string]interface{}{
		"from": common.HexToAddress(""),
		"to":   (*common.Address)(nil),
	}
	s.Equal(want, got)
}
