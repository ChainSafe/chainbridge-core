// Copyright 2020 ChainSafe Systems
// SPDX-License-Identifier: LGPL-3.0-only

package txtrie

import (
	"bytes"
	"errors"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	"github.com/ethereum/go-ethereum/ethdb/memorydb"
	"github.com/ethereum/go-ethereum/rlp"
	ethtrie "github.com/ethereum/go-ethereum/trie"
)

var (
	// from https://github.com/ethereum/go-ethereum/blob/bcb308745010675671991522ad2a9e811938d7fb/trie/trie.go#L32
	emptyRoot = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
)

func CreateNewTrie(root common.Hash, transactions types.Transactions) (*ethtrie.Trie, error) {
	if transactions == nil {
		return nil, errors.New("transactions cannot be nil")
	}
	db := memorydb.New()
	trie, err := ethtrie.New(emptyRoot, ethtrie.NewDatabase(db))
	if err != nil {
		return nil, err
	}
	for i, tx := range transactions {
		key, err := rlp.EncodeToBytes(uint(i))
		if err != nil {
			return nil, err
		}
		value, err := rlp.EncodeToBytes(tx)
		if err != nil {
			return nil, err
		}
		trie.Update(key, value)
	}
	if trie.Hash().Hex() != root.Hex() {
		return nil, errors.New("transaction roots don't match")
	}
	return trie, nil
}

func RetrieveProof(trie *ethtrie.Trie, key []byte) ([]byte, []byte, error) {
	nodeIterator := trie.NodeIterator(key)
	trieIterator := ethtrie.NewIterator(nodeIterator)
	proof := make([][][]byte, 0)
	trieIterator.Next()
	value := trieIterator.Prove()
	for _, v := range value {
		n := make([][]byte, 0, 17)
		err := rlp.DecodeBytes(v, &n)
		if err != nil {
			return nil, nil, err
		}
		proof = append(proof, n)
	}
	buf := &bytes.Buffer{}
	err := rlp.Encode(buf, proof)
	if err != nil {
		return nil, nil, err
	}
	leafKey := keybytesToHex(nodeIterator.LeafKey())
	leafKey = leafKey[:len(leafKey)-1]
	return buf.Bytes(), leafKey, nil
}

func keybytesToHex(str []byte) []byte {
	l := len(str)*2 + 1
	var nibbles = make([]byte, l)
	for i, b := range str {
		nibbles[i*2] = b / 16
		nibbles[i*2+1] = b % 16
	}
	nibbles[l-1] = 16
	return nibbles
}

// VerifyProof verifies merkle proof on path key against the provided root
func VerifyProof(root common.Hash, key []byte, proof ethdb.KeyValueStore) (bool, error) {
	exists, _, err := ethtrie.VerifyProof(root, key, proof)

	if err != nil {
		return false, err
	}

	return exists != nil, nil
}
