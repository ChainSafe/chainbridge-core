package calls_test

import (
	"encoding/hex"
	"errors"
	"math/big"
	"testing"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	mock_listener "github.com/ChainSafe/chainbridge-core/chains/evm/calls/mock"
	"github.com/ChainSafe/chainbridge-core/chains/evm/voter/proposal"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/ethereum/go-ethereum/common"
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

	s.Equal(message.ProposalStatus{}, status)
	s.NotNil(err)
}

func (s *ProposalStatusTestSuite) TestProposalStatusFailedUnpack() {
	s.mockContractCaller.EXPECT().CallContract(gomock.Any(), gomock.Any(), nil).Return([]byte("invalid"), nil)

	status, err := calls.ProposalStatus(s.mockContractCaller, &proposal.Proposal{})

	s.Equal(message.ProposalStatus{}, status)
	s.NotNil(err)
}

func (s *ProposalStatusTestSuite) TestProposalStatusSuccessfulCall() {
	proposalStatus, _ := hex.DecodeString("0000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000001c0000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000001f")
	s.mockContractCaller.EXPECT().CallContract(gomock.Any(), gomock.Any(), nil).Return(proposalStatus, nil)

	status, err := calls.ProposalStatus(s.mockContractCaller, &proposal.Proposal{})

	s.Nil(err)
	s.Equal(status.YesVotesTotal, uint8(3))
	s.Equal(status.Status, message.ProposalStatusExecuted)
}

func TestPrepareWithdrawInput(t *testing.T) {
	handlerAddress := common.HexToAddress("0x3167776db165D8eA0f51790CA2bbf44Db5105ADF")
	tokenAddress := common.HexToAddress("0x3f709398808af36ADBA86ACC617FeB7F5B7B193E")
	recipientAddress := common.HexToAddress("0x8e5F72B158BEDf0ab50EDa78c70dFC118158C272")
	amountOrTokenId := big.NewInt(1)

	inputBytes, err := calls.PrepareWithdrawInput(
		handlerAddress,
		tokenAddress,
		recipientAddress,
		amountOrTokenId,
	)
	if err != nil {
		t.Fatalf("could not prepare withdraw input: %v", err)
	}

	if len(inputBytes) == 0 {
		t.Fatal("prepared input byte slice empty")
	}
}
