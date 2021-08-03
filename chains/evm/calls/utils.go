package calls

import (
	"context"
	"strings"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
)

// UTILS

func sliceTo32Bytes(in []byte) [32]byte {
	var res [32]byte
	copy(res[:], in)
	return res
}

func toCallArg(msg ethereum.CallMsg) map[string]interface{} {
	arg := map[string]interface{}{
		"from": msg.From,
		"to":   msg.To,
	}
	if len(msg.Data) > 0 {
		arg["data"] = hexutil.Bytes(msg.Data)
	}
	if msg.Value != nil {
		arg["value"] = (*hexutil.Big)(msg.Value)
	}
	if msg.Gas != 0 {
		arg["gas"] = hexutil.Uint64(msg.Gas)
	}
	if msg.GasPrice != nil {
		arg["gasPrice"] = (*hexutil.Big)(msg.GasPrice)
	}
	return arg
}

func mintRole(chainClient ChainClient, erc20Contract common.Address) ([32]byte, error) {
	a, err := abi.JSON(strings.NewReader(ERC20PresetMinterPauserABI))
	if err != nil {
		return [32]byte{}, err
	}
	input, err := a.Pack("MINTER_ROLE")
	if err != nil {
		return [32]byte{}, err
	}
	msg := ethereum.CallMsg{From: common.Address{}, To: &erc20Contract, Data: input}
	out, err := chainClient.CallContract(context.TODO(), toCallArg(msg), nil)
	if err != nil {
		return [32]byte{}, err
	}
	res, err := a.Unpack("MINTER_ROLE", out)
	if err != nil {
		return [32]byte{}, err
	}
	out0 := *abi.ConvertType(res[0], new([32]byte)).(*[32]byte)
	return out0, nil
}
