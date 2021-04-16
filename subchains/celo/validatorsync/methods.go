//Copyright 2020 ChainSafe Systems
//SPDX-License-Identifier: LGPL-3.0-only
package validatorsync

import (
	"math/big"

	"github.com/celo-org/celo-blockchain/consensus/istanbul"
	"github.com/celo-org/celo-blockchain/core/types"
	blscrypto "github.com/celo-org/celo-blockchain/crypto/bls"
	"github.com/celo-org/celo-bls-go/bls"
	"github.com/pkg/errors"
)

var (
	ErrorWrongInitialValidators = errors.New("wrong initial validators")
)

func applyValidatorsDiff(extra *types.IstanbulExtra, validators []*istanbul.ValidatorData) ([]*istanbul.ValidatorData, error) {
	var addedValidators []*istanbul.ValidatorData
	for i, addr := range extra.AddedValidators {
		addedValidators = append(addedValidators, &istanbul.ValidatorData{Address: addr, BLSPublicKey: extra.AddedValidatorsPublicKeys[i]})
	}

	if extra.RemovedValidators.BitLen() > len(validators) {
		return nil, ErrorWrongInitialValidators
	}
	var newValidators []*istanbul.ValidatorData
	if extra.RemovedValidators.BitLen() == 0 {
		newValidators = append(make([]*istanbul.ValidatorData, 0), validators...)
	} else {
		for i := 0; i < extra.RemovedValidators.BitLen(); i++ {
			if extra.RemovedValidators.Bit(i) == 1 {
				continue
			}
			newValidators = append(newValidators, validators[i])
		}
	}
	newValidators = append(newValidators, addedValidators...)
	return newValidators, nil
}

func aggregatePublicKeys(validators []*istanbul.ValidatorData) (*bls.PublicKey, error) {
	publicKeys := make([]blscrypto.SerializedPublicKey, len(validators))
	for i := range validators {
		publicKeys[i] = validators[i].BLSPublicKey
	}

	publicKeyObjs := make([]*bls.PublicKey, len(publicKeys))
	for i := range publicKeys {
		publicKeyObj, err := bls.DeserializePublicKeyCached(publicKeys[i][:])
		if err != nil {
			return nil, err
		}
		publicKeyObjs[i] = publicKeyObj
		publicKeyObj.Destroy()
	}
	apk, err := bls.AggregatePublicKeys(publicKeyObjs)
	if err != nil {
		return nil, err
	}
	defer apk.Destroy()

	return apk, nil
}

func computeLastBlockOfEpochForProvidedBlock(block *big.Int, epochSize uint64) *big.Int {
	epochNumber := istanbul.GetEpochNumber(block.Uint64(), epochSize)
	lastBlock := istanbul.GetEpochLastBlockNumber(epochNumber, epochSize)
	return big.NewInt(0).SetUint64(lastBlock)
}
