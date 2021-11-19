package calls

import (
	"bytes"
	"context"
	"fmt"
	"math/big"
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ChainSafe/chainbridge-core/chains/evm/voter/proposal"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/ChainSafe/chainbridge-core/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
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

func PrepareAdminSetResourceInput(handler common.Address, resourceID types.ResourceID, addr common.Address) ([]byte, error) {
	log.Debug().Msgf("resourceID %x", resourceID)
	a, err := abi.JSON(strings.NewReader(consts.BridgeABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("adminSetResource", handler, resourceID, addr)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func PrepareAdminSetGenericResourceInput(
	handler common.Address,
	rId types.ResourceID,
	addr common.Address,
	depositFunctionSig [4]byte,
	depositerOffset *big.Int,
	executeFunctionSig [4]byte,
) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(consts.BridgeABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("adminSetGenericResource", handler, rId, addr, depositFunctionSig, depositerOffset, executeFunctionSig)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func AdminSetGenericResource(
	client ClientDispatcher,
	fabric TxFabric,
	gasPriceClient GasPricer,
	handler common.Address,
	rID types.ResourceID,
	addr common.Address,
	depositFunctionSig [4]byte,
	depositerOffset *big.Int,
	executeFunctionSig [4]byte,
) (common.Hash, error) {
	input, err := PrepareAdminSetGenericResourceInput(handler, rID, addr, depositFunctionSig, depositerOffset, executeFunctionSig)
	if err != nil {
		return common.Hash{}, err
	}

	gasLimit := uint64(2000000)
	h, err := Transact(client, fabric, gasPriceClient, &handler, input, gasLimit, big.NewInt(0))
	if err != nil {
		return common.Hash{}, fmt.Errorf("setting generic resource failed %w", err)
	}

	return h, nil
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

func ConstructErc20DepositData(destRecipient []byte, amount *big.Int) []byte {
	var data []byte
	data = append(data, math.PaddedBigBytes(amount, 32)...)
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(destRecipient))), 32)...)
	data = append(data, destRecipient...)
	return data
}

func ConstructErc721DepositData(destRecipient []byte, tokenId *big.Int, metadata []byte) []byte {
	var data []byte
	data = append(data, math.PaddedBigBytes(tokenId, 32)...)                               // Token Id
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(destRecipient))), 32)...) // Length of recipient
	data = append(data, destRecipient...)                                                  // Recipient
	data = append(data, metadata...)                                                       // Metadata
	return data
}

func ConstructGenericDepositData(metadata []byte) []byte {
	var data []byte
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(metadata))), 32)...)
	data = append(data, metadata...)
	return data
}

