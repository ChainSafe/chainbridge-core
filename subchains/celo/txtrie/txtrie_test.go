// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package txtrie

import (
	"fmt"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/trie"
)

var (
	emptyHash = common.HexToHash("")
)

func computeEthReferenceTrieHash(transactions types.Transactions) (common.Hash, error) {
	newTrie, err := trie.New(emptyRoot, trie.NewDatabase(nil))
	if err != nil {
		return emptyHash, err
	}

	for i, tx := range transactions {

		key, err := rlp.EncodeToBytes(uint(i))
		if err != nil {
			return emptyHash, err
		}

		value, err := rlp.EncodeToBytes(tx)
		if err != nil {
			return emptyHash, err
		}

		err = newTrie.TryUpdate(key, value)
		if err != nil {
			return emptyHash, err
		}
	}

	return newTrie.Hash(), nil
}

func TestAddEmptyTrie(t *testing.T) {
	emptyTransactions := make([]*types.Transaction, 0)
	_, err := CreateNewTrie(emptyRoot, types.Transactions(emptyTransactions))
	if err != nil {
		t.Fatal(err)
	}
}

func TestAddSingleTrieUpdate(t *testing.T) {
	vals := GetTransactions1()
	root, err := computeEthReferenceTrieHash(vals)
	if err != nil {
		t.Fatal(err)
	}

	tr, err := CreateNewTrie(root, types.Transactions(vals))
	if err != nil {
		t.Fatal(err)
	}
	keyRlp, err := rlp.EncodeToBytes(uint(0))
	if err != nil {
		t.Fatal(err.Error())
	}
	proof, key, err := RetrieveProof(tr, keyRlp)
	if err != nil {
		t.Fatal(err.Error())
	}

	if proof == nil {
		t.Fatal("proof is nil")
	}
	if common.Bytes2Hex(key) != "0800" {
		t.Fatal(fmt.Sprintf("wrong RLP key is %s", common.Bytes2Hex(key)))
	}

}
