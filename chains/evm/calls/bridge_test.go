package calls_test

import (
	"encoding/hex"
	"errors"
	"testing"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	mock_listener "github.com/ChainSafe/chainbridge-core/chains/evm/calls/mock"
	"github.com/ChainSafe/chainbridge-core/chains/evm/voter/proposal"
	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

func TestPrepareSetDepositNonceInput(t *testing.T) {
	domainId := uint8(0)
	depositNonce := uint64(0)

	bytes, err := calls.PrepareSetDepositNonceInput(domainId, depositNonce)
	if err != nil {
		t.Fatalf("could not prepare set deposit nonce input: %v", err)
	}

	if len(bytes) == 0 {
		t.Fatal("byte slice returned is empty")
	}
}

type ProposalStatusTestSuite struct {
	suite.Suite
	mockContractCaller *mock_listener.MockContractCallerClient
}

func TestRunProposalStatusTestSuite(t *testing.T) {
	suite.Run(t, new(ProposalStatusTestSuite))
}

func (s *ProposalStatusTestSuite) SetupSuite()    {}
func (s *ProposalStatusTestSuite) TearDownSuite() {}
func (s *ProposalStatusTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	s.mockContractCaller = mock_listener.NewMockContractCallerClient(gomockController)
}
func (s *ProposalStatusTestSuite) TearDownTest() {}

func (s *ProposalStatusTestSuite) TestProposalStatusFailedContractCall() {
	s.mockContractCaller.EXPECT().CallContract(gomock.Any(), gomock.Any(), nil).Return(nil, errors.New("error"))

	status, err := calls.ProposalStatus(s.mockContractCaller, &proposal.Proposal{})

	s.Equal(relayer.ProposalStatusInactive, status)
	s.NotNil(err)
}

func (s *ProposalStatusTestSuite) TestProposalStatusFailedUnpack() {
	s.mockContractCaller.EXPECT().CallContract(gomock.Any(), gomock.Any(), nil).Return([]byte("invalid"), nil)

	status, err := calls.ProposalStatus(s.mockContractCaller, &proposal.Proposal{})

	s.Equal(relayer.ProposalStatusInactive, status)
	s.NotNil(err)
}

func (s *ProposalStatusTestSuite) TestProposalStatusSuccessfulCall() {
	proposalStatus, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")
	s.mockContractCaller.EXPECT().CallContract(gomock.Any(), gomock.Any(), nil).Return(proposalStatus, nil)

	status, err := calls.ProposalStatus(s.mockContractCaller, &proposal.Proposal{})

	s.Equal(relayer.ProposalStatusInactive, status)
	s.Nil(err)
}
