package bridge

import (
	"bytes"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/client"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/deposit"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"math/big"
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/voter/proposal"
	"github.com/ChainSafe/chainbridge-core/relayer/message"
	"github.com/ChainSafe/chainbridge-core/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

type BridgeContract struct {
	contracts.Contract
}

func NewBridgeContract(
	client client.ContractCallerDispatcherClient,
	bridgeContractAddress common.Address,
	transactor transactor.Transactor,
) *BridgeContract {
	a, err := abi.JSON(strings.NewReader(consts.BridgeABI))
	if err != nil {
		log.Fatal().Msg("Unable to load BridgeABI") // TODO
	}
	b := common.FromHex(consts.BridgeBin)
	return &BridgeContract{contracts.NewContract(bridgeContractAddress, a, b, client, transactor)}
}

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

func (c BridgeContract) Erc20Deposit(
	recipient common.Address,
	amount *big.Int,
	resourceID types.ResourceID,
	destDomainID uint8,
	opts transactor.TransactOptions,
) (*common.Hash, error) {
	data := deposit.ConstructErc20DepositData(recipient.Bytes(), amount)
	txHash, err := c.Deposit(resourceID, destDomainID, data, opts)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	// TODO - log success
	return txHash, err
}

func (c *BridgeContract) Erc721Deposit(
	tokenId *big.Int,
	metadata string,
	recipient common.Address,
	resourceID types.ResourceID,
	destDomainID uint8,
	opts transactor.TransactOptions,
) (*common.Hash, error) {
	data := deposit.ConstructErc721DepositData(recipient.Bytes(), tokenId, []byte(metadata))
	txHash, err := c.Deposit(resourceID, destDomainID, data, opts)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	// TODO - log success
	return txHash, err
}

func (c BridgeContract) GenericDeposit(
	metadata []byte,
	resourceID types.ResourceID,
	destDomainID uint8,
	opts transactor.TransactOptions,
) (*common.Hash, error) {
	data := deposit.ConstructGenericDepositData(metadata)
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
	return c.ExecuteTransaction(
		"voteProposal",
		opts,
		proposal.Source, proposal.DepositNonce, proposal.ResourceId, proposal.Data,
	)
}

func (c *BridgeContract) SimulateVoteProposal(proposal *proposal.Proposal) error {
	_, err := c.CallContract(
		"voteProposal",
		proposal.Source, proposal.DepositNonce, proposal.ResourceId, proposal.Data,
	)
	return err
}

func (c *BridgeContract) Pause(opts transactor.TransactOptions) (*common.Hash, error) {
	return c.ExecuteTransaction(
		"adminPauseTransfers",
		opts,
	)
}

func (c *BridgeContract) Unpause(opts transactor.TransactOptions) (*common.Hash, error) {
	return c.ExecuteTransaction(
		"adminUnpauseTransfers",
		opts,
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
	// @dev withdrawal data should include:
	// tokenAddress
	// recipientAddress
	// realAmount
	data := bytes.Buffer{}
	data.Write(common.LeftPadBytes(tokenAddress.Bytes(), 32))
	data.Write(common.LeftPadBytes(recipientAddress.Bytes(), 32))
	data.Write(common.LeftPadBytes(amountOrTokenId.Bytes(), 32))

	return c.ExecuteTransaction("adminWithdraw", opts, handlerAddress, data.Bytes())
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

func idAndNonce(srcId uint8, nonce uint64) *big.Int {
	var data []byte
	data = append(data, big.NewInt(int64(nonce)).Bytes()...)
	data = append(data, uint8(srcId))
	return big.NewInt(0).SetBytes(data)
}
