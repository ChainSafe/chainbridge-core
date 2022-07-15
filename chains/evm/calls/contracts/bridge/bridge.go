package bridge

import (
	"bytes"
	"math/big"
	"strconv"
	"strings"

	"github.com/ChainSafe/sygma-core/chains/evm/calls"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/consts"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts/deposit"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/transactor"
	"github.com/ethereum/go-ethereum/common/hexutil"

	"github.com/ChainSafe/sygma-core/chains/evm/executor/proposal"
	"github.com/ChainSafe/sygma-core/relayer/message"
	"github.com/ChainSafe/sygma-core/types"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/rs/zerolog/log"
)

type BridgeContract struct {
	contracts.Contract
}

func NewBridgeContract(
	client calls.ContractCallerDispatcher,
	bridgeContractAddress common.Address,
	transactor transactor.Transactor,
) *BridgeContract {
	a, _ := abi.JSON(strings.NewReader(consts.BridgeABI))
	b := common.FromHex(consts.BridgeBin)
	return &BridgeContract{contracts.NewContract(bridgeContractAddress, a, b, client, transactor)}
}

func (c *BridgeContract) AddRelayer(
	relayerAddr common.Address,
	opts transactor.TransactOptions,
) (*common.Hash, error) {
	log.Debug().Msgf("Adding new relayer %s", relayerAddr.String())
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
	log.Debug().Msgf("Setting generic resource %s", hexutil.Encode(rID[:]))
	return c.ExecuteTransaction(
		"adminSetGenericResource",
		opts,
		handler, rID, addr, depositFunctionSig, depositerOffset, executeFunctionSig,
	)
}

func (c *BridgeContract) AdminSetResource(
	handlerAddr common.Address,
	rID types.ResourceID,
	targetContractAddr common.Address,
	opts transactor.TransactOptions,
) (*common.Hash, error) {
	log.Debug().Msgf("Setting resource %s", hexutil.Encode(rID[:]))
	return c.ExecuteTransaction(
		"adminSetResource",
		opts,
		handlerAddr, rID, targetContractAddr,
	)
}

func (c *BridgeContract) SetDepositNonce(
	domainId uint8,
	depositNonce uint64,
	opts transactor.TransactOptions,
) (*common.Hash, error) {
	log.Debug().Msgf("Setting deposit nonce %d for %d", depositNonce, domainId)
	return c.ExecuteTransaction(
		"adminSetDepositNonce",
		opts,
		domainId, depositNonce,
	)
}

func (c *BridgeContract) AdminChangeRelayerThreshold(
	threshold uint64,
	opts transactor.TransactOptions,
) (*common.Hash, error) {
	log.Debug().Msgf("Setting threshold %d", threshold)
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
	log.Debug().Msgf("Setting burnable input for %s", tokenContractAddr.String())
	return c.ExecuteTransaction(
		"adminSetBurnable",
		opts,
		handlerAddr, tokenContractAddr,
	)
}

