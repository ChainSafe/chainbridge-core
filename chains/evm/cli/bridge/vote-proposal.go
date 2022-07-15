package bridge

import (
	"fmt"

	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts/bridge"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/flags"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/initialize"
	"github.com/ChainSafe/sygma-core/chains/evm/cli/logger"
	"github.com/ChainSafe/sygma-core/chains/evm/executor/proposal"
	"github.com/ChainSafe/sygma-core/util"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var voteProposalCmd = &cobra.Command{
	Use:   "vote-proposal",
	Short: "Vote on proposal",
	Long:  "The vote-proposal subcommand votes on an on-chain proposal. Valid relayer private key required for transaction to be successful.",
	PreRun: func(cmd *cobra.Command, args []string) {
		logger.LoggerMetadata(cmd.Name(), cmd.Flags())
	},
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		return util.CallPersistentPreRun(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		c, err := initialize.InitializeClient(url, senderKeyPair)
		if err != nil {
			return err
		}
		t, err := initialize.InitializeTransactor(gasPrice, evmtransaction.NewTransaction, c, prepare)
		if err != nil {
			return err
		}
		return VoteProposalCmd(cmd, args, bridge.NewBridgeContract(c, BridgeAddr, t))
	},
	Args: func(cmd *cobra.Command, args []string) error {
		err := ValidateVoteProposalFlags(cmd, args)
		if err != nil {
			return err
		}

		err = ProcessVoteProposalFlags(cmd, args)
		return err
	},
}

func BindVoteProposalCmdFlags(cmd *cobra.Command) {
	cmd.Flags().StringVar(&Bridge, "bridge", "", "bridge contract address")
	cmd.Flags().StringVar(&Data, "data", "", "hex proposal metadata")
	cmd.Flags().Uint8Var(&DomainID, "domain", 0, "origin domain ID of proposal")
	cmd.Flags().Uint64Var(&DepositNonce, "deposit-nonce", 0, "deposit nonce of proposal to vote on")
	cmd.Flags().StringVar(&ResourceID, "resource", "", "resource id of asset")
	flags.MarkFlagsAsRequired(cmd, "bridge", "deposit-nonce", "domain", "resource", "data")
}

func init() {
	BindVoteProposalCmdFlags(voteProposalCmd)
}

func ValidateVoteProposalFlags(cmd *cobra.Command, args []string) error {
	if !common.IsHexAddress(Bridge) {
		return fmt.Errorf("invalid bridge address: %s", Bridge)
	}
	return nil
}

func ProcessVoteProposalFlags(cmd *cobra.Command, args []string) error {
	var err error
	BridgeAddr = common.HexToAddress(Bridge)
	DataBytes = common.Hex2Bytes(Data)

	ResourceIdBytesArr, err = flags.ProcessResourceID(ResourceID)
	return err
}

type Voter interface {
	SimulateVoteProposal(proposal *proposal.Proposal) error
	VoteProposal(proposal *proposal.Proposal, opts transactor.TransactOptions) (*common.Hash, error)
}

func VoteProposalCmd(cmd *cobra.Command, args []string, voter Voter) error {
	prop := &proposal.Proposal{
		Source:       uint8(DomainID),
		DepositNonce: DepositNonce,
		Data:         DataBytes,
		ResourceId:   ResourceIdBytesArr,
	}

	err := voter.SimulateVoteProposal(prop)
	if err != nil {
		return err
	}

	h, err := voter.VoteProposal(prop, transactor.TransactOptions{})
	if err != nil {
		return err
	}

	log.Info().Msgf("Successfully voted on proposal with hash: %s", h.Hex())
	return nil
}
