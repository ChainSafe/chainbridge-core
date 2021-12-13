package contracts

import (
	"context"
	"fmt"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"
)

type Contract struct {
	contractAddress common.Address
	ABI             abi.ABI
	bytecode        []byte
	client          calls.ContractCallerDispatcher
	transactor.Transactor
}

func NewContract(
	contractAddress common.Address,
	abi abi.ABI,
	bytecode []byte,
	client calls.ContractCallerDispatcher,
	transactor transactor.Transactor,
) Contract {
	return Contract{
		contractAddress: contractAddress,
		ABI:             abi,
		bytecode:        bytecode,
		client:          client,
		Transactor:      transactor,
	}
}

func (c *Contract) ContractAddress() *common.Address {
	return &c.contractAddress
}

func (c *Contract) PackMethod(method string, args ...interface{}) ([]byte, error) {
	input, err := c.ABI.Pack(method, args...)
	if err != nil {
		log.Error().Err(fmt.Errorf("pack method error: %v", err))
		return []byte{}, err
	}
	return input, nil
}

func (c *Contract) UnpackResult(method string, output []byte) ([]interface{}, error) {
	res, err := c.ABI.Unpack(method, output)
	if err != nil {
		log.Error().Err(fmt.Errorf("unpack output error: %v", err))
		return nil, err
	}
	return res, err
}

func (c *Contract) ExecuteTransaction(method string, opts transactor.TransactOptions, args ...interface{}) (*common.Hash, error) {
	input, err := c.PackMethod(method, args...)
	if err != nil {
		return nil, err
	}
	h, err := c.Transact(&c.contractAddress, input, opts)
	if err != nil {
		log.Error().
			Str("contract", c.contractAddress.String()).
			Err(err).
			Msgf("error on executing %s", method)
		return nil, err
	}
	log.Debug().
		Str("txHash", h.String()).
		Str("contract", c.contractAddress.String()).
		Msgf("method %s executed", method)
	return h, err
}

func (c *Contract) CallContract(method string, args ...interface{}) ([]interface{}, error) {
	input, err := c.PackMethod(method, args...)
	if err != nil {
		return nil, err
	}
	msg := ethereum.CallMsg{From: c.client.From(), To: &c.contractAddress, Data: input}
	out, err := c.client.CallContract(context.TODO(), calls.ToCallArg(msg), nil)
	if err != nil {
		log.Error().
			Str("contract", c.contractAddress.String()).
			Err(err).
			Msgf("error on calling %s", method)
		return nil, err
	}
	if len(out) == 0 {
		// Make sure we have a contract to operate on, and bail out otherwise.
		if code, err := c.client.CodeAt(context.Background(), c.contractAddress, nil); err != nil {
			return nil, err
		} else if len(code) == 0 {
			return nil, fmt.Errorf("no code at provided address %s", c.contractAddress.String())
		}
	}
	log.Debug().
		Str("contract", c.contractAddress.String()).
		Msgf("method %s called", method)
	return c.UnpackResult(method, out)
}

func (c *Contract) DeployContract(params ...interface{}) (common.Address, error) {
	input, err := c.PackMethod("", params...)
	if err != nil {
		return common.Address{}, err
	}
	opts := transactor.TransactOptions{GasLimit: consts.DefaultDeployGasLimit}
	hash, err := c.Transact(nil, append(c.bytecode, input...), opts)
	if err != nil {
		return common.Address{}, err
	}
	tx, _, err := c.client.GetTransactionByHash(*hash)
	if err != nil {
		return common.Address{}, err
	}
	address := crypto.CreateAddress(c.client.From(), tx.Nonce())
	c.contractAddress = address
	log.Debug().
		Str("txHash", hash.String()).
		Str("deployedAddress", address.String()).
		Msgf("successful contract deployment")
	return address, nil
}