func (c *BridgeContract) deposit(
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

func (c *BridgeContract) Erc20Deposit(
	recipient common.Address,
	amount *big.Int,
	resourceID types.ResourceID,
	destDomainID uint8,
	opts transactor.TransactOptions,
) (*common.Hash, error) {
	log.Debug().
		Str("recipient", recipient.String()).
		Str("resourceID", hexutil.Encode(resourceID[:])).
		Str("amount", amount.String()).
		Msgf("ERC20 deposit")
	var data []byte
	if opts.Priority == 0 {
		data = deposit.ConstructErc20DepositData(recipient.Bytes(), amount)
	} else {
		data = deposit.ConstructErc20DepositDataWithPriority(recipient.Bytes(), amount, opts.Priority)
	}
	txHash, err := c.deposit(resourceID, destDomainID, data, opts)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
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
	log.Debug().
		Str("recipient", recipient.String()).
		Str("resourceID", hexutil.Encode(resourceID[:])).
		Str("tokenID", tokenId.String()).
		Msgf("ERC721 deposit")
	var data []byte
	if opts.Priority == 0 {
		data = deposit.ConstructErc721DepositData(recipient.Bytes(), tokenId, []byte(metadata))
	} else {
		data = deposit.ConstructErc721DepositDataWithPriority(recipient.Bytes(), tokenId, []byte(metadata), opts.Priority)
	}
	txHash, err := c.deposit(resourceID, destDomainID, data, opts)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	return txHash, err
}

func (c *BridgeContract) GenericDeposit(
	metadata []byte,
	resourceID types.ResourceID,
	destDomainID uint8,
	opts transactor.TransactOptions,
) (*common.Hash, error) {
	log.Debug().
		Str("resourceID", hexutil.Encode(resourceID[:])).
		Msgf("Generic deposit")
	data := deposit.ConstructGenericDepositData(metadata)
	txHash, err := c.deposit(resourceID, destDomainID, data, opts)
	if err != nil {
		log.Error().Err(err)
		return nil, err
	}
	return txHash, err
}

func (c *BridgeContract) ExecuteProposal(
	proposal *proposal.Proposal,
	opts transactor.TransactOptions,
) (*common.Hash, error) {
	log.Debug().
		Str("depositNonce", strconv.FormatUint(proposal.DepositNonce, 10)).
		Str("resourceID", hexutil.Encode(proposal.ResourceId[:])).
		Str("handler", proposal.HandlerAddress.String()).
		Msgf("Execute proposal")
	return c.ExecuteTransaction(
		"executeProposal",
		opts,
		proposal.Source, proposal.DepositNonce, proposal.Data, proposal.ResourceId, true,
	)
}

func (c *BridgeContract) VoteProposal(
	proposal *proposal.Proposal,
	opts transactor.TransactOptions,
) (*common.Hash, error) {
	log.Debug().
		Str("depositNonce", strconv.FormatUint(proposal.DepositNonce, 10)).
		Str("resourceID", hexutil.Encode(proposal.ResourceId[:])).
		Str("handler", proposal.HandlerAddress.String()).
		Msgf("Vote proposal")
	return c.ExecuteTransaction(
		"voteProposal",
		opts,
		proposal.Source, proposal.DepositNonce, proposal.ResourceId, proposal.Data,
	)
}

func (c *BridgeContract) SimulateVoteProposal(proposal *proposal.Proposal) error {
	log.Debug().
		Str("depositNonce", strconv.FormatUint(proposal.DepositNonce, 10)).
		Str("resourceID", hexutil.Encode(proposal.ResourceId[:])).
		Str("handler", proposal.HandlerAddress.String()).
		Msgf("Simulate vote proposal")
	_, err := c.CallContract(
		"voteProposal",
		proposal.Source, proposal.DepositNonce, proposal.ResourceId, proposal.Data,
	)
	return err
}

func (c *BridgeContract) Pause(opts transactor.TransactOptions) (*common.Hash, error) {
	log.Debug().Msg("Pause transfers")
	return c.ExecuteTransaction(
		"adminPauseTransfers",
		opts,
	)
}

func (c *BridgeContract) Unpause(opts transactor.TransactOptions) (*common.Hash, error) {
	log.Debug().Msg("Unpause transfers")
	return c.ExecuteTransaction(
		"adminUnpauseTransfers",
		opts,
	)
}

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
	log.Debug().Msg("Getting threshold")
	res, err := c.CallContract("_relayerThreshold")
	if err != nil {
		return 0, err
	}
	out := *abi.ConvertType(res[0], new(uint8)).(*uint8)
	return out, nil
}

func (c *BridgeContract) IsRelayer(relayerAddress common.Address) (bool, error) {
	log.Debug().Msgf("Getting is %s a relayer", relayerAddress.String())
	res, err := c.CallContract("isRelayer", relayerAddress)
	if err != nil {
		return false, err
	}
	out := abi.ConvertType(res[0], new(bool)).(*bool)
	return *out, nil
}

func (c *BridgeContract) ProposalStatus(p *proposal.Proposal) (message.ProposalStatus, error) {
	log.Debug().
		Str("depositNonce", strconv.FormatUint(p.DepositNonce, 10)).
		Str("resourceID", hexutil.Encode(p.ResourceId[:])).
		Str("handler", p.HandlerAddress.String()).
		Msg("Getting proposal status")
	res, err := c.CallContract("getProposal", p.Source, p.DepositNonce, p.GetDataHash())
	if err != nil {
		return message.ProposalStatus{}, err
	}
	out := *abi.ConvertType(res[0], new(message.ProposalStatus)).(*message.ProposalStatus)
	return out, nil
}

func (c *BridgeContract) IsProposalVotedBy(by common.Address, p *proposal.Proposal) (bool, error) {
	log.Debug().
		Str("depositNonce", strconv.FormatUint(p.DepositNonce, 10)).
		Str("resourceID", hexutil.Encode(p.ResourceId[:])).
		Str("handler", p.HandlerAddress.String()).
		Msgf("Getting is proposal voted by %s", by.String())
	res, err := c.CallContract("_hasVotedOnProposal", idAndNonce(p.Source, p.DepositNonce), p.GetDataHash(), by)
	if err != nil {
		return false, err
	}
	out := *abi.ConvertType(res[0], new(bool)).(*bool)
	return out, nil
}

func (c *BridgeContract) GetHandlerAddressForResourceID(
	resourceID types.ResourceID,
) (common.Address, error) {
	log.Debug().Msgf("Getting handler address for resource %s", hexutil.Encode(resourceID[:]))
	res, err := c.CallContract("_resourceIDToHandlerAddress", resourceID)
	if err != nil {
		return common.Address{}, err
	}
	out := *abi.ConvertType(res[0], new(common.Address)).(*common.Address)
	return out, nil
}

func idAndNonce(srcId uint8, nonce uint64) *big.Int {
	var data []byte
	data = append(data, big.NewInt(int64(nonce)).Bytes()...)
	data = append(data, uint8(srcId))
	return big.NewInt(0).SetBytes(data)
}
