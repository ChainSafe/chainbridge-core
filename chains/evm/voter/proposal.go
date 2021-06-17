package voter

import (
	"context"
	"math/big"
	"strconv"
	"strings"

	"github.com/status-im/keycard-go/hexutils"

	"github.com/ethereum/go-ethereum/accounts/abi"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
)

type Proposal struct {
	Source         uint8  // Source where message was initiated
	Destination    uint8  // Destination chain of message
	DepositNonce   uint64 // Nonce for the deposit
	ResourceId     [32]byte
	Payload        []interface{} // data associated with event sequence
	Data           []byte
	HandlerAddress common.Address
	BridgeAddress  common.Address
}

func (p *Proposal) Status(evmCaller ChainClient) (relayer.ProposalStatus, error) {
	definition := "[{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"originChainID\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"depositNonce\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"}],\"name\":\"getProposal\",\"outputs\":[{\"components\":[{\"internalType\":\"bytes32\",\"name\":\"_resourceID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"_dataHash\",\"type\":\"bytes32\"},{\"internalType\":\"address[]\",\"name\":\"_yesVotes\",\"type\":\"address[]\"},{\"internalType\":\"address[]\",\"name\":\"_noVotes\",\"type\":\"address[]\"},{\"internalType\":\"enumBridge.ProposalStatus\",\"name\":\"_status\",\"type\":\"uint8\"},{\"internalType\":\"uint256\",\"name\":\"_proposedBlock\",\"type\":\"uint256\"}],\"internalType\":\"structBridge.Proposal\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
	a, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		return relayer.ProposalStatusInactive, err // Not sure what status to use here
	}
	log.Debug().Str("datahash", p.GetDataHash().String()).Uint8("source", p.Source).Uint64("nonce", p.DepositNonce).Msg("Getting proposal status")
	input, err := a.Pack("getProposal", p.Source, p.DepositNonce, p.GetDataHash())
	if err != nil {
		return relayer.ProposalStatusInactive, err
	}
	log.Debug().Msg(hexutils.BytesToHex(input))

	msg := ethereum.CallMsg{From: common.Address{}, To: &p.BridgeAddress, Data: input}
	out, err := evmCaller.CallContract(context.TODO(), toCallArg(msg), nil)
	log.Debug().Msg(strconv.Itoa(len(out)))
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
	log.Debug().Msgf("GETTING PROPOSAL STATUS %+v", out0)
	return relayer.ProposalStatus(out0.Status), nil
}

func (p *Proposal) VotedBy(evmCaller ChainClient, by common.Address) (bool, error) {
	definition := "[{\"inputs\":[{\"internalType\":\"uint72\",\"name\":\"\",\"type\":\"uint72\"},{\"internalType\":\"bytes32\",\"name\":\"\",\"type\":\"bytes32\"},{\"internalType\":\"address\",\"name\":\"\",\"type\":\"address\"}],\"name\":\"_hasVotedOnProposal\",\"outputs\":[{\"internalType\":\"bool\",\"name\":\"\",\"type\":\"bool\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
	//input, err := buildDataUnsafe([]byte("_hasVotedOnProposal(uint72,bytes32,address)"), idAndNonce(p.Source, p.DepositNonce).Bytes(), p.GetDataHash().Bytes(), by.Bytes())
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

func (p *Proposal) Execute(client ChainClient) error {
	definition := "[{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"chainID\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"depositNonce\",\"type\":\"uint64\"},{\"internalType\":\"bytes\",\"name\":\"data\",\"type\":\"bytes\"},{\"internalType\":\"bytes32\",\"name\":\"resourceID\",\"type\":\"bytes32\"}],\"name\":\"executeProposal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
	a, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		return err // Not sure what status to use here
	}
	input, err := a.Pack("executeProposal", p.Source, p.DepositNonce, p.Data, p.ResourceId)
	if err != nil {
		return err
	}
	gasLimit := uint64(6721975)
	gp, err := client.GasPrice()
	if err != nil {
		return err
	}
	tx := evmtransaction.NewTransaction(client.Nonce(), p.BridgeAddress, big.NewInt(0), gasLimit, gp, input)
	h, err := client.SignAndSendTransaction(context.TODO(), tx)
	if err != nil {
		return err
	}
	log.Debug().Str("hash", h.Hex()).Msgf("Executed")
	return nil
}

func (p *Proposal) Vote(client ChainClient) error {
	definition := "[{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"chainID\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"depositNonce\",\"type\":\"uint64\"},{\"internalType\":\"bytes32\",\"name\":\"resourceID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes32\",\"name\":\"dataHash\",\"type\":\"bytes32\"}],\"name\":\"voteProposal\",\"outputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"function\"}]"
	a, err := abi.JSON(strings.NewReader(definition))
	if err != nil {
		return err // Not sure what status to use here
	}
	input, err := a.Pack("voteProposal", p.Source, p.DepositNonce, p.ResourceId, p.GetDataHash())
	if err != nil {
		return err
	}
	gasLimit := uint64(6721975)
	gp, err := client.GasPrice()
	if err != nil {
		return err
	}
	tx := evmtransaction.NewTransaction(client.Nonce(), p.BridgeAddress, big.NewInt(0), gasLimit, gp, input)
	h, err := client.SignAndSendTransaction(context.TODO(), tx)
	if err != nil {
		return err
	}
	log.Debug().Str("hash", h.Hex()).Msgf("Executed")
	return nil
}

// CreateProposalDataHash constructs and returns proposal data hash
func (p *Proposal) GetDataHash() common.Hash {
	log.Debug().Str("handler", p.HandlerAddress.String()).Str("data", hexutils.BytesToHex(p.Data)).Msg("Creating datahash")
	return crypto.Keccak256Hash(append(p.HandlerAddress.Bytes(), p.Data...))
}

func idAndNonce(srcId uint8, nonce uint64) *big.Int {
	var data []byte
	data = append(data, big.NewInt(int64(nonce)).Bytes()...)
	data = append(data, uint8(srcId))
	return big.NewInt(0).SetBytes(data)
}
