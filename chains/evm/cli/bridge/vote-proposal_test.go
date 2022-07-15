package bridge_test

import (
	"encoding/hex"
	"errors"
	"testing"

	"github.com/ChainSafe/sygma-core/chains/evm/calls"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/sygma-core/chains/evm/cli"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/bridge"
	mock_bridge "github.com/ChainSafe/sygma-core/chains/evm/cli/bridge/mock"
	"github.com/ChainSafe/sygma-core/chains/evm/executor/proposal"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/spf13/cobra"
	"github.com/stretchr/testify/suite"
)

type VoteProposalTestSuite struct {
	suite.Suite
	gomockController    *gomock.Controller
	mockVoteProposalCmd *cobra.Command
	voter               *mock_bridge.MockVoter
}

func TestRunVoteProposalTestSuite(t *testing.T) {
	suite.Run(t, new(VoteProposalTestSuite))
}

func (s *VoteProposalTestSuite) SetupSuite()    {}
func (s *VoteProposalTestSuite) TearDownSuite() {}
func (s *VoteProposalTestSuite) SetupTest() {
	s.gomockController = gomock.NewController(s.T())
	s.voter = mock_bridge.NewMockVoter(s.gomockController)

	s.mockVoteProposalCmd = &cobra.Command{
		RunE: func(cmd *cobra.Command, args []string) error {
			return bridge.VoteProposalCmd(
				cmd,
				args,
				s.voter)
		},
		Args: func(cmd *cobra.Command, args []string) error {
			err := bridge.ValidateVoteProposalFlags(cmd, args)
			if err != nil {
				return err
			}

			err = bridge.ProcessVoteProposalFlags(cmd, args)
			return err
		},
	}
	cli.BindEVMCLIFlags(s.mockVoteProposalCmd)
	bridge.BindVoteProposalCmdFlags(s.mockVoteProposalCmd)
}

func (s *VoteProposalTestSuite) TestValidate_InvalidBridgeAddress() {
	rootCmdArgs := []string{
		"--url", "test-url",
		"--private-key", "test-private-key",
		"--bridge", "invalid",
		"--data", "hex-data",
		"--resource", "test-resource",
		"--domain", "1",
	}
	s.mockVoteProposalCmd.SetArgs(rootCmdArgs)

	err := s.mockVoteProposalCmd.Execute()

	s.NotNil(err)
	s.Equal(err.Error(), "invalid bridge address: invalid")
}

func (s *VoteProposalTestSuite) TestValidate_InvalidResourceID() {
	rootCmdArgs := []string{
		"--url", "test-url",
		"--private-key", "test-private-key",
		"--bridge", "0x829bd824b016326a401d083b33d092293333a830",
		"--data", "hex-data",
		"--resource", "111",
		"--domain", "1",
	}
	s.mockVoteProposalCmd.SetArgs(rootCmdArgs)

	err := s.mockVoteProposalCmd.Execute()

	s.NotNil(err)
	s.Equal(err.Error(), "failed decoding resourceID hex string: encoding/hex: odd length hex string")
}

func (s *VoteProposalTestSuite) TestValidate_FailedSimulateCall() {
	rootCmdArgs := []string{
		"--url", "test-url",
		"--private-key", "test-private-key",
		"--bridge", "0x829bd824b016326a401d083b33d092293333a830",
		"--data", "00000000000000000000000000000000000000000000000000000000000f424000000000000000000000000000000000000000000000000000000000000000148e0a907331554af72563bd8d43051c2e64be5d35",
		"--resource", "0x000000000000000000000075df75bcdca8ea2360c562b4aadbaf3dfaf5b19b00",
		"--deposit-nonce", "2",
		"--domain", "1",
	}
	s.mockVoteProposalCmd.SetArgs(rootCmdArgs)
	resourceID, _ := hex.DecodeString("000000000000000000000075df75bcdca8ea2360c562b4aadbaf3dfaf5b19b00")
	s.voter.EXPECT().SimulateVoteProposal(&proposal.Proposal{
		Source:       1,
		DepositNonce: 2,
		Data:         common.Hex2Bytes("00000000000000000000000000000000000000000000000000000000000f424000000000000000000000000000000000000000000000000000000000000000148e0a907331554af72563bd8d43051c2e64be5d35"),
		ResourceId:   calls.SliceTo32Bytes(resourceID),
	}).Return(errors.New("failed simulating call"))

	err := s.mockVoteProposalCmd.Execute()

	s.NotNil(err)
	s.Equal(err.Error(), "failed simulating call")
}

