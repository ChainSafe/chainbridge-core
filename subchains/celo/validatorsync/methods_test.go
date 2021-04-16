//Copyright 2020 ChainSafe Systems
//SPDX-License-Identifier: LGPL-3.0-only
package validatorsync

import (
	"math/big"
	"testing"

	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/consensus/istanbul"
	"github.com/celo-org/celo-blockchain/core/types"
	"github.com/celo-org/celo-blockchain/crypto"
	blscrypto "github.com/celo-org/celo-blockchain/crypto/bls"
	"github.com/stretchr/testify/suite"
)

type WriterTestSuite struct {
	suite.Suite
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(WriterTestSuite))
}

func (s *WriterTestSuite) SetupSuite()    {}
func (s *WriterTestSuite) TearDownSuite() {}
func (s *WriterTestSuite) SetupTest()     {}
func (s *WriterTestSuite) TearDownTest()  {}

func (s *WriterTestSuite) TestApplyValidatorsDiffOnlyAdd() {
	startVals := make([]*istanbul.ValidatorData, 2)
	startVals[0] = &istanbul.ValidatorData{Address: common.Address{0x0f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	startVals[1] = &istanbul.ValidatorData{Address: common.Address{0x1f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	addedAddresses := []common.Address{{0x3f}}
	extra := &types.IstanbulExtra{
		AddedValidators:           addedAddresses,
		RemovedValidators:         big.NewInt(0),
		AddedValidatorsPublicKeys: []blscrypto.SerializedPublicKey{{0x3f}},
	}
	resVals, err := applyValidatorsDiff(extra, startVals)
	s.Nil(err)
	s.Equal(3, len(resVals))
	s.Equal(resVals[2].BLSPublicKey, blscrypto.SerializedPublicKey{0x3f})
}

func (s *WriterTestSuite) TestApplyValidatorsDiff() {
	startVals := make([]*istanbul.ValidatorData, 3)
	startVals[0] = &istanbul.ValidatorData{Address: common.Address{0x0f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	startVals[1] = &istanbul.ValidatorData{Address: common.Address{0x1f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	startVals[2] = &istanbul.ValidatorData{Address: common.Address{0x2f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	addedAddresses := []common.Address{{0x3f}}

	extra := &types.IstanbulExtra{
		AddedValidators:           addedAddresses,
		RemovedValidators:         big.NewInt(4),
		AddedValidatorsPublicKeys: []blscrypto.SerializedPublicKey{{0x3f}},
	}
	resVals, err := applyValidatorsDiff(extra, startVals)
	s.Nil(err)
	s.Equal(len(resVals), 3)
	s.Equal(resVals[2].BLSPublicKey, blscrypto.SerializedPublicKey{0x3f})
}

func (s *WriterTestSuite) TestApplyValidatorsDiffEmptyStartVals() {
	startVals := make([]*istanbul.ValidatorData, 0)
	addedAddresses := []common.Address{{0x3f}}

	extra := &types.IstanbulExtra{
		AddedValidators:           addedAddresses,
		RemovedValidators:         big.NewInt(0),
		AddedValidatorsPublicKeys: []blscrypto.SerializedPublicKey{{0x3f}},
	}
	resVals, err := applyValidatorsDiff(extra, startVals)
	s.Nil(err)
	s.Equal(len(resVals), 1)
	s.Equal(resVals[0].BLSPublicKey, blscrypto.SerializedPublicKey{0x3f})
}

func (s *WriterTestSuite) TestApplyValidatorsDiffWithRemovedOnEmptyVals() {
	startVals := make([]*istanbul.ValidatorData, 0)
	addedAddresses := []common.Address{{0x3f}}

	extra := &types.IstanbulExtra{
		AddedValidators:           addedAddresses,
		RemovedValidators:         big.NewInt(1),
		AddedValidatorsPublicKeys: []blscrypto.SerializedPublicKey{{0x3f}},
	}
	resVals, err := applyValidatorsDiff(extra, startVals)
	s.Nil(resVals)
	s.NotNil(err)
	s.Equal(err, ErrorWrongInitialValidators)
}

func (s *WriterTestSuite) TestDefineBlocksEpochLastBlockNumber() {
	s.Equal(computeLastBlockOfEpochForProvidedBlock(big.NewInt(0), 2335), big.NewInt(0))

	s.Equal(computeLastBlockOfEpochForProvidedBlock(big.NewInt(11), 12), big.NewInt(12))

	s.Equal(computeLastBlockOfEpochForProvidedBlock(big.NewInt(251), 12), big.NewInt(252))

}

func (s *WriterTestSuite) TestAggregatePublicKeys() {
	startVals := make([]*istanbul.ValidatorData, 3)
	testKey1, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	testKey2, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f292")
	testKey3, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f293")
	blsPK1, _ := blscrypto.ECDSAToBLS(testKey1)
	blsPK2, _ := blscrypto.ECDSAToBLS(testKey2)
	blsPK3, _ := blscrypto.ECDSAToBLS(testKey3)
	pubKey1, _ := blscrypto.PrivateToPublic(blsPK1)
	pubKey2, _ := blscrypto.PrivateToPublic(blsPK2)
	pubKey3, _ := blscrypto.PrivateToPublic(blsPK3)

	startVals[0] = &istanbul.ValidatorData{Address: common.Address{0x0f}, BLSPublicKey: pubKey1}
	startVals[1] = &istanbul.ValidatorData{Address: common.Address{0x1f}, BLSPublicKey: pubKey2}
	startVals[2] = &istanbul.ValidatorData{Address: common.Address{0x2f}, BLSPublicKey: pubKey3}
	apk, err := aggregatePublicKeys(startVals)
	s.Nil(err)
	s.NotNil(apk)
	// checking that function is clear
	apk2, err := aggregatePublicKeys(startVals)
	s.Nil(err)
	s.Equal(apk, apk2)
}
