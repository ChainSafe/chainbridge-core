package itx

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/forwarder"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/ethereum/go-ethereum/crypto"
	signer "github.com/ethereum/go-ethereum/signer/core"
	"github.com/rs/zerolog/log"
)

type ForwarderContract interface {
	GetNonce(from common.Address) (*big.Int, error)
	ExecuteData(forwardReq forwarder.ForwardRequest, domainSeparator *[32]byte, typeHash *[32]byte, suffixData []byte, sig []byte) ([]byte, error)
	ContractAddress() *common.Address
}

type NonceStorer interface {
	StoreNonce(chainID *big.Int, nonce *big.Int) error
	GetNonce(chainID *big.Int) (*big.Int, error)
}

type GsnForwarder struct {
	kp                *secp256k1.Keypair
	nonce             *big.Int
	nonceLock         sync.Mutex
	chainID           *big.Int
	forwarderContract ForwarderContract
	nonceStore        NonceStorer
}

// NewGsnForwarder creates an instance of GsnForwarder
func NewGsnForwarder(chainID *big.Int, kp *secp256k1.Keypair, forwarderContract ForwarderContract, nonceStore NonceStorer) *GsnForwarder {
	return &GsnForwarder{
		chainID:           chainID,
		kp:                kp,
		forwarderContract: forwarderContract,
		nonceStore:        nonceStore,
	}
}

// LockNonce locks mutex for nonce to prevent nonce duplication
func (c *GsnForwarder) LockNonce() {
	c.nonceLock.Lock()
}

// UnlockNonce unlocks mutext for nonce and stores nonce into storage.
//
// Nonce is stored on unlock, because current nonce should always be the correct one when unlocking.
func (c *GsnForwarder) UnlockNonce() {
	err := c.nonceStore.StoreNonce(c.chainID, c.nonce)
	if err != nil {
		log.Error().Err(fmt.Errorf("failed storing nonce: %v", err))
	}

	c.nonceLock.Unlock()
}

// UnsafeNonce returns current valid nonce for a forwarded transaction.
//
// If nonce is not set, looks for nonce in storage and on contract and returns the
// higher one. Nonce in storage can be higher if there are pending transactions after
// relayer has been manually shutdown.
func (c *GsnForwarder) UnsafeNonce() (*big.Int, error) {
	if c.nonce == nil {
		storedNonce, err := c.nonceStore.GetNonce(c.chainID)
		if err != nil {
			return nil, err
		}

		from := common.HexToAddress(c.kp.Address())
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
	return nonce, nil
}

// UnsafeIncreaseNonce increases nonce value by 1. Should be used
// while nonce is locked.
func (c *GsnForwarder) UnsafeIncreaseNonce() {
	c.nonce.Add(c.nonce, big.NewInt(1))
}

func (c *GsnForwarder) ForwarderAddress() common.Address {
	return *c.forwarderContract.ContractAddress()
}

func (c *GsnForwarder) ChainId() *big.Int {
	return c.chainID
}

// ForwarderData returns ABI packed and signed byte data for a forwarded transaction
func (c *GsnForwarder) ForwarderData(to common.Address, data []byte, opts transactor.TransactOptions) ([]byte, error) {
	from := c.kp.Address()
	forwarderHash, domainSeperator, typeHash, err := c.typedHash(
		from,
		to.String(),
		data,
		math.NewHexOrDecimal256(opts.Value.Int64()),
		math.NewHexOrDecimal256(int64(opts.GasLimit)),
		opts.Nonce,
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

	forwardReq := forwarder.ForwardRequest{
		From:       common.HexToAddress(from),
		To:         to,
		Value:      opts.Value,
		Gas:        big.NewInt(int64(opts.GasLimit)),
		Nonce:      opts.Nonce,
		Data:       data,
		ValidUntil: big.NewInt(0),
	}
	suffixData := common.Hex2Bytes("0x")
	return c.forwarderContract.ExecuteData(forwardReq, domainSeperator, typeHash, suffixData, sig)
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
