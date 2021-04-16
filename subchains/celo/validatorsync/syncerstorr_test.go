//Copyright 2020 ChainSafe Systems
//SPDX-License-Identifier: LGPL-3.0-only
package validatorsync

import (
	"errors"
	"math/big"
	"os"
	"testing"

	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/consensus/istanbul"
	"github.com/celo-org/celo-blockchain/crypto"
	blscrypto "github.com/celo-org/celo-blockchain/crypto/bls"
	"github.com/stretchr/testify/suite"
	"github.com/syndtr/goleveldb/leveldb"
)

type SyncerDBTestSuite struct {
	suite.Suite
	syncer *ValidatorsStore
}

func TestRunSyncerDBTestSuite(t *testing.T) {
	suite.Run(t, new(SyncerDBTestSuite))
}
func (s *SyncerDBTestSuite) SetupSuite()    {}
func (s *SyncerDBTestSuite) TearDownSuite() {}
func (s *SyncerDBTestSuite) SetupTest() {
	db, err := leveldb.OpenFile("./test/db", nil)
	if err != nil {
		s.Fail(err.Error())
	}
	syncer := NewValidatorsStore(db)
	s.syncer = syncer
}
func (s *SyncerDBTestSuite) TearDownTest() {
	s.syncer.Close()
	os.RemoveAll("./test")
}

func (s *SyncerDBTestSuite) TestSetValidatorsForBlock() {
	chainID := uint8(1)
	startVals := make([]*istanbul.ValidatorData, 3)
	startVals[0] = &istanbul.ValidatorData{Address: common.Address{0x0f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	startVals[1] = &istanbul.ValidatorData{Address: common.Address{0x1f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	startVals[2] = &istanbul.ValidatorData{Address: common.Address{0x2f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	err := s.syncer.SetValidatorsForBlock(big.NewInt(420), startVals, chainID)
	s.Nil(err)
	v, err := s.syncer.GetValidatorsForBlock(big.NewInt(420), chainID)
	s.Nil(err)
	s.Equal(3, len(v))
	s.Equal(common.Address{0x0f}, v[0].Address)
	b, err := s.syncer.GetLatestKnownEpochLastBlock(chainID)
	s.Nil(err)
	s.Equal(0, b.Cmp(big.NewInt(420)))

	validators, err := s.syncer.GetValidatorsForBlock(big.NewInt(420), chainID)
	s.Nil(err)
	s.Equal(3, len(validators))
	s.Equal(common.Address{0x0f}, validators[0].Address)

}

func (s *SyncerDBTestSuite) TestGetLatestKnownBlockWithEmptyDB() {
	chainID := uint8(1)
	v, err := s.syncer.GetLatestKnownEpochLastBlock(chainID)
	s.Nil(err)
	s.Equal(0, v.Cmp(big.NewInt(0)))
}

func (s *SyncerDBTestSuite) TestTestSetValidatorsForBlockForDifferentChains() {
	chainID1 := uint8(1)
	startVals1 := make([]*istanbul.ValidatorData, 3)
	startVals1[0] = &istanbul.ValidatorData{Address: common.Address{0x0f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	startVals1[1] = &istanbul.ValidatorData{Address: common.Address{0x1f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	startVals1[2] = &istanbul.ValidatorData{Address: common.Address{0x2f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	err := s.syncer.SetValidatorsForBlock(big.NewInt(420), startVals1, chainID1)
	s.Nil(err)

	chainID2 := uint8(2)
	startVals2 := make([]*istanbul.ValidatorData, 2)
	startVals2[0] = &istanbul.ValidatorData{Address: common.Address{0x3f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	startVals2[1] = &istanbul.ValidatorData{Address: common.Address{0x4f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	err = s.syncer.SetValidatorsForBlock(big.NewInt(420), startVals2, chainID2)
	s.Nil(err)

	v, err := s.syncer.GetValidatorsForBlock(big.NewInt(420), chainID1)
	s.Nil(err)
	s.Equal(3, len(v))
	s.Equal(common.Address{0x0f}, v[0].Address)
	b, err := s.syncer.GetLatestKnownEpochLastBlock(chainID1)
	s.Nil(err)
	s.Equal(0, b.Cmp(big.NewInt(420)))

	validators, err := s.syncer.GetValidatorsForBlock(big.NewInt(420), chainID1)
	s.Nil(err)
	s.Equal(3, len(validators))
	s.Equal(common.Address{0x0f}, validators[0].Address)

	v, err = s.syncer.GetValidatorsForBlock(big.NewInt(420), chainID2)
	s.Nil(err)
	s.Equal(2, len(v))
	s.Equal(common.Address{0x3f}, v[0].Address)
	b, err = s.syncer.GetLatestKnownEpochLastBlock(chainID2)
	s.Nil(err)
	s.Equal(0, b.Cmp(big.NewInt(420)))

	validators, err = s.syncer.GetValidatorsForBlock(big.NewInt(420), chainID2)
	s.Nil(err)
	s.Equal(2, len(validators))
	s.Equal(common.Address{0x3f}, validators[0].Address)
}

func (s *SyncerDBTestSuite) TestGetAPKForBlock() {
	chainID := uint8(1)
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

	err := s.syncer.SetValidatorsForBlock(big.NewInt(12), startVals, chainID)
	s.Nil(err)
	apk, err := s.syncer.GetAPKForBlock(big.NewInt(11), chainID, 12)
	s.Nil(err)
	s.NotNil(apk)
}

func (s *SyncerDBTestSuite) TestGetAPKForBlockNotExistsBlockErr() {
	chainID := uint8(1)
	startVals := make([]*istanbul.ValidatorData, 3)
	startVals[0] = &istanbul.ValidatorData{Address: common.Address{0x0f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	startVals[1] = &istanbul.ValidatorData{Address: common.Address{0x1f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	startVals[2] = &istanbul.ValidatorData{Address: common.Address{0x2f}, BLSPublicKey: blscrypto.SerializedPublicKey{}}
	err := s.syncer.SetValidatorsForBlock(big.NewInt(12), startVals, chainID)
	s.Nil(err)
	apk, err := s.syncer.GetAPKForBlock(big.NewInt(11), chainID, 13)
	s.NotNil(err)
	s.True(errors.Is(err, ErrNoBlockInStore))
	s.Nil(apk)
}
