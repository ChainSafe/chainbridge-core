package deposit

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common/math"
)

func constructMainDepositData(tokenStats *big.Int, destRecipient []byte) []byte {
	var data []byte
	data = append(data, math.PaddedBigBytes(tokenStats, 32)...)                            // Amount (ERC20) or Token Id (ERC721)
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(destRecipient))), 32)...) // length of recipient
	data = append(data, destRecipient...)                                                  // Recipient
	return data
}

func ConstructErc20DepositData(destRecipient []byte, amount *big.Int) []byte {
	data := constructMainDepositData(amount, destRecipient)
	return data
}

func ConstructErc20DepositDataWithPriority(destRecipient []byte, amount *big.Int, priority uint8) []byte {
	data := constructMainDepositData(amount, destRecipient)
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len([]uint8{priority}))), 1)...) // Length of priority
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(priority)), 1)...)               // Priority
	return data
}

func ConstructErc721DepositData(destRecipient []byte, tokenId *big.Int, metadata []byte) []byte {
	data := constructMainDepositData(tokenId, destRecipient)
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(metadata))), 32)...) // Length of metadata
	data = append(data, metadata...)                                                  // Metadata
	return data
}

func ConstructErc721DepositDataWithPriority(destRecipient []byte, tokenId *big.Int, metadata []byte, priority uint8) []byte {
	data := constructMainDepositData(tokenId, destRecipient)
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(metadata))), 32)...)         // Length of metadata
	data = append(data, metadata...)                                                          // Metadata
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len([]uint8{priority}))), 1)...) // Length of priority
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(priority)), 1)...)               // Priority
	return data
}

func ConstructGenericDepositData(metadata []byte) []byte {
	var data []byte
	data = append(data, math.PaddedBigBytes(big.NewInt(int64(len(metadata))), 32)...) // Length of metadata
	data = append(data, metadata...)                                                  // Metadata
	return data
}
