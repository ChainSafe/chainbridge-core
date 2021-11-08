package calls_test

import (
	"math/big"
	"testing"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type GetSolidityFunctionSigTestSuite struct {
	suite.Suite
	gomockController *gomock.Controller
}

func TestRunGetSolidityFunctionSigTestSuite(t *testing.T) {
	suite.Run(t, new(GetSolidityFunctionSigTestSuite))
}

func (s *GetSolidityFunctionSigTestSuite) SetupSuite()    {}
func (s *GetSolidityFunctionSigTestSuite) TearDownSuite() {}
func (s *GetSolidityFunctionSigTestSuite) SetupTest() {
	s.gomockController = gomock.NewController(s.T())
}
func (s *GetSolidityFunctionSigTestSuite) TearDownTest() {}

func (s *GetSolidityFunctionSigTestSuite) TestReturnsValidSolidityFunctionSig() {
	sig := calls.GetSolidityFunctionSig([]byte("store(bytes32)"))

	s.Equal(sig, [4]byte{0x65, 0x4c, 0xf8, 0x8c})
}

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
	got := calls.ToCallArg(msg)
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
	got := calls.ToCallArg(msg)
	want := map[string]interface{}{
		"from": common.HexToAddress(""),
		"to":   (*common.Address)(nil),
	}
	s.Equal(want, got)
}
