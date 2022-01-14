package bridge_test

import (
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/bridge"
	mock_calls "github.com/ChainSafe/chainbridge-core/chains/evm/calls/mock"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"

	mock_transactor "github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor/mock"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor/signAndSend"
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
	bridgeContract     *bridge.BridgeContract
	proposal           proposal.Proposal
}

var (
	testContractAddress   = "0x5f75ce92326e304962b22749bd71e36976171285"
	testInteractorAddress = "0x1B33100D4f077f027042c01241D617b264d77931"
	testRelayerAddress    = "0x0E223343BE5E126d7Cd1F6228F8F86fA04aD80fe"
	testHandlerAddress    = "0xb157b07c616860546464b733a056be414167a09b"
	testResourceId        = [32]byte{66}
	testDomainId          = uint8(0)
)

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
	s.bridgeContract = bridge.NewBridgeContract(s.mockContractCaller, common.HexToAddress(testContractAddress), s.mockTransactor)
	s.proposal = *proposal.NewProposal(
		uint8(1),
		uint64(1),
		testResourceId,
		[]byte{},
		common.HexToAddress(testHandlerAddress),
		s.bridgeAddress,
	)
}
func (s *ProposalStatusTestSuite) TearDownTest() {}

func (s *ProposalStatusTestSuite) TestProposalContractCall_StatusFailed() {
	s.mockContractCaller.EXPECT().CallContract(gomock.Any(), gomock.Any(), nil).Return(nil, errors.New("error"))
	s.mockContractCaller.EXPECT().From().Times(1).Return(common.Address{})
	bc := bridge.NewBridgeContract(s.mockContractCaller, common.Address{}, s.mockTransactor)
	status, err := bc.ProposalStatus(&proposal.Proposal{})
	s.Equal(message.ProposalStatus{}, status)
	s.NotNil(err)
}

func (s *ProposalStatusTestSuite) TestProposalStatus_FailedUnpack() {
	s.mockContractCaller.EXPECT().CallContract(gomock.Any(), gomock.Any(), nil).Return([]byte("invalid"), nil)
	s.mockContractCaller.EXPECT().From().Times(1).Return(common.Address{})
	bc := bridge.NewBridgeContract(s.mockContractCaller, common.Address{}, s.mockTransactor)
	status, err := bc.ProposalStatus(&proposal.Proposal{})

	s.Equal(message.ProposalStatus{}, status)
	s.NotNil(err)
}

func (s *ProposalStatusTestSuite) TestProposalStatus_SuccessfulCall() {
	proposalStatus, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000001c0000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000001f")
	s.mockContractCaller.EXPECT().CallContract(gomock.Any(), gomock.Any(), nil).Return(proposalStatus, nil)
	s.mockContractCaller.EXPECT().From().Times(1).Return(common.Address{})
	bc := bridge.NewBridgeContract(s.mockContractCaller, common.Address{}, s.mockTransactor)
	status, err := bc.ProposalStatus(&proposal.Proposal{})

	s.Nil(err)
	s.Equal(status.YesVotesTotal, uint8(3))
	s.Equal(status.Status, message.ProposalStatusExecuted)
}

