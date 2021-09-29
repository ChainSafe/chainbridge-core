package calls

import (
	"fmt"
	"github.com/ChainSafe/chainbridge-core/chains/evm/voter"
	"math/big"
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

func PrepareSetBurnableInput(handler, tokenAddress common.Address) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(consts.BridgeABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("adminSetBurnable", handler, tokenAddress)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func PrepareAdminSetResourceInput(handler common.Address, rId [32]byte, addr common.Address) ([]byte, error) {
	log.Debug().Msgf("ResourceID %x", rId)
	a, err := abi.JSON(strings.NewReader(consts.BridgeABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("adminSetResource", handler, rId, addr)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func PrepareErc20DepositInput(destDomainID uint8, resourceID [32]byte, data []byte) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(consts.BridgeABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("deposit", destDomainID, resourceID, data)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func PrepareExecuteProposalInput(sourceDomainID uint8, depositNonce uint64, resourceID [32]byte, calldata []byte, revertOnFail bool) ([]byte, error){
	a, err := abi.JSON(strings.NewReader(consts.BridgeABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("executeProposal", sourceDomainID, depositNonce, calldata, resourceID, revertOnFail)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func PrepareVoteProposalInput(sourceDomainID uint8, resourceID [32]byte, calldata []byte) ([]byte, error){
	a, err := abi.JSON(strings.NewReader(consts.BridgeABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("voteProposal", sourceDomainID, resourceID, calldata)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func PrepareAddRelayerInput(relayer common.Address) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(consts.BridgeABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("adminAddRelayer", relayer)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}
func PrepareIsRelayerInput(address common.Address) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(consts.BridgeABI))
	if err != nil {
		return nil, err
	}

	data, err := a.Pack("isRelayer", address)
	if err != nil {
		log.Error().Err(fmt.Errorf("unpack output error: %v", err))
		return nil, err
	}
	return data, nil
}

func ParseIsRelayerOutput(output []byte) (bool, error) {
	a, err := abi.JSON(strings.NewReader(consts.BridgeABI))
	if err != nil {
		return false, err
	}

	res, err := a.Unpack("isRelayer", output)
	if err != nil {
		log.Error().Err(fmt.Errorf("unpack output error: %v", err))
		return false, err
	}

	b := abi.ConvertType(res[0], new(bool)).(*bool)
	return *b, nil
}

func Deposit(client ChainClient, fabric TxFabric, bridgeAddress, recipient common.Address, amount *big.Int, resourceID [32]byte, destDomainID uint8) error {
	data := ConstructErc20DepositData(recipient.Bytes(), amount)
	input, err := PrepareErc20DepositInput(destDomainID, resourceID, data)
	if err != nil {
		return err
	}
	gasLimit := uint64(2000000)
	h, err := Transact(client, fabric, &bridgeAddress, input, gasLimit)
	if err != nil {
		return fmt.Errorf("deposit failed %w", err)
	}
	log.Debug().Str("hash", h.String()).Msgf("Deposit sent")
	return nil
}


func ExecuteProposal(client ClientDispatcher, fabric TxFabric, proposal *voter.Proposal) (common.Hash, error) {
	// revertOnFail should be constantly false, true is used only for internal contract calls when you need to execute proposal in voteProposal function right after it becomes Passed becouse of votes
	input, err := PrepareExecuteProposalInput(proposal.Source, proposal.DepositNonce, proposal.ResourceId, proposal.Data, true)
	if err != nil {
		return common.Hash{}, err
	}
	gasLimit := uint64(300000)
	h, err := Transact(client, fabric, &proposal.BridgeAddress, input, gasLimit)
	if err != nil {
		return common.Hash{}, fmt.Errorf("execute proposal failed %w", err)
	}
	return h, nil
}


func VoteProposal(client ClientDispatcher, fabric TxFabric, proposal *voter.Proposal) (common.Hash, error) {
	// revertOnFail should be constantly false, true is used only for internal contract calls when you need to execute proposal in voteProposal function right after it becomes Passed becouse of votes
	input, err := PrepareVoteProposalInput(proposal.Source, proposal.ResourceId, proposal.Data)
	if err != nil {
		return common.Hash{}, err
	}
	gasLimit := uint64(300000)
	h, err := Transact(client, fabric, &proposal.BridgeAddress, input, gasLimit)
	if err != nil {
		return common.Hash{}, fmt.Errorf("vote proposal failed %w", err)
	}
	return h, nil
}