func PrepareDepositInput(destDomainID uint8, resourceID types.ResourceID, data []byte) ([]byte, error) {
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

func Deposit(
	client ClientDispatcher,
	fabric TxFabric,
	gasPriceClient GasPricer,
	bridgeAddress common.Address,
	resourceID types.ResourceID,
	destDomainID uint8,
	data []byte,
) (*common.Hash, error) {
	input, err := PrepareDepositInput(destDomainID, resourceID, data)
	if err != nil {
		return nil, err
	}

	gasLimit := uint64(2000000)
	h, err := Transact(client, fabric, gasPriceClient, &bridgeAddress, input, gasLimit, big.NewInt(0))
	if err != nil {
		return nil, fmt.Errorf("deposit failed %w", err)
	}
	log.Debug().Str("hash", h.String()).Msgf("Deposit sent")
	return &h, nil
}

func PrepareExecuteProposalInput(sourceDomainID uint8, depositNonce uint64, resourceID types.ResourceID, calldata []byte, revertOnFail bool) ([]byte, error) {
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

func ExecuteProposal(client ClientDispatcher, fabric TxFabric, gasPriceClient GasPricer, proposal *proposal.Proposal) (common.Hash, error) {
	// revertOnFail should be constantly false, true is used only for internal contract calls when you need to execute proposal in voteProposal function right after it becomes Passed becouse of votes
	input, err := PrepareExecuteProposalInput(proposal.Source, proposal.DepositNonce, proposal.ResourceId, proposal.Data, true)
	if err != nil {
		return common.Hash{}, err
	}
	gasLimit := uint64(300000)
	h, err := Transact(client, fabric, gasPriceClient, &proposal.BridgeAddress, input, gasLimit, big.NewInt(0))
	if err != nil {
		return common.Hash{}, fmt.Errorf("execute proposal failed %w", err)
	}
	return h, nil
}

func PrepareVoteProposalInput(sourceDomainID uint8, depositNonce uint64, resourceID types.ResourceID, calldata []byte) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(consts.BridgeABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("voteProposal", sourceDomainID, depositNonce, resourceID, calldata)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func VoteProposal(client ClientDispatcher, fabric TxFabric, gasPriceClient GasPricer, proposal *proposal.Proposal) (common.Hash, error) {
	// revertOnFail should be constantly false, true is used only for internal contract calls when you need to execute proposal in voteProposal function right after it becomes Passed becouse of votes
	input, err := PrepareVoteProposalInput(proposal.Source, proposal.DepositNonce, proposal.ResourceId, proposal.Data)
	if err != nil {
		return common.Hash{}, err
	}
	gasLimit := uint64(300000)
	h, err := Transact(client, fabric, gasPriceClient, &proposal.BridgeAddress, input, gasLimit, big.NewInt(0))
	if err != nil {
		return common.Hash{}, fmt.Errorf("vote proposal failed %w", err)
	}
	return h, nil
}

func PrepareSetDepositNonceInput(domainID uint8, depositNonce uint64) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(consts.BridgeABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("adminSetDepositNonce", domainID, depositNonce)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func PrepareSetThresholdInput(threshold *big.Int) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(consts.BridgeABI))
	if err != nil {
		return []byte{}, err
	}
	input, err := a.Pack("adminChangeRelayerThreshold", threshold)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func GetThreshold(evmCaller ContractCallerClient, bridgeAddress *common.Address) (uint8, error) {
	a, err := abi.JSON(strings.NewReader(consts.BridgeABI))
	if err != nil {
		return 0, err
	}
	input, err := a.Pack("_relayerThreshold")
	if err != nil {
		return 0, err
	}
	msg := ethereum.CallMsg{From: common.Address{}, To: bridgeAddress, Data: input}
	out, err := evmCaller.CallContract(context.TODO(), ToCallArg(msg), nil)
	if err != nil {
		return 0, err
	}

	res, err := a.Unpack("_relayerThreshold", out)
	if err != nil {
		return 0, err
	}

	out0 := *abi.ConvertType(res[0], new(uint8)).(*uint8)
	return out0, nil
}

func ProposalStatus(evmCaller ContractCallerClient, p *proposal.Proposal) (message.ProposalStatus, error) {
	a, err := abi.JSON(strings.NewReader(consts.BridgeABI))
	if err != nil {
		return message.ProposalStatus{}, err
	}
	input, err := a.Pack("getProposal", p.Source, p.DepositNonce, p.GetDataHash())
	if err != nil {
		return message.ProposalStatus{}, err
	}

	msg := ethereum.CallMsg{From: common.Address{}, To: &p.BridgeAddress, Data: input}
	out, err := evmCaller.CallContract(context.TODO(), ToCallArg(msg), nil)
	if err != nil {
		return message.ProposalStatus{}, err
	}
	res, err := a.Unpack("getProposal", out)
	if err != nil {
		return message.ProposalStatus{}, err
	}

	ps := *abi.ConvertType(res[0], new(message.ProposalStatus)).(*message.ProposalStatus)
	return ps, nil
}

func idAndNonce(srcId uint8, nonce uint64) *big.Int {
	var data []byte
	data = append(data, big.NewInt(int64(nonce)).Bytes()...)
	data = append(data, uint8(srcId))
	return big.NewInt(0).SetBytes(data)
}

func IsProposalVotedBy(evmCaller ContractCallerClient, by common.Address, p *proposal.Proposal) (bool, error) {
	a, err := abi.JSON(strings.NewReader(consts.BridgeABI))
	if err != nil {
		return false, err
	}
	input, err := a.Pack("_hasVotedOnProposal", idAndNonce(p.Source, p.DepositNonce), p.GetDataHash(), by)
	if err != nil {
		return false, err
	}
	msg := ethereum.CallMsg{From: common.Address{}, To: &p.BridgeAddress, Data: input}
	out, err := evmCaller.CallContract(context.TODO(), ToCallArg(msg), nil)
	if err != nil {
		return false, err
	}
	res, err := a.Unpack("_hasVotedOnProposal", out)
	if err != nil {
		return false, err
	}
	out0 := *abi.ConvertType(res[0], new(bool)).(*bool)
	return out0, nil
}

// public function to generate bytedata for adminWithdraw contract method
// Used to manually withdraw funds from ERC safes
func PrepareWithdrawInput(
	handlerAddress,
	tokenAddress,
	recipientAddress common.Address,
	realAmount *big.Int,
) ([]byte, error) {
	a, err := abi.JSON(strings.NewReader(consts.BridgeABI))
	if err != nil {
		return []byte{}, err
	}

	// @dev withdrawal data should include:
	// tokenAddress
	// recipientAddress
	// realAmount
	data := bytes.Buffer{}
	data.Write(common.LeftPadBytes(tokenAddress.Bytes(), 32))
	data.Write(common.LeftPadBytes(recipientAddress.Bytes(), 32))
	data.Write(common.LeftPadBytes(realAmount.Bytes(), 32))

	input, err := a.Pack(
		"adminWithdraw",
		handlerAddress,
		data.Bytes(),
	)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

// public function to Withdraw funds from ERC safes
func Withdraw(client ClientDispatcher, txFabric TxFabric, gasPricer GasPricer, gasLimit uint64, bridgeAddress, handlerAddress, tokenAddress, recipientAddress common.Address, amountOrTokenId *big.Int) (*common.Hash, error) {
	withdrawInput, err := PrepareWithdrawInput(
		handlerAddress,
		tokenAddress,
		recipientAddress,
		amountOrTokenId,
	)
	if err != nil {
		return nil, fmt.Errorf("withdrawal input error: %v", err)
	}
	h, err := Transact(
		client,
		txFabric,
		gasPricer,
		&bridgeAddress,
		withdrawInput,
		gasLimit,
		big.NewInt(0),
	)
	if err != nil {
		return nil, fmt.Errorf("withdrawal failed %w", err)
	}
	log.Debug().Str("hash", h.String()).Msgf("Withdrawal sent")

	return &h, nil
}
