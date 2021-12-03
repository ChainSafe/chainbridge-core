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

type Forwarder interface {
	GetNonce(from common.Address) (*big.Int, error)
	Address() common.Address
	ABI() *abi.ABI
}

type NonceStorer interface {
	StoreNonce(chainID *big.Int, nonce *big.Int) error
	GetNonce(chainID *big.Int) (*big.Int, error)
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

type GsnForwarder struct {
	kp                *secp256k1.Keypair
	nonce             *big.Int
	nonceLock         sync.Mutex
	chainID           *big.Int
	forwarderContract Forwarder
	nonceStore        NonceStorer
}

func NewGsnForwarder(chainID *big.Int, kp *secp256k1.Keypair, forwarderContract Forwarder, nonceStore NonceStorer) *GsnForwarder {
	return &GsnForwarder{
		chainID:           chainID,
		kp:                kp,
		forwarderContract: forwarderContract,
		nonceStore:        nonceStore,
	}
}

func (c *GsnForwarder) NextNonce(from common.Address) (*big.Int, error) {
	c.nonceLock.Lock()
	defer c.nonceLock.Unlock()

	if c.nonce == nil {
		storedNonce, err := c.nonceStore.GetNonce(c.chainID)
		if err != nil {
			return nil, err
		}

		contractNonce, err := c.forwarderContract.GetNonce(from)
		if err != nil {
			return nil, err
		}

		var nonce *big.Int
		if storedNonce.Cmp(contractNonce) >= 0 {
			nonce = storedNonce
		} else {
			nonce = contractNonce
		}

		c.nonce = nonce
	}

	nonce := big.NewInt(c.nonce.Int64())
	c.nonce.Add(c.nonce, big.NewInt(1))

	err := c.nonceStore.StoreNonce(c.chainID, nonce)
	if err != nil {
		return nonce, err
	}

	return nonce, nil
}

func (c *GsnForwarder) ForwarderAddress() common.Address {
	return c.forwarderContract.Address()
}

func (c *GsnForwarder) ChainId() *big.Int {
	return c.chainID
}

func (c *GsnForwarder) ForwarderData(to common.Address, data []byte, opts transactor.TransactOptions) ([]byte, error) {
	from := c.kp.Address()
	nonce, err := c.NextNonce(common.HexToAddress(from))
	if err != nil {
		return nil, err
	}

	forwarderHash, domainSeperator, typeHash, err := c.typedHash(
		from,
		to.String(),
		data,
		math.NewHexOrDecimal256(opts.Value.Int64()),
		math.NewHexOrDecimal256(opts.GasLimit.Int64()),
		nonce,
		c.ForwarderAddress().Hex(),
	)
	if err != nil {
		return nil, err
	}

	sig, err := crypto.Sign(forwarderHash, c.kp.PrivateKey())
	if err != nil {
		return nil, err
	}
	sig[64] += 27 // Transform V from 0/1 to 27/28

	forwardReq := ForwardRequest{
		From:       common.HexToAddress(from),
		To:         to,
		Value:      opts.Value,
		Gas:        opts.GasLimit,
		Nonce:      nonce,
		Data:       data,
		ValidUntil: big.NewInt(0),
	}
	suffixData := common.Hex2Bytes("0x")
	return c.forwarderContract.ABI().Pack("execute", forwardReq, domainSeperator, typeHash, suffixData, sig)
}

func (c *GsnForwarder) typedHash(
	from, to string,
	data []byte,
	value, gas *math.HexOrDecimal256,
	nonce *big.Int,
	verifyingContract string,
) ([]byte, *[32]byte, *[32]byte, error) {
	chainId := math.NewHexOrDecimal256(c.chainID.Int64())
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