func (s *ProposalStatusTestSuite) TestPrepare_WithdrawInput_Success() {
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

func (s *ProposalStatusTestSuite) TestDeployContract_Success() {
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

func (s *ProposalStatusTestSuite) TestBridge_AddRelayer_Success() {
	s.mockTransactor.EXPECT().Transact(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(&common.Hash{1, 2, 3, 4, 5}, nil)
	res, err := s.bridgeContract.AddRelayer(common.HexToAddress(testRelayerAddress), signAndSend.DefaultTransactionOptions)
	s.Equal(
		&common.Hash{1, 2, 3, 4, 5},
		res,
	)
	s.Nil(err)
}

func (s *ProposalStatusTestSuite) TestBridge_AdminSetGenericResource_Success() {
	s.mockTransactor.EXPECT().Transact(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(&common.Hash{10, 11, 12, 13, 14}, nil)
	res, err := s.bridgeContract.AdminSetGenericResource(common.HexToAddress(testHandlerAddress), [32]byte{1}, common.HexToAddress(testInteractorAddress), [4]byte{2}, big.NewInt(0), [4]byte{3}, signAndSend.DefaultTransactionOptions)
	s.Equal(
		&common.Hash{10, 11, 12, 13, 14},
		res,
	)
	s.Nil(err)
}

func (s *ProposalStatusTestSuite) TestBridge_AdminSetResource_Success() {
	s.mockTransactor.EXPECT().Transact(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(&common.Hash{10, 11, 12, 13, 14, 15, 16}, nil)
	res, err := s.bridgeContract.AdminSetResource(common.HexToAddress(testHandlerAddress), testResourceId, common.HexToAddress(testContractAddress), signAndSend.DefaultTransactionOptions)
	s.Equal(
		&common.Hash{10, 11, 12, 13, 14, 15, 16},
		res,
	)
	s.Nil(err)
}

func (s *ProposalStatusTestSuite) TestBridge_SetDepositNonce_Success() {
	s.mockTransactor.EXPECT().Transact(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(&common.Hash{20, 21, 22, 23}, nil)
	res, err := s.bridgeContract.SetDepositNonce(testDomainId, uint64(0), signAndSend.DefaultTransactionOptions)
	s.Equal(
		&common.Hash{20, 21, 22, 23},
		res,
	)
	s.Nil(err)
}

func (s *ProposalStatusTestSuite) TestBridge_SetThresholdInput_Success() {
	s.mockTransactor.EXPECT().Transact(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(&common.Hash{22, 23, 24, 25}, nil)
	res, err := s.bridgeContract.SetThresholdInput(uint64(2), signAndSend.DefaultTransactionOptions)
	s.Equal(
		&common.Hash{22, 23, 24, 25},
		res,
	)
	s.Nil(err)
}

func (s *ProposalStatusTestSuite) TestBridge_SetBurnableInput_Success() {
	s.mockTransactor.EXPECT().Transact(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(&common.Hash{25, 26, 27, 28}, nil)
	res, err := s.bridgeContract.SetBurnableInput(common.HexToAddress(testHandlerAddress), common.HexToAddress(testContractAddress), signAndSend.DefaultTransactionOptions)
	s.Equal(
		&common.Hash{25, 26, 27, 28},
		res,
	)
	s.Nil(err)
}

func (s *ProposalStatusTestSuite) TestBridge_Erc20Deposit_Success() {
	s.mockTransactor.EXPECT().Transact(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(&common.Hash{30, 31, 32, 33}, nil)
	res, err := s.bridgeContract.Erc20Deposit(common.HexToAddress(testInteractorAddress), big.NewInt(10), testResourceId, testDomainId, signAndSend.DefaultTransactionOptions)
	s.Equal(
		&common.Hash{30, 31, 32, 33},
		res,
	)
	s.Nil(err)
}

func (s *ProposalStatusTestSuite) TestBridge_Erc721Deposit_Success() {
	s.mockTransactor.EXPECT().Transact(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(&common.Hash{32, 33, 34, 35}, nil)
	res, err := s.bridgeContract.Erc721Deposit(big.NewInt(55), "token_uri", common.HexToAddress(testInteractorAddress), testResourceId, testDomainId, signAndSend.DefaultTransactionOptions)
	s.Equal(
		&common.Hash{32, 33, 34, 35},
		res,
	)
	s.Nil(err)
}

func (s *ProposalStatusTestSuite) TestBridge_GenericDeposit_Success() {
	s.mockTransactor.EXPECT().Transact(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(&common.Hash{35, 36, 37, 38}, nil)
	res, err := s.bridgeContract.GenericDeposit([]byte{1, 2, 3}, testResourceId, testDomainId, signAndSend.DefaultTransactionOptions)
	s.Equal(
		&common.Hash{35, 36, 37, 38},
		res,
	)
	s.Nil(err)
}

func (s *ProposalStatusTestSuite) TestBridge_ExecuteProposal_Success() {
	s.mockTransactor.EXPECT().Transact(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(&common.Hash{36, 37, 38}, nil)
	res, err := s.bridgeContract.ExecuteProposal(&s.proposal, signAndSend.DefaultTransactionOptions)
	s.Equal(
		&common.Hash{36, 37, 38},
		res,
	)
	s.Nil(err)
}

func (s *ProposalStatusTestSuite) TestBridge_VoteProposal_Success() {
	s.mockTransactor.EXPECT().Transact(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(&common.Hash{37, 38, 39}, nil)
	res, err := s.bridgeContract.VoteProposal(&s.proposal, signAndSend.DefaultTransactionOptions)
	s.Equal(
		&common.Hash{37, 38, 39},
		res,
	)
	s.Nil(err)
}

func (s *ProposalStatusTestSuite) TestBridge_SimulateVoteProposal_Success() {
	s.mockContractCaller.EXPECT().From().Return(common.HexToAddress(testInteractorAddress))
	s.mockContractCaller.EXPECT().CallContract(
		gomock.Any(),
		gomock.Any(),
		nil,
	).Return([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10}, nil)
	err := s.bridgeContract.SimulateVoteProposal(&s.proposal)
	s.Nil(err)
}

func (s *ProposalStatusTestSuite) TestBridge_Pause_Success() {
	s.mockTransactor.EXPECT().Transact(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(&common.Hash{40, 41, 42}, nil)
	res, err := s.bridgeContract.Pause(signAndSend.DefaultTransactionOptions)
	s.Equal(
		&common.Hash{40, 41, 42},
		res,
	)
	s.Nil(err)
}

func (s *ProposalStatusTestSuite) TestBridge_Unpause_Success() {
	s.mockTransactor.EXPECT().Transact(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(&common.Hash{42, 43, 44}, nil)
	res, err := s.bridgeContract.Unpause(signAndSend.DefaultTransactionOptions)
	s.Equal(
		&common.Hash{42, 43, 44},
		res,
	)
	s.Nil(err)
}

func (s *ProposalStatusTestSuite) TestBridge_Withdraw_Success() {
	s.mockTransactor.EXPECT().Transact(
		gomock.Any(),
		gomock.Any(),
		gomock.Any(),
	).Return(&common.Hash{44, 45, 46}, nil)
	res, err := s.bridgeContract.Withdraw(
		common.HexToAddress(testHandlerAddress),
		common.HexToAddress(testContractAddress),
		common.HexToAddress(testInteractorAddress),
		big.NewInt(5), signAndSend.DefaultTransactionOptions,
	)
	s.Equal(
		&common.Hash{44, 45, 46},
		res,
	)
	s.Nil(err)
}

func (s *ProposalStatusTestSuite) TestBridge_GetThreshold_Success() {
	s.mockContractCaller.EXPECT().From().Return(common.HexToAddress(testInteractorAddress))
	s.mockContractCaller.EXPECT().CallContract(
		gomock.Any(),
		gomock.Any(),
		nil,
	).Return([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2}, nil)
	res, err := s.bridgeContract.GetThreshold()
	s.Equal(
		uint8(2),
		res,
	)
	s.Nil(err)
}

func (s *ProposalStatusTestSuite) TestBridge_IsRelayer_Success() {
	s.mockContractCaller.EXPECT().From().Return(common.HexToAddress(testInteractorAddress))
	s.mockContractCaller.EXPECT().CallContract(
		gomock.Any(),
		gomock.Any(),
		nil,
	).Return([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, nil)
	res, err := s.bridgeContract.IsRelayer(common.HexToAddress(testInteractorAddress))
	s.Equal(
		true,
		res,
	)
	s.Nil(err)
}

func (s *ProposalStatusTestSuite) TestBridge_ProposalStatus_Success() {
	proposalStatus, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000001c0000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000001f")
	s.mockContractCaller.EXPECT().From().Return(common.HexToAddress(testInteractorAddress))
	s.mockContractCaller.EXPECT().CallContract(
		gomock.Any(),
		gomock.Any(),
		nil,
	).Return(proposalStatus, nil)
	res, err := s.bridgeContract.ProposalStatus(&s.proposal)
	s.Equal(
		message.ProposalStatus(message.ProposalStatus{Status: 0x3, YesVotes: big.NewInt(28), YesVotesTotal: 0x3, ProposedBlock: big.NewInt(31)}),
		res,
	)
	s.Nil(err)
}

func (s *ProposalStatusTestSuite) TestBridge_IsProposalVotedBy_Success() {
	s.mockContractCaller.EXPECT().From().Return(common.HexToAddress(testInteractorAddress))
	s.mockContractCaller.EXPECT().CallContract(
		gomock.Any(),
		gomock.Any(),
		nil,
	).Return([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}, nil)
	res, err := s.bridgeContract.IsProposalVotedBy(common.HexToAddress(testHandlerAddress), &s.proposal)
	s.Equal(
		true,
		res,
	)
	s.Nil(err)
}

func (s *ProposalStatusTestSuite) TestBridge_GetHandlerAddressForResourceID_Success() {
	s.mockContractCaller.EXPECT().From().Return(common.HexToAddress(testInteractorAddress))
	s.mockContractCaller.EXPECT().CallContract(
		gomock.Any(),
		gomock.Any(),
		nil,
	).Return([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 3, 4, 5}, nil)
	res, err := s.bridgeContract.GetHandlerAddressForResourceID(testResourceId)
	s.Equal(
		common.HexToAddress("0x0000000000000000000000000000000102030405"),
		res,
	)
	s.Nil(err)
}
