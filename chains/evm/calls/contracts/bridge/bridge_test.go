package bridge_test

import (
	"encoding/hex"
	"errors"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/bridge"
	mock_calls "github.com/ChainSafe/chainbridge-core/chains/evm/calls/mock"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"math/big"
	"testing"

	mock_transactor "github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor/mock"
	"github.com/ChainSafe/chainbridge-core/chains/evm/voter/proposal"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type ProposalStatusTestSuite struct {
	suite.Suite
	mockContractCaller *mock_calls.MockContractCallerDispatcher
	mockTransactor     *mock_transactor.MockTransactor
	bridgeAddress      common.Address
}

func TestRunProposalStatusTestSuite(t *testing.T) {
	suite.Run(t, new(ProposalStatusTestSuite))
}

func (s *ProposalStatusTestSuite) SetupSuite()    {}
func (s *ProposalStatusTestSuite) TearDownSuite() {}
func (s *ProposalStatusTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	s.mockContractCaller = mock_calls.NewMockContractCallerDispatcher(gomockController)
	s.mockTransactor = mock_transactor.NewMockTransactor(gomockController)
	s.bridgeAddress = common.HexToAddress("0x3162226db165D8eA0f51720CA2bbf44Db2105ADF")
}
func (s *ProposalStatusTestSuite) TearDownTest() {}

func (s *ProposalStatusTestSuite) TestProposalStatusFailedContractCall() {
	s.mockContractCaller.EXPECT().CallContract(gomock.Any(), gomock.Any(), nil).Return(nil, errors.New("error"))
	s.mockContractCaller.EXPECT().From().Times(1).Return(common.Address{})
	bc := bridge.NewBridgeContract(s.mockContractCaller, common.Address{}, s.mockTransactor)
	status, err := bc.ProposalStatus(&proposal.Proposal{})
	s.Equal(message.ProposalStatus{}, status)
	s.NotNil(err)
}

func (s *ProposalStatusTestSuite) TestProposalStatusFailedUnpack() {
	s.mockContractCaller.EXPECT().CallContract(gomock.Any(), gomock.Any(), nil).Return([]byte("invalid"), nil)
	s.mockContractCaller.EXPECT().From().Times(1).Return(common.Address{})
	bc := bridge.NewBridgeContract(s.mockContractCaller, common.Address{}, s.mockTransactor)
	status, err := bc.ProposalStatus(&proposal.Proposal{})

	s.Equal(message.ProposalStatus{}, status)
	s.NotNil(err)
}

func (s *ProposalStatusTestSuite) TestProposalStatusSuccessfulCall() {
	proposalStatus, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000001c0000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000001f")
	s.mockContractCaller.EXPECT().CallContract(gomock.Any(), gomock.Any(), nil).Return(proposalStatus, nil)
	s.mockContractCaller.EXPECT().From().Times(1).Return(common.Address{})
	bc := bridge.NewBridgeContract(s.mockContractCaller, common.Address{}, s.mockTransactor)
	status, err := bc.ProposalStatus(&proposal.Proposal{})

	s.Nil(err)
	s.Equal(status.YesVotesTotal, uint8(3))
	s.Equal(status.Status, message.ProposalStatusExecuted)
}

func (s *ProposalStatusTestSuite) TestPrepareWithdrawInput() {
	handlerAddress := common.HexToAddress("0x3167776db165D8eA0f51790CA2bbf44Db5105ADF")
	tokenAddress := common.HexToAddress("0x3f709398808af36ADBA86ACC617FeB7F5B7B193E")
	recipientAddress := common.HexToAddress("0x8e5F72B158BEDf0ab50EDa78c70dFC118158C272")
	amountOrTokenId := big.NewInt(1)

	s.mockTransactor.EXPECT().Transact(&s.bridgeAddress, gomock.Any(), gomock.Any()).Times(1).Return(
		&common.Hash{}, nil,
	)

	bc := bridge.NewBridgeContract(s.mockContractCaller, s.bridgeAddress, s.mockTransactor)
	inputBytes, err := bc.Withdraw(
		handlerAddress, tokenAddress, recipientAddress, amountOrTokenId, transactor.TransactOptions{},
	)
	s.NotNil(inputBytes)
	s.Nil(err)
}

func (s *ProposalStatusTestSuite) TestDeployContract() {
	bc := bridge.NewBridgeContract(s.mockContractCaller, s.bridgeAddress, s.mockTransactor)
	a := common.HexToAddress("0x12847465a15b58D4351AfB47F0CBbeebE93B06e3")
	address, err := bc.PackMethod("",
		uint8(1),
		[]common.Address{a},
		big.NewInt(0).SetUint64(1),
		big.NewInt(0).SetUint64(1),
		big.NewInt(0),
	)
	s.Nil(err)
	s.NotNil(address)
}
