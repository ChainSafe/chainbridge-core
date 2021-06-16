package voter

import (
	"context"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/relayer"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
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
	//function getProposal(uint8 originChainID, uint64 depositNonce, bytes32 dataHash) view returns((bytes32,bytes32,address[],address[],uint8,uint256))
	input, err := buildDataUnsafe([]byte("getProposal(uint8,uint64,bytes32"), big.NewInt(0).SetUint64(uint64(p.Source)).Bytes(), big.NewInt(0).SetUint64(p.DepositNonce).Bytes(), p.GetDataHash().Bytes())
	if err != nil {
		return relayer.ProposalStatusActive, err // Not sure what status to use here
	}
	msg := ethereum.CallMsg{From: common.Address{}, To: &p.BridgeAddress, Data: input}
	out, err := evmCaller.CallContract(context.TODO(), toCallArg(msg), nil)
	if err != nil {
		return relayer.ProposalStatusActive, err
	}
	type bridgeProposal struct {
		ResourceID    [32]byte
		DataHash      [32]byte
		YesVotes      []common.Address
		NoVotes       []common.Address
		Status        uint8
		ProposedBlock *big.Int
	}

	out0 := *abi.ConvertType(out[0], new(bridgeProposal)).(*bridgeProposal)
	return relayer.ProposalStatus(out0.Status), nil
}

func (p *Proposal) VotedBy(evmCaller ChainClient, by common.Address) (bool, error) {
	//_hasVotedOnProposal(uint72 , bytes32 , address ) constant returns(bool)
	input, err := buildDataUnsafe([]byte("_hasVotedOnProposal(uint72,bytes32,address"), idAndNonce(p.Source, p.DepositNonce).Bytes(), p.GetDataHash().Bytes(), by.Bytes())
	if err != nil {
		return false, err // Not sure what status to use here
	}
	msg := ethereum.CallMsg{From: common.Address{}, To: &p.BridgeAddress, Data: input}
	out, err := evmCaller.CallContract(context.TODO(), toCallArg(msg), nil)
	if err != nil {
		return false, err
	}
	var b bool
	out0 := *abi.ConvertType(out[0], b).(*bool)
	return out0, nil
}

func (p *Proposal) Execute(client ChainClient) error {
	gasLimit := uint64(6721975)
	//executeProposal(uint8 chainID, uint64 depositNonce, bytes data, bytes32 resourceID) returns()
	data, err := buildDataUnsafe(
		[]byte("executeProposal(uint8,uint64,bytes,bytes32)"),
		big.NewInt(int64(p.Source)).Bytes(),
		big.NewInt(int64(p.DepositNonce)).Bytes(),
		p.Data,
		p.ResourceId[:])

	if err != nil {
		return err
	}
	gp, err := client.GasPrice()
	if err != nil {
		return err
	}
	tx := evmtransaction.NewTransaction(client.Nonce(), p.BridgeAddress, big.NewInt(0), gasLimit, gp, data)
	h, err := client.SignAndSendTransaction(context.TODO(), tx)
	if err != nil {
		return err
	}
	log.Debug().Str("hash", h.Hex()).Msgf("Executed")
	return nil
}

func (p *Proposal) Vote(client ChainClient) error {
	gasLimit := uint64(6721975)
	//voteProposal(uint8 chainID, uint64 depositNonce, bytes32 resourceID, bytes32 dataHash)
	data, err := buildDataUnsafe(
		[]byte("voteProposal(uint8,uint64,bytes,bytes32,bytes32)"),
		big.NewInt(int64(p.Source)).Bytes(),
		big.NewInt(int64(p.DepositNonce)).Bytes(),
		p.ResourceId[:],
		p.GetDataHash().Bytes(),
	)
	if err != nil {
		return err
	}
	gp, err := client.GasPrice()
	if err != nil {
		return err
	}
	tx := evmtransaction.NewTransaction(client.Nonce(), p.BridgeAddress, big.NewInt(0), gasLimit, gp, data)

	h, err := client.SignAndSendTransaction(context.TODO(), tx)
	if err != nil {
		return err
	}
	log.Debug().Str("hash", h.Hex()).Msgf("Voted")
	return nil
}

// CreateProposalDataHash constructs and returns proposal data hash
func (p *Proposal) GetDataHash() common.Hash {
	return crypto.Keccak256Hash(append(p.HandlerAddress.Bytes(), p.Data...))
}

func idAndNonce(srcId uint8, nonce uint64) *big.Int {
	var data []byte
	data = append(data, big.NewInt(int64(nonce)).Bytes()...)
	data = append(data, uint8(srcId))
	return big.NewInt(0).SetBytes(data)
}
