package voter_test

import (
	"encoding/hex"
	"errors"
	mock_calls "github.com/ChainSafe/chainbridge-core/chains/evm/calls/mock"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm/voter"
	mock_voter "github.com/ChainSafe/chainbridge-core/chains/evm/voter/mock"
	"github.com/ChainSafe/chainbridge-core/chains/evm/voter/proposal"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"math/big"
	"testing"
	"time"
)

var (
	proposalVotedResponse, _    = hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000001")
	threshold, _                = hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000003")
	proposalNotVotedResponse, _ = hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000")
	executedProposalStatus, _   = hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000001c0000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000001f")
	inactiveProposalStatus, _   = hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
)

type VoterTestSuite struct {
	suite.Suite
	voter              *voter.EVMVoter
	mockMessageHandler *mock_voter.MockMessageHandler
	mockClient         *mock_voter.MockChainClient
	mockgasPricer      *mock_calls.MockGasPricer
}

func TestRunVoterTestSuite(t *testing.T) {
	suite.Run(t, new(VoterTestSuite))
}

func (s *VoterTestSuite) SetupSuite()    {}
func (s *VoterTestSuite) TearDownSuite() {}
func (s *VoterTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	s.mockMessageHandler = mock_voter.NewMockMessageHandler(gomockController)
	s.mockClient = mock_voter.NewMockChainClient(gomockController)
	s.mockgasPricer = mock_calls.NewMockGasPricer(gomockController)
	s.voter = voter.NewVoter(
		s.mockMessageHandler,
		s.mockClient,
		evmtransaction.NewTransaction,
		s.mockgasPricer,
	)
	voter.Sleep = func(d time.Duration) {}
}
func (s *VoterTestSuite) TearDownTest() {}

func (s *VoterTestSuite) TestVoteProposal_HandleMessageError() {
	s.mockMessageHandler.EXPECT().HandleMessage(gomock.Any()).Return(nil, errors.New("error"))

	err := s.voter.VoteProposal(&message.Message{})

	s.NotNil(err)
}

func (s *VoterTestSuite) TestVoteProposal_SimulateVoteProposalError() {
	s.mockMessageHandler.EXPECT().HandleMessage(gomock.Any()).Return(&proposal.Proposal{
		Source:       0,
		DepositNonce: 0,
	}, nil)
	s.mockClient.EXPECT().RelayerAddress().Return(common.Address{})
	//calls.IsProposalVotedBy
	s.mockClient.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Return(proposalNotVotedResponse, nil)
	//calls.ProposalStatus
	s.mockClient.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Return(inactiveProposalStatus, nil)
	//calls.GetThreshold
	s.mockClient.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Return(threshold, nil)

	//calls.repetitiveSimulateVote called 6 times
	s.mockClient.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Times(6).Return(nil, errors.New("error"))

	err := s.voter.VoteProposal(&message.Message{})

	s.NotNil(err)
}

func (s *VoterTestSuite) TestVoteProposal_SimulateVoteProposal() {
	s.mockMessageHandler.EXPECT().HandleMessage(gomock.Any()).Return(&proposal.Proposal{
		Source:       0,
		DepositNonce: 0,
	}, nil)
	s.mockClient.EXPECT().RelayerAddress().Return(common.Address{})
	s.mockClient.EXPECT().LockNonce().Times(1)
	s.mockClient.EXPECT().UnlockNonce().Times(1)
	s.mockClient.EXPECT().UnsafeNonce().Return(big.NewInt(1), nil)
	s.mockgasPricer.EXPECT().GasPrice().Return([]*big.Int{big.NewInt(10)}, nil)
	s.mockClient.EXPECT().SignAndSendTransaction(gomock.Any(), gomock.Any()).Times(1).Return(common.Hash{}, nil)
	s.mockClient.EXPECT().From().Return(common.Address{})
	s.mockClient.EXPECT().WaitAndReturnTxReceipt(gomock.Any()).Times(1).Return(nil, nil)
	s.mockClient.EXPECT().UnsafeIncreaseNonce().Times(1)

	//calls.IsProposalVotedBy
	s.mockClient.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Return(proposalNotVotedResponse, nil)
	//calls.ProposalStatus
	s.mockClient.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Return(inactiveProposalStatus, nil)
	//calls.GetThreshold
	s.mockClient.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Return(threshold, nil)
	//calls.repetitiveSimulateVote
	s.mockClient.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Return(proposalVotedResponse, nil)

	err := s.voter.VoteProposal(&message.Message{})

	s.Nil(err)
}
func (s *VoterTestSuite) TestVoteProposal_IsProposalVotedByError() {
	s.mockMessageHandler.EXPECT().HandleMessage(gomock.Any()).Return(&proposal.Proposal{
		Source:       0,
		DepositNonce: 0,
	}, nil)
	s.mockClient.EXPECT().RelayerAddress().Return(common.Address{})
	s.mockClient.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte{}, errors.New("error"))

	err := s.voter.VoteProposal(&message.Message{})

	s.NotNil(err)
}

func (s *VoterTestSuite) TestVoteProposal_ProposalAlreadyVoted() {
	s.mockMessageHandler.EXPECT().HandleMessage(gomock.Any()).Return(&proposal.Proposal{
		Source:       0,
		DepositNonce: 0,
	}, nil)
	s.mockClient.EXPECT().RelayerAddress().Return(common.Address{})
	s.mockClient.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Return(proposalVotedResponse, nil)

	err := s.voter.VoteProposal(&message.Message{})

	s.Nil(err)
}

func (s *VoterTestSuite) TestVoteProposal_ProposalStatusFail() {
	s.mockMessageHandler.EXPECT().HandleMessage(gomock.Any()).Return(&proposal.Proposal{
		Source:       0,
		DepositNonce: 0,
	}, nil)
	s.mockClient.EXPECT().RelayerAddress().Return(common.Address{})
	s.mockClient.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Return(proposalNotVotedResponse, nil)
	s.mockClient.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte{}, errors.New("error"))

	err := s.voter.VoteProposal(&message.Message{})

	s.NotNil(err)
}

func (s *VoterTestSuite) TestVoteProposal_ExecutedProposal() {
	s.mockMessageHandler.EXPECT().HandleMessage(gomock.Any()).Return(&proposal.Proposal{
		Source:       0,
		DepositNonce: 0,
	}, nil)
	s.mockClient.EXPECT().RelayerAddress().Return(common.Address{})
	s.mockClient.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Return(proposalNotVotedResponse, nil)
	s.mockClient.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Return(executedProposalStatus, nil)

	err := s.voter.VoteProposal(&message.Message{})

	s.Nil(err)
}

func (s *VoterTestSuite) TestVoteProposal_GetThresholdFail() {
	s.mockMessageHandler.EXPECT().HandleMessage(gomock.Any()).Return(&proposal.Proposal{
		Source:       0,
		DepositNonce: 0,
	}, nil)
	s.mockClient.EXPECT().RelayerAddress().Return(common.Address{})
	s.mockClient.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Return(proposalNotVotedResponse, nil)
	s.mockClient.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Return(inactiveProposalStatus, nil)
	s.mockClient.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte{}, errors.New("error"))

	err := s.voter.VoteProposal(&message.Message{})

	s.NotNil(err)
}
