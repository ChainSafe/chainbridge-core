package deposit

import (
	"github.com/ethereum/go-ethereum/common/math"
	"math/big"
)

func ConstructErc20DepositData(destRecipient []byte, amount *big.Int) []byte {
	var data []byte
	data = append(data, math.PaddedBigBytes(amount, 32)...)
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(destRecipient))), 32)...)
	data = append(data, destRecipient...)
	return data
}

func ConstructErc721DepositData(destRecipient []byte, tokenId *big.Int, metadata []byte) []byte {
	var data []byte
	data = append(data, math.PaddedBigBytes(tokenId, 32)...)                               // Token Id
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(destRecipient))), 32)...) // Length of recipient
	data = append(data, destRecipient...)                                                  // Recipient
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(metadata))), 32)...)      // Length of metadata
	data = append(data, metadata...)                                                       // Metadata
	return data
}

func ConstructGenericDepositData(metadata []byte) []byte {
	var data []byte
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(metadata))), 32)...)
	data = append(data, metadata...)
	return data
}
