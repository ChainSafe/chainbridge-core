package forwarder

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ChainSafe/chainbridge-core/chains/evm/transactor"
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	signer "github.com/ethereum/go-ethereum/signer/core"
)

// A forwarder client for sending transactions via a global GSN forwarder
type GsnForwarder struct {
	nonce             *big.Int
	nonceLock         sync.Mutex
	chainID           uint8
	forwarderContract Forwarder
}

type Forwarder interface {
	GetNonce(from common.Address) (*big.Int, error)
	GetAddress() common.Address
	GetABI() *abi.ABI
}

type ForwardRequest struct {
	From       common.Address
	To         common.Address
	Value      *big.Int
	Gas        *big.Int
	Nonce      *big.Int
	Data       []byte
	ValidUntil *big.Int
}

func NewGsnForwarder(chainID uint8, forwarderContract Forwarder) *GsnForwarder {
	return &GsnForwarder{
		chainID:           chainID,
		forwarderContract: forwarderContract,
	}
}

func (c *GsnForwarder) GetNonce(from common.Address) (*big.Int, error) {
	c.nonceLock.Lock()
	defer c.nonceLock.Unlock()

	if c.nonce == nil {
		nonce, err := c.forwarderContract.GetNonce(from)
		if err != nil {
			return nil, err
		}

		c.nonce = nonce
	}

	nonce := big.NewInt(c.nonce.Int64())
	c.nonce.Add(c.nonce, big.NewInt(1))
	return nonce, nil
}

func (c *GsnForwarder) GetForwarderAddress() common.Address {
	return c.forwarderContract.GetAddress()
}

func (c *GsnForwarder) GetChainId() uint8 {
	return c.chainID
}

func (c *GsnForwarder) GetForwarderData(to common.Address, data []byte, kp *secp256k1.Keypair, opts transactor.TransactOptions) ([]byte, error) {
	from := kp.Address()
	forwarderHash, domainSeperator, typeHash, err := c.typedHash(
		from,
		to.String(),
		data,
		math.NewHexOrDecimal256(opts.Value.Int64()),
		math.NewHexOrDecimal256(opts.GasLimit.Int64()),
		opts.Nonce,
		c.GetForwarderAddress().Hex(),
	)
	if err != nil {
		return nil, err
	}

	sig, err := crypto.Sign(forwarderHash, kp.PrivateKey())
	if err != nil {
		return nil, err
	}
	sig[64] += 27 // Transform V from 0/1 to 27/28

	var suffixData = common.Hex2Bytes("0x")
	forwardReq := ForwardRequest{
		From:       common.HexToAddress(from),
		To:         to,
		Value:      opts.Value,
		Gas:        opts.GasLimit,
		Nonce:      opts.Nonce,
		Data:       data,
		ValidUntil: big.NewInt(0),
	}
	return c.forwarderContract.GetABI().Pack("execute", forwardReq, domainSeperator, typeHash, suffixData, sig)
}

func (c *GsnForwarder) typedHash(
	from, to string,
	data []byte,
	value, gas *math.HexOrDecimal256,
	nonce *big.Int,
	verifyingContract string,
) ([]byte, *[32]byte, *[32]byte, error) {
	chainId := math.NewHexOrDecimal256(int64(c.chainID))
	typedData := signer.TypedData{
		Types: signer.Types{
			"EIP712Domain": []signer.Type{
				{Name: "name", Type: "string"},
				{Name: "version", Type: "string"},
				{Name: "chainId", Type: "uint256"},
				{Name: "verifyingContract", Type: "address"},
			},
			"ForwardRequest": []signer.Type{
				{Name: "from", Type: "address"},
				{Name: "to", Type: "address"},
				{Name: "value", Type: "uint256"},
				{Name: "gas", Type: "uint256"},
				{Name: "nonce", Type: "uint256"},
				{Name: "data", Type: "bytes"},
				{Name: "validUntil", Type: "uint256"},
			},
		},
		PrimaryType: "ForwardRequest",
		Domain: signer.TypedDataDomain{
			Name:              "GSN Relayed Transaction",
			ChainId:           chainId,
			Version:           "2",
			VerifyingContract: verifyingContract,
		},
		Message: signer.TypedDataMessage{
			"from":       from,
			"to":         to,
			"value":      value,
			"gas":        gas,
			"data":       data,
			"nonce":      math.NewHexOrDecimal256(nonce.Int64()),
			"validUntil": math.NewHexOrDecimal256(0),
		},
	}
	domainSeparator, err := typedData.HashStruct("EIP712Domain", typedData.Domain.Map())
	if err != nil {
		return nil, nil, nil, err
	}

	var fixedSizeDomainSeperator [32]byte
	copy(fixedSizeDomainSeperator[:], domainSeparator)

	typedDataHash, err := typedData.HashStruct(typedData.PrimaryType, typedData.Message)
	if err != nil {
		return nil, nil, nil, err
	}
	var fixedSizeTypeHash [32]byte
	copy(fixedSizeTypeHash[:], typedData.TypeHash(typedData.PrimaryType))

	rawData := []byte(fmt.Sprintf("\x19\x01%s%s", string(domainSeparator), string(typedDataHash)))
	return crypto.Keccak256(rawData), &fixedSizeDomainSeperator, &fixedSizeTypeHash, nil
}