func (s *VoteProposalTestSuite) TestValidate_FailedVoteCall() {
	rootCmdArgs := []string{
		"--url", "test-url",
		"--private-key", "test-private-key",
		"--bridge", "0x829bd824b016326a401d083b33d092293333a830",
		"--data", "00000000000000000000000000000000000000000000000000000000000f424000000000000000000000000000000000000000000000000000000000000000148e0a907331554af72563bd8d43051c2e64be5d35",
		"--resource", "0x000000000000000000000075df75bcdca8ea2360c562b4aadbaf3dfaf5b19b00",
		"--deposit-nonce", "2",
		"--domain", "1",
	}
	s.mockVoteProposalCmd.SetArgs(rootCmdArgs)
	resourceID, _ := hex.DecodeString("000000000000000000000075df75bcdca8ea2360c562b4aadbaf3dfaf5b19b00")
	s.voter.EXPECT().SimulateVoteProposal(&proposal.Proposal{
		Source:       1,
		DepositNonce: 2,
		Data:         common.Hex2Bytes("00000000000000000000000000000000000000000000000000000000000f424000000000000000000000000000000000000000000000000000000000000000148e0a907331554af72563bd8d43051c2e64be5d35"),
		ResourceId:   calls.SliceTo32Bytes(resourceID),
	}).Return(nil)
	s.voter.EXPECT().VoteProposal(&proposal.Proposal{
		Source:       1,
		DepositNonce: 2,
		Data:         common.Hex2Bytes("00000000000000000000000000000000000000000000000000000000000f424000000000000000000000000000000000000000000000000000000000000000148e0a907331554af72563bd8d43051c2e64be5d35"),
		ResourceId:   calls.SliceTo32Bytes(resourceID),
	}, transactor.TransactOptions{}).Return(&common.Hash{}, errors.New("failed vote call"))

	err := s.mockVoteProposalCmd.Execute()

	s.NotNil(err)
	s.Equal(err.Error(), "failed vote call")
}

func (s *VoteProposalTestSuite) TestValidate_SuccessfulVote() {
	rootCmdArgs := []string{
		"--url", "test-url",
		"--private-key", "test-private-key",
		"--bridge", "0x829bd824b016326a401d083b33d092293333a830",
		"--data", "00000000000000000000000000000000000000000000000000000000000f424000000000000000000000000000000000000000000000000000000000000000148e0a907331554af72563bd8d43051c2e64be5d35",
		"--resource", "0x000000000000000000000075df75bcdca8ea2360c562b4aadbaf3dfaf5b19b00",
		"--deposit-nonce", "2",
		"--domain", "1",
	}
	s.mockVoteProposalCmd.SetArgs(rootCmdArgs)
	resourceID, _ := hex.DecodeString("000000000000000000000075df75bcdca8ea2360c562b4aadbaf3dfaf5b19b00")
	s.voter.EXPECT().SimulateVoteProposal(&proposal.Proposal{
		Source:       1,
		DepositNonce: 2,
		Data:         common.Hex2Bytes("00000000000000000000000000000000000000000000000000000000000f424000000000000000000000000000000000000000000000000000000000000000148e0a907331554af72563bd8d43051c2e64be5d35"),
		ResourceId:   calls.SliceTo32Bytes(resourceID),
	}).Return(nil)
	s.voter.EXPECT().VoteProposal(&proposal.Proposal{
		Source:       1,
		DepositNonce: 2,
		Data:         common.Hex2Bytes("00000000000000000000000000000000000000000000000000000000000f424000000000000000000000000000000000000000000000000000000000000000148e0a907331554af72563bd8d43051c2e64be5d35"),
		ResourceId:   calls.SliceTo32Bytes(resourceID),
	}, transactor.TransactOptions{}).Return(&common.Hash{}, nil)

	err := s.mockVoteProposalCmd.Execute()

	s.Nil(err)
}
