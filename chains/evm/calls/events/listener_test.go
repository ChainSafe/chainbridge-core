package events_test

import (
	"encoding/hex"
	"strings"
	"testing"

	"github.com/ChainSafe/sygma-core/chains/evm/calls/consts"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/events"
	"github.com/ChainSafe/sygma-core/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

var (
	logData = "00000000000000000000000000000000000000000000000000000000000000020000000000000000000000d606a00c1a39da53ea7bb3ab570bbe40b156eb6600000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000120000000000000000000000000000000000000000000000000000000000000005400000000000000000000000000000000000000000000000000000000000f424000000000000000000000000000000000000000000000000000000000000000148e0a907331554af72563bd8d43051c2e64be5d350000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"
)

type EvmClientTestSuite struct {
	suite.Suite
	gomockController *gomock.Controller
	listener         events.Listener
}

func TestRunEvmClientTestSuite(t *testing.T) {
	suite.Run(t, new(EvmClientTestSuite))
}

func (s *EvmClientTestSuite) SetupSuite() {
	s.gomockController = gomock.NewController(s.T())
	s.listener = *events.NewListener(nil)
}

func (s *EvmClientTestSuite) TestUnpackDepositEventLogFailedUnpack() {
	abi, _ := abi.JSON(strings.NewReader(consts.BridgeABI))
	_, err := s.listener.UnpackDeposit(abi, []byte("invalid"))
	s.NotNil(err)
}

func (s *EvmClientTestSuite) TestUnpackDepositEventLogValidData() {
	abi, _ := abi.JSON(strings.NewReader(consts.BridgeABI))
	logDataBytes, _ := hex.DecodeString(logData)
	expectedRID := types.ResourceID(types.ResourceID{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0xd6, 0x6, 0xa0, 0xc, 0x1a, 0x39, 0xda, 0x53, 0xea, 0x7b, 0xb3, 0xab, 0x57, 0xb, 0xbe, 0x40, 0xb1, 0x56, 0xeb, 0x66, 0x0})
	dl, err := s.listener.UnpackDeposit(abi, logDataBytes)
	s.Nil(err)
	s.Equal(dl.SenderAddress.String(), "0x0000000000000000000000000000000000000000")
	s.Equal(dl.DepositNonce, uint64(1))
	s.Equal(dl.DestinationDomainID, uint8(2))
	s.Equal(dl.ResourceID, expectedRID)
	s.Equal(dl.HandlerResponse, []byte{})
}
