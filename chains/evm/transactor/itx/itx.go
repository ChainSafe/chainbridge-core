package itx

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/transactor"
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rpc"
)

type Forwarder interface {
	GetNonce(from common.Address) (*big.Int, error)
	GetForwarderAddress() common.Address
	GetChainId() uint8
	GetForwarderData(to common.Address, data []byte, kp *secp256k1.Keypair, opts transactor.TransactOptions) ([]byte, error)
}

type RelayTx struct {
	to   common.Address
	data []byte
	opts transactor.TransactOptions
}

type SignedRelayTx struct {
	*RelayTx
	txID common.Hash
	sig  []byte
}

type ITXTransactor struct {
	forwarder Forwarder
	rpcClient *rpc.Client
	kp        *secp256k1.Keypair
}

func NewITXTransactor(url string, forwarder Forwarder, kp *secp256k1.Keypair) (*ITXTransactor, error) {
	rpcClient, err := rpc.DialContext(context.TODO(), url)
	if err != nil {
		return nil, err
	}

	return &ITXTransactor{
		rpcClient: rpcClient,
		forwarder: forwarder,
		kp:        kp,
	}, nil
}

func (itx *ITXTransactor) Transact(to common.Address, data []byte, opts transactor.TransactOptions) (common.Hash, error) {
	nonce, err := itx.forwarder.GetNonce(itx.kp.CommonAddress())
	if err != nil {
		return common.Hash{}, err
	}
	opts.Nonce = nonce
	opts.ChainID = itx.forwarder.GetChainId()

	forwarderData, err := itx.forwarder.GetForwarderData(to, data, itx.kp, opts)
	if err != nil {
		return common.Hash{}, err
	}

	signedTx, err := itx.signRelayTx(&RelayTx{
		to:   itx.forwarder.GetForwarderAddress(),
		data: forwarderData,
		opts: opts,
	})
	if err != nil {
		return common.Hash{}, err
	}

	return itx.sendTransaction(context.Background(), signedTx)
}

func (itx *ITXTransactor) signRelayTx(tx *RelayTx) (*SignedRelayTx, error) {
	uint256Type, _ := abi.NewType("uint256", "uint256", nil)
	addressType, _ := abi.NewType("address", "address", nil)
	bytesType, _ := abi.NewType("bytes", "bytes", nil)
	stringType, _ := abi.NewType("string", "string", nil)
	arguments := abi.Arguments{
		{Type: addressType},
		{Type: bytesType},
		{Type: uint256Type},
		{Type: uint256Type},
		{Type: stringType},
	}

	packed, err := arguments.Pack(tx.to, tx.data, tx.opts.GasLimit, big.NewInt(int64(tx.opts.ChainID)), tx.opts.Priority)
	if err != nil {
		return nil, err
	}

	txID := crypto.Keccak256Hash(packed)
	msg := fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(txID), string(txID.Bytes()))
	hash := crypto.Keccak256Hash([]byte(msg))
	sig, err := crypto.Sign(hash.Bytes(), itx.kp.PrivateKey())
	if err != nil {
		return nil, nil
	}

	return &SignedRelayTx{
		RelayTx: tx,
		sig:     sig,
		txID:    txID,
	}, nil
}

func (itx *ITXTransactor) sendTransaction(ctx context.Context, signedTx *SignedRelayTx) (common.Hash, error) {
	type txResponse struct {
		RelayTransactionHash hexutil.Bytes
	}

	sig := "0x" + common.Bytes2Hex(signedTx.sig)
	txArg := map[string]interface{}{
		"to":       signedTx.to.String(),
		"data":     "0x" + common.Bytes2Hex(signedTx.data),
		"gas":      signedTx.opts.GasLimit.String(),
		"schedule": signedTx.opts.Priority,
	}

	var resp txResponse
	err := itx.rpcClient.CallContext(ctx, &resp, "relay_sendTransaction", txArg, sig)
	if err != nil {
		return common.Hash{}, err
	}

	return common.HexToHash(resp.RelayTransactionHash.String()), nil
}
