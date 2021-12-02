package itx

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ChainSafe/chainbridge-core/chains/evm/transactor"
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
)

var (
	DefaultITXOpts = transactor.TransactOptions{
		GasLimit: big.NewInt(consts.DefaultGasLimit),
		GasPrice: big.NewInt(1),
		Priority: "slow",
		Value:    big.NewInt(0),
	}
)

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

type RelayCaller interface {
	CallContext(ctx context.Context, result interface{}, method string, args ...interface{}) error
}

type Forwarder interface {
	ForwarderAddress() common.Address
	ChainId() *big.Int
	ForwarderData(to common.Address, data []byte, opts transactor.TransactOptions) ([]byte, error)
}

type ITXTransactor struct {
	forwarder   Forwarder
	relayCaller RelayCaller
	kp          *secp256k1.Keypair
}

func NewITXTransactor(relayCaller RelayCaller, forwarder Forwarder, kp *secp256k1.Keypair) *ITXTransactor {
	return &ITXTransactor{
		relayCaller: relayCaller,
		forwarder:   forwarder,
		kp:          kp,
	}
}

func (itx *ITXTransactor) Transact(to common.Address, data []byte, opts transactor.TransactOptions) (common.Hash, error) {
	err := transactor.MergeTransactionOptions(&opts, &DefaultITXOpts)
	if err != nil {
		return common.Hash{}, err
	}
	opts.ChainID = itx.forwarder.ChainId()

	forwarderData, err := itx.forwarder.ForwarderData(to, data, opts)
	if err != nil {
		return common.Hash{}, err
	}

	// increase gas limit because of forwarder overhead
	opts.GasLimit = opts.GasLimit.Mul(
		opts.GasLimit, big.NewInt(2),
	)
	signedTx, err := itx.signRelayTx(&RelayTx{
		to:   itx.forwarder.ForwarderAddress(),
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
	packed, err := arguments.Pack(
		tx.to,
		tx.data,
		tx.opts.GasLimit,
		tx.opts.ChainID,
		tx.opts.Priority,
	)
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
	sig := "0x" + common.Bytes2Hex(signedTx.sig)
	txArg := map[string]interface{}{
		"to":       signedTx.to.String(),
		"data":     "0x" + common.Bytes2Hex(signedTx.data),
		"gas":      signedTx.opts.GasLimit.String(),
		"schedule": signedTx.opts.Priority,
	}

	resp := struct {
		RelayTransactionHash hexutil.Bytes
	}{}
	err := itx.relayCaller.CallContext(ctx, &resp, "relay_sendTransaction", txArg, sig)
	if err != nil {
		return common.Hash{}, err
	}

	return common.HexToHash(resp.RelayTransactionHash.String()), nil
}
