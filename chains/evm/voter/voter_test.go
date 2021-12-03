package voter_test

import (
	"errors"
	"testing"
	"time"

	"github.com/ChainSafe/chainbridge-core/chains/evm/voter"
	mock_voter "github.com/ChainSafe/chainbridge-core/chains/evm/voter/mock"
	"github.com/ChainSafe/chainbridge-core/chains/evm/voter/proposal"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type VoterTestSuite struct {
	suite.Suite
	voter              *voter.EVMVoter
	mockMessageHandler *mock_voter.MockMessageHandler
	mockClient         *mock_voter.MockChainClient
	mockBridgeContract *mock_voter.MockBridgeContract
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
	s.mockBridgeContract = mock_voter.NewMockBridgeContract(gomockController)
	s.voter = voter.NewVoter(
		s.mockMessageHandler,
		s.mockClient,
		s.mockBridgeContract,
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

	s.mockBridgeContract.EXPECT().IsProposalVotedBy(gomock.Any(), gomock.Any()).Return(false, nil)
	s.mockBridgeContract.EXPECT().ProposalStatus(gomock.Any()).Return(message.ProposalStatus{Status: message.ProposalStatusActive}, nil)
	s.mockBridgeContract.EXPECT().GetThreshold().Return(uint8(1), nil)
	s.mockBridgeContract.EXPECT().SimulateVoteProposal(gomock.Any()).Times(6).Return(errors.New("error"))

	err := s.voter.VoteProposal(&message.Message{})

	s.NotNil(err)
}

func (s *VoterTestSuite) TestVoteProposal_SimulateVoteProposal() {
	s.mockMessageHandler.EXPECT().HandleMessage(gomock.Any()).Return(&proposal.Proposal{
		Source:       0,
		DepositNonce: 0,
	}, nil)
	s.mockClient.EXPECT().RelayerAddress().Return(common.Address{})

	s.mockBridgeContract.EXPECT().IsProposalVotedBy(gomock.Any(), gomock.Any()).Return(false, nil)
	s.mockBridgeContract.EXPECT().ProposalStatus(gomock.Any()).Return(message.ProposalStatus{Status: message.ProposalStatusActive}, nil)
	s.mockBridgeContract.EXPECT().GetThreshold().Return(uint8(1), nil)
	s.mockBridgeContract.EXPECT().SimulateVoteProposal(gomock.Any()).Times(1).Return(nil)
	s.mockBridgeContract.EXPECT().VoteProposal(gomock.Any(), gomock.Any()).Return(&common.Hash{}, nil)

	err := s.voter.VoteProposal(&message.Message{})

	s.Nil(err)
}

func (s *VoterTestSuite) TestVoteProposal_IsProposalVotedByError() {
	s.mockMessageHandler.EXPECT().HandleMessage(gomock.Any()).Return(&proposal.Proposal{
		Source:       0,
		DepositNonce: 0,
	}, nil)
	s.mockClient.EXPECT().RelayerAddress().Return(common.Address{})
	s.mockBridgeContract.EXPECT().IsProposalVotedBy(gomock.Any(), gomock.Any()).Return(false, errors.New("error"))

	err := s.voter.VoteProposal(&message.Message{})

	s.NotNil(err)
}

func (s *VoterTestSuite) TestVoteProposal_ProposalAlreadyVoted() {
	s.mockMessageHandler.EXPECT().HandleMessage(gomock.Any()).Return(&proposal.Proposal{
		Source:       0,
		DepositNonce: 0,
	}, nil)
	s.mockClient.EXPECT().RelayerAddress().Return(common.Address{})
	s.mockBridgeContract.EXPECT().IsProposalVotedBy(gomock.Any(), gomock.Any()).Return(true, nil)

	err := s.voter.VoteProposal(&message.Message{})

	s.Nil(err)
}

func (s *VoterTestSuite) TestVoteProposal_ProposalStatusFail() {
	s.mockMessageHandler.EXPECT().HandleMessage(gomock.Any()).Return(&proposal.Proposal{
		Source:       0,
		DepositNonce: 0,
	}, nil)
	s.mockClient.EXPECT().RelayerAddress().Return(common.Address{})
	s.mockBridgeContract.EXPECT().IsProposalVotedBy(gomock.Any(), gomock.Any()).Return(false, nil)
	s.mockBridgeContract.EXPECT().ProposalStatus(gomock.Any()).Return(message.ProposalStatus{}, errors.New("error"))

	err := s.voter.VoteProposal(&message.Message{})

	s.NotNil(err)
}

func (s *VoterTestSuite) TestVoteProposal_ExecutedProposal() {
	s.mockMessageHandler.EXPECT().HandleMessage(gomock.Any()).Return(&proposal.Proposal{
		Source:       0,
		DepositNonce: 0,
	}, nil)
	s.mockClient.EXPECT().RelayerAddress().Return(common.Address{})
	s.mockBridgeContract.EXPECT().IsProposalVotedBy(gomock.Any(), gomock.Any()).Return(false, nil)
	s.mockBridgeContract.EXPECT().ProposalStatus(gomock.Any()).Return(message.ProposalStatus{Status: message.ProposalStatusExecuted}, nil)

	err := s.voter.VoteProposal(&message.Message{})

	s.Nil(err)
}

func (s *VoterTestSuite) TestVoteProposal_GetThresholdFail() {
	s.mockMessageHandler.EXPECT().HandleMessage(gomock.Any()).Return(&proposal.Proposal{
		Source:       0,
		DepositNonce: 0,
	}, nil)
	s.mockClient.EXPECT().RelayerAddress().Return(common.Address{})
	s.mockBridgeContract.EXPECT().IsProposalVotedBy(gomock.Any(), gomock.Any()).Return(false, nil)
	s.mockBridgeContract.EXPECT().ProposalStatus(gomock.Any()).Return(message.ProposalStatus{Status: message.ProposalStatusActive}, nil)
	s.mockBridgeContract.EXPECT().GetThreshold().Return(uint8(0), errors.New("error"))

	err := s.voter.VoteProposal(&message.Message{})

	s.NotNil(err)
}
