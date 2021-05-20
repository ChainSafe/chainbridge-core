//Copyright 2020 ChainSafe Systems
//SPDX-License-Identifier: LGPL-3.0-only
package validatorsync

import (
	"errors"
	"math/big"
	"os"
	"testing"
	"time"

	mock_validatorsync "github.com/ChainSafe/chainbridge-core/subchains/celo/validatorsync/mock"
	"github.com/celo-org/celo-blockchain/common"
	"github.com/celo-org/celo-blockchain/common/hexutil"
	"github.com/celo-org/celo-blockchain/core/types"
	"github.com/celo-org/celo-blockchain/crypto"
	blscrypto "github.com/celo-org/celo-blockchain/crypto/bls"
	"github.com/celo-org/celo-blockchain/rlp"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/syndtr/goleveldb/leveldb"
)

type SyncTestSuite struct {
	suite.Suite
	store  *ValidatorsStore
	client *mock_validatorsync.MockHeaderByNumberGetter
}

func TestRunSyncTestSuite(t *testing.T) {
	suite.Run(t, new(SyncTestSuite))
}
func (s *SyncTestSuite) SetupSuite()    {}
func (s *SyncTestSuite) TearDownSuite() {}
func (s *SyncTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	db, err := leveldb.OpenFile("./test/db", nil)
	if err != nil {
		s.Fail(err.Error())
	}
	syncer := NewValidatorsStore(db)
	s.store = syncer
	s.client = mock_validatorsync.NewMockHeaderByNumberGetter(gomockController)

}
func (s *SyncTestSuite) TearDownTest() {
	s.store.Close()
	os.RemoveAll("./test")
}

func generateBlockHeader() (*types.Header, error) {
	testKey2, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f292")
	testKey3, _ := crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f293")
	blsPK2, _ := blscrypto.ECDSAToBLS(testKey2)
	blsPK3, _ := blscrypto.ECDSAToBLS(testKey3)
	pubKey2, _ := blscrypto.PrivateToPublic(blsPK2)
	pubKey3, _ := blscrypto.PrivateToPublic(blsPK3)
	extra, err := rlp.EncodeToBytes(&types.IstanbulExtra{
		AddedValidators: []common.Address{
			common.BytesToAddress(hexutil.MustDecode("0x44add0ec310f115a0e603b2d7db9f067778eaf8a")),
			common.BytesToAddress(hexutil.MustDecode("0x294fc7e8f22b3bcdcf955dd7ff3ba2ed833f8212")),
		},
		AddedValidatorsPublicKeys: []blscrypto.SerializedPublicKey{
			pubKey2,
			pubKey3,
		},
		RemovedValidators:    big.NewInt(0),
		Seal:                 []byte{},
		AggregatedSeal:       types.IstanbulAggregatedSeal{},
		ParentAggregatedSeal: types.IstanbulAggregatedSeal{},
	})
	if err != nil {
		return nil, err
	}
	h := &types.Header{
		Extra: append(make([]byte, types.IstanbulExtraVanity), extra...),
	}
	return h, nil
}

func (s *SyncTestSuite) TestStoreBlockValidatorsWIthEmptyDB() {
	header, err := generateBlockHeader()
	s.Nil(err)
	stopChn := make(chan struct{})
	errChn := make(chan error)
	chainID := uint8(1)
	// First iteration
	s.client.EXPECT().HeaderByNumber(gomock.Any(), gomock.Any()).Return(header, nil)
	//Second iteration
	s.client.EXPECT().HeaderByNumber(gomock.Any(), gomock.Any()).Return(header, nil)
	//Third iteration
	s.client.EXPECT().HeaderByNumber(gomock.Any(), gomock.Any()).Return(header, nil)
	// Erroring to stop routine
	e := errors.New("some error occured")
	s.client.EXPECT().HeaderByNumber(gomock.Any(), gomock.Any()).Return(nil, e)

	go func() {
		select {
		case err := <-errChn:
			s.True(errors.Is(err, e))
		case <-time.After(time.Second * 10):
			// Closing this goroutine after 10 seconds
			return
		}
	}()

	SyncBlockValidators(stopChn, errChn, s.client, s.store, chainID, 12)

	vals, err := s.store.GetValidatorsForBlock(big.NewInt(0), chainID)
	s.Nil(err)
	vals2, err := s.store.GetValidatorsForBlock(big.NewInt(12), chainID)
	s.Nil(err)
	vals4, err := s.store.GetValidatorsForBlock(big.NewInt(24), chainID)
	s.Nil(err)
	// zero epoch (2 validators
	s.Equal(2, len(vals))

	// first epoch 1-12 blocs (4 validators)
	s.Equal(4, len(vals2))

	// second epoch 13-24 (6 validators
	s.Equal(6, len(vals4))

	lb, err := s.store.GetLatestKnownEpochLastBlock(chainID)
	s.Nil(err)
	s.Equal(0, lb.Cmp(big.NewInt(24)))

	apk, err := s.store.GetAPKForBlock(big.NewInt(1), chainID, 12)
	s.Nil(err)
	s.NotNil(apk)

}
