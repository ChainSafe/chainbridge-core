package bridge

import (
	"bytes"
	"context"
	"fmt"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/client"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"math/big"
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/voter/proposal"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/ChainSafe/chainbridge-core/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

type BridgeContract struct {
	client                client.ContractCallerDispatcherClient
	bridgeContractAddress common.Address
	abi                   abi.ABI
	transactor.Transactor
}

func NewBridgeContract(client client.ContractCallerDispatcherClient, bridgeContractAddress common.Address, transactor transactor.Transactor) *BridgeContract {
	a, err := abi.JSON(strings.NewReader(consts.BridgeABI))
	if err != nil {
		log.Fatal().Msg("Unable to load BridgeABI") // TODO
	}

	return &BridgeContract{
		client:                client,
		bridgeContractAddress: bridgeContractAddress,
		abi:                   a,
		Transactor:            transactor,
	}
}

func (c *BridgeContract) PackMethod(method string, args ...interface{}) ([]byte, error) {
	input, err := c.abi.Pack(method, args...)
	if err != nil {
		return []byte{}, err
	}
	return input, nil
}

func (c *BridgeContract) UnpackResult(method string, output []byte) ([]interface{}, error) {
	res, err := c.abi.Unpack(method, output)
	if err != nil {
		log.Error().Err(fmt.Errorf("unpack output error: %v", err))
		return nil, err
	}
	return res, err
}

// public function to generate bytedata for adminWithdraw contract method
// Used to manually withdraw funds from ERC safes
func (c *BridgeContract) PrepareWithdrawInput(
	handlerAddress,
	tokenAddress,
	recipientAddress common.Address,
	realAmount *big.Int,
) ([]byte, error) {
	// @dev withdrawal data should include:
	// tokenAddress
	// recipientAddress
	// realAmount
	data := bytes.Buffer{}
	data.Write(common.LeftPadBytes(tokenAddress.Bytes(), 32))
	data.Write(common.LeftPadBytes(recipientAddress.Bytes(), 32))
	data.Write(common.LeftPadBytes(realAmount.Bytes(), 32))

	return c.PackMethod("adminWithdraw", handlerAddress, data.Bytes())
}

// ------------------

func (c *BridgeContract) AddRelayer(
	relayerAddr common.Address,
	opts transactor.TransactOptions,
) (*common.Hash, error) {
	return c.ExecuteTransaction(
		"adminAddRelayer",
		opts,
		relayerAddr,
	)
}

func (c *BridgeContract) AdminSetGenericResource(
	handler common.Address,
	rID types.ResourceID,
	addr common.Address,
	depositFunctionSig [4]byte,
	depositerOffset *big.Int,
	executeFunctionSig [4]byte,
	opts transactor.TransactOptions,
) (*common.Hash, error) {
	return c.ExecuteTransaction(
		"adminSetGenericResource",
		opts,
		handler, rID, addr, depositFunctionSig, depositerOffset, executeFunctionSig,
	)
}

func (c *BridgeContract) AdminSetResource(
	handlerAddr common.Address,
	resourceIdBytesArr types.ResourceID,
	targetContractAddr common.Address,
	opts transactor.TransactOptions,
) (*common.Hash, error) {
	return c.ExecuteTransaction(
		"adminSetResource",
		opts,
		handlerAddr, resourceIdBytesArr, targetContractAddr,
	)
}

func (c *BridgeContract) SetDepositNonce(
	domainId uint8,
	depositNonce uint64,
	opts transactor.TransactOptions,
) (*common.Hash, error) {
	return c.ExecuteTransaction(
		"adminSetDepositNonce",
		opts,
		domainId, depositNonce,
	)
}

func (c *BridgeContract) SetThresholdInput(
	threshold uint64,
	opts transactor.TransactOptions,
) (*common.Hash, error) {
	return c.ExecuteTransaction(
		"adminChangeRelayerThreshold",
		opts,
		big.NewInt(0).SetUint64(threshold),
	)
}

func (c *BridgeContract) SetBurnableInput(
	handlerAddr common.Address,
	tokenContractAddr common.Address,
	opts transactor.TransactOptions,
) (*common.Hash, error) {
	return c.ExecuteTransaction(
		"adminSetBurnable",
		opts,
		handlerAddr, tokenContractAddr,
	)
}

func (c *BridgeContract) Deposit(
	resourceID types.ResourceID,
	destDomainID uint8,
	data []byte,
	opts transactor.TransactOptions,
) (*common.Hash, error) {
	return c.ExecuteTransaction(
		"deposit",
		opts,
		destDomainID, resourceID, data,
	)
}

func (c *BridgeContract) Erc721Deposit(
	tokenId *big.Int,
	metadata string,
	recipient common.Address,
	resourceID types.ResourceID,
	destDomainID uint8,
	opts transactor.TransactOptions,
) (*common.Hash, error) {
	data := ConstructErc721DepositData(recipient.Bytes(), tokenId, []byte(metadata))
	txHash, err := c.Deposit(resourceID, destDomainID, data, opts)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	// TODO - log success
	return txHash, err
}

func (c *BridgeContract) ExecuteProposal(
	proposal *proposal.Proposal,
	opts transactor.TransactOptions,
) (*common.Hash, error) {
	opts = transactor.MergeTransactionOptions(opts, transactor.TransactOptions{GasLimit: 300000})
	// revertOnFail should be constantly false, true is used only for internal contract calls
	// when you need to execute proposal in voteProposal function right after it becomes Passed becouse of votes
	return c.ExecuteTransaction(
		"executeProposal",
		opts,
		proposal.Source, proposal.DepositNonce, proposal.ResourceId, proposal.Data, true,
	)
}

func (c *BridgeContract) VoteProposal(
	proposal *proposal.Proposal,
	opts transactor.TransactOptions,
) (*common.Hash, error) {
	opts = transactor.MergeTransactionOptions(opts, transactor.TransactOptions{GasLimit: 300000})
	// revertOnFail should be constantly false, true is used only for internal contract calls when you need to execute proposal in voteProposal function right after it becomes Passed becouse of votes
	return c.ExecuteTransaction(
		"voteProposal",
		opts,
		proposal.Source, proposal.DepositNonce, proposal.ResourceId, proposal.Data,
	)
}

// public function to Withdraw funds from ERC safes
func (c *BridgeContract) Withdraw(
	handlerAddress,
	tokenAddress,
	recipientAddress common.Address,
	amountOrTokenId *big.Int,
	opts transactor.TransactOptions,
) (*common.Hash, error) {
	withdrawInput, err := c.PrepareWithdrawInput(
		handlerAddress, tokenAddress, recipientAddress, amountOrTokenId,
	)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	return c.executeTransaction(withdrawInput, opts, "Withdraw")
}

func (c *BridgeContract) GetThreshold() (uint8, error) {
	res, err := c.CallContract("_relayerThreshold")
	if err != nil {
		return 0, err
	}
	t := *abi.ConvertType(res[0], new(uint8)).(*uint8)
	return t, nil
}

func (c *BridgeContract) IsRelayer(relayerAddress common.Address) (bool, error) {
	res, err := c.CallContract("isRelayer", relayerAddress)
	if err != nil {
		return false, err
	}
	b := abi.ConvertType(res[0], new(bool)).(*bool)
	return *b, nil
}

func (c *BridgeContract) ProposalStatus(p *proposal.Proposal) (message.ProposalStatus, error) {
	res, err := c.CallContract("getProposal", p.Source, p.DepositNonce, p.GetDataHash())
	if err != nil {
		return message.ProposalStatus{}, err
	}
	ps := *abi.ConvertType(res[0], new(message.ProposalStatus)).(*message.ProposalStatus)
	return ps, nil
}

func (c *BridgeContract) IsProposalVotedBy(by common.Address, p *proposal.Proposal) (bool, error) {
	res, err := c.CallContract("_hasVotedOnProposal", idAndNonce(p.Source, p.DepositNonce), p.GetDataHash(), by)
	if err != nil {
		return false, err
	}
	out0 := *abi.ConvertType(res[0], new(bool)).(*bool)
	return out0, nil
}

func (c *BridgeContract) CallContract(method string, args ...interface{}) ([]interface{}, error) {
	input, err := c.PackMethod(method, args...)
	if err != nil {
		return nil, err
	}
	msg := ethereum.CallMsg{From: common.Address{}, To: &c.bridgeContractAddress, Data: input}
	out, err := c.client.CallContract(context.TODO(), client.ToCallArg(msg), nil)
	if err != nil {
		return nil, err
	}
	if len(out) == 0 {
		// Make sure we have a contract to operate on, and bail out otherwise.
		if code, err := c.client.CodeAt(context.Background(), c.bridgeContractAddress, nil); err != nil {
			return nil, err
		} else if len(code) == 0 {
			return nil, fmt.Errorf("no code at provided address %s", c.bridgeContractAddress.String())
		}
	}
	return c.UnpackResult(method, out)
}

func (c BridgeContract) ExecuteTransaction(method string, opts transactor.TransactOptions, args ...interface{}) (*common.Hash, error) {
	input, err := c.PackMethod(method, args...)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	return c.executeTransaction(input, opts, method)
}

func (c *BridgeContract) executeTransaction(
	input []byte,
	opts transactor.TransactOptions,
	opName string,
) (*common.Hash, error) {
	h, err := c.Transact(&c.bridgeContractAddress, input, opts)
	if err != nil {
		log.Error().Err(err).Msg(opName)
		return nil, err
	}
	log.Debug().Str("hash", h.String()).Msgf("%s sent", opName)
	return h, err
}

func idAndNonce(srcId uint8, nonce uint64) *big.Int {
	var data []byte
	data = append(data, big.NewInt(int64(nonce)).Bytes()...)
	data = append(data, uint8(srcId))
	return big.NewInt(0).SetBytes(data)
}
