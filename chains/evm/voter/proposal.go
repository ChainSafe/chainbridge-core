package voter

import (
	"context"
	"math/big"
	"strings"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"

	"github.com/status-im/keycard-go/hexutils"

	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
)

type TxFabric func(nonce uint64, to *common.Address, amount *big.Int, gasLimit uint64, gasPrice *big.Int, data []byte) evmclient.CommonTransaction

func NewProposal(source uint8, depositNonce uint64, resourceId [32]byte, data []byte, handlerAddress, bridgeAddress common.Address) *Proposal {
	return &Proposal{
		Source:         source,
		DepositNonce:   depositNonce,
		ResourceId:     resourceId,
		Data:           data,
		HandlerAddress: handlerAddress,
		BridgeAddress:  bridgeAddress,
	}
}

type Proposal struct {
	Source         uint8  // Source where message was initiated
	DepositNonce   uint64 // Nonce for the deposit
	ResourceId     [32]byte
	Payload        []interface{} // data associated with event sequence
	Data           []byte
	HandlerAddress common.Address
	BridgeAddress  common.Address
}

func (p *Proposal) Status(evmCaller ChainClient) (relayer.ProposalStatus, error) {
	definition := "[{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"originDomainID\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"depositNonce\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"}],\"name\":\"getProposal\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"_resourceID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_dataHash\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"_yesVotes\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_noVotes\",\"type\":\"address[]\"},{\"internalType\":\"enumBridge.ProposalStatus\",\"name\":\"_status\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"_proposedBlock\",\"type\":\"uint256\"}],\"internalType\":\"structBridge.Proposal\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
	a, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		return relayer.ProposalStatusInactive, err // Not sure what status to use here
	}
	input, err := a.Pack("getProposal", p.Source, p.DepositNonce, p.GetDataHash())
	if err != nil {
		return relayer.ProposalStatusInactive, err
	}

	msg := ethereum.CallMsg{From: common.Address{}, To: &p.BridgeAddress, Data: input}
	out, err := evmCaller.CallContract(context.TODO(), toCallArg(msg), nil)
	if err != nil {
		return relayer.ProposalStatusInactive, err
	}
	type bridgeProposal struct {
		ResourceID    [32]byte
		DataHash      [32]byte
		YesVotes      []common.Address
		NoVotes       []common.Address
		Status        uint8
		ProposedBlock *big.Int
	}
	res, err := a.Unpack("getProposal", out)
	out0 := *abi.ConvertType(res[0], new(bridgeProposal)).(*bridgeProposal)
	return relayer.ProposalStatus(out0.Status), nil
}

func (p *Proposal) VotedBy(evmCaller ChainClient, by common.Address) (bool, error) {
	definition := "[{\"inputs\":[{\"internalType\":\"uint72\",\"name\":\"\",\"type\":\"uint72\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"_hasVotedOnProposal\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
	a, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		return false, err // Not sure what status to use here
	}
	input, err := a.Pack("_hasVotedOnProposal", idAndNonce(p.Source, p.DepositNonce), p.GetDataHash(), by)
	if err != nil {
		return false, err
	}
	msg := ethereum.CallMsg{From: common.Address{}, To: &p.BridgeAddress, Data: input}
	out, err := evmCaller.CallContract(context.TODO(), toCallArg(msg), nil)
	if err != nil {
		return false, err
	}
	res, err := a.Unpack("_hasVotedOnProposal", out)
	out0 := *abi.ConvertType(res[0], new(bool)).(*bool)
	return out0, nil
}

func (p *Proposal) Execute(client ChainClient, fabric TxFabric) error {
	log.Debug().Str("rID", hexutils.BytesToHex(p.ResourceId[:])).Uint64("depositNonce", p.DepositNonce).Msg("Executing proposal")
	definition := "[{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"domainID\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"depositNonce\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"resourceID\",\"type\":\"bytes32\"}],\"name\":\"executeProposal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
	a, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		return err // Not sure what status to use here
	}
	input, err := a.Pack("executeProposal", p.Source, p.DepositNonce, p.Data, p.ResourceId)
	if err != nil {
		return err
	}
	gasLimit := uint64(2000000)
	gp, err := client.GasPrice()
	if err != nil {
		return err
	}
	client.LockNonce()
	n, err := client.UnsafeNonce()
	if err != nil {
		return err
	}
	tx := fabric(n.Uint64(), &p.BridgeAddress, big.NewInt(0), gasLimit, gp, input)
	hash, err := client.SignAndSendTransaction(context.TODO(), tx)
	if err != nil {
		return err
	}
	log.Debug().Str("hash", hash.String()).Uint64("nonce", n.Uint64()).Msgf("Executed")
	err = client.UnsafeIncreaseNonce()
	if err != nil {
		return err
	}
	client.UnlockNonce()
	return nil
}

func (p *Proposal) Vote(client ChainClient, fabric TxFabric) error {
	log.Debug().Str("rID", hexutils.BytesToHex(p.ResourceId[:])).Uint64("depositNonce", p.DepositNonce).Uint8("domainID", p.Source).Msg("Voting proposal")
	definition := "[{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"domainID\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"depositNonce\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"resourceID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"}],\"name\":\"voteProposal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
	a, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		return err // Not sure what status to use here
	}
	input, err := a.Pack("voteProposal", p.Source, p.DepositNonce, p.ResourceId, p.GetDataHash())
	if err != nil {
		return err
	}
	gasLimit := uint64(1000000)
	gp, err := client.GasPrice()
	if err != nil {
		return err
	}
	client.LockNonce()
	n, err := client.UnsafeNonce()
	if err != nil {
		return err
	}
	tx := fabric(n.Uint64(), &p.BridgeAddress, big.NewInt(0), gasLimit, gp, input)
	hash, err := client.SignAndSendTransaction(context.TODO(), tx)
	if err != nil {
		return err
	}
	log.Debug().Str("hash", hash.String()).Uint64("nonce", n.Uint64()).Msgf("Voted")
	err = client.UnsafeIncreaseNonce()
	if err != nil {
		return err
	}
	client.UnlockNonce()
	return nil
}

// GetDataHash constructs and returns proposal data hash
func (p *Proposal) GetDataHash() common.Hash {
	return crypto.Keccak256Hash(append(p.HandlerAddress.Bytes(), p.Data...))
}

func idAndNonce(srcId uint8, nonce uint64) *big.Int {
	var data []byte
	data = append(data, big.NewInt(int64(nonce)).Bytes()...)
	data = append(data, uint8(srcId))
	return big.NewInt(0).SetBytes(data)
}
