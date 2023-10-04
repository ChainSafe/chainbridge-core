package store_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ChainSafe/sygma-core/mock"
	"github.com/ChainSafe/sygma-core/store"
	"github.com/stretchr/testify/suite"
	"github.com/syndtr/goleveldb/leveldb"
	"go.uber.org/mock/gomock"
)

type NonceStoreTestSuite struct {
	suite.Suite
	nonceStore           *store.NonceStore
	keyValueReaderWriter *mock.MockKeyValueReaderWriter
}

func TestRunNonceStoreTestSuite(t *testing.T) {
	suite.Run(t, new(NonceStoreTestSuite))
}

func (s *NonceStoreTestSuite) SetupSuite()    {}
func (s *NonceStoreTestSuite) TearDownSuite() {}
func (s *NonceStoreTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	s.keyValueReaderWriter = mock.NewMockKeyValueReaderWriter(gomockController)
	s.nonceStore = store.NewNonceStore(s.keyValueReaderWriter)
}
func (s *NonceStoreTestSuite) TearDownTest() {}

func (s *NonceStoreTestSuite) TestStoreBlock_FailedStore() {
	key := "chain:1:nonce"
	s.keyValueReaderWriter.EXPECT().SetByKey([]byte(key), []byte{5}).Return(errors.New("error"))

	err := s.nonceStore.StoreNonce(big.NewInt(1), big.NewInt(5))

	s.NotNil(err)
}

func (s *NonceStoreTestSuite) TestStoreBlock_SuccessfulStore() {
	key := "chain:1:nonce"
	s.keyValueReaderWriter.EXPECT().SetByKey([]byte(key), []byte{5}).Return(nil)

	err := s.nonceStore.StoreNonce(big.NewInt(1), big.NewInt(5))

	s.Nil(err)
}

func (s *NonceStoreTestSuite) TestGetNonce_FailedFetch() {
	key := "chain:1:nonce"
	s.keyValueReaderWriter.EXPECT().GetByKey([]byte(key)).Return(nil, errors.New("error"))

	_, err := s.nonceStore.GetNonce(big.NewInt(1))

	s.NotNil(err)
}

func (s *NonceStoreTestSuite) TestGetNonce_NonceNotFound() {
	key := "chain:1:nonce"
	s.keyValueReaderWriter.EXPECT().GetByKey([]byte(key)).Return(nil, leveldb.ErrNotFound)

	nonce, err := s.nonceStore.GetNonce(big.NewInt(1))

	s.Nil(err)
	s.Equal(nonce, big.NewInt(0))
}

func (s *NonceStoreTestSuite) TestGetNonce_SuccessfulFetch() {
	key := "chain:1:nonce"
	s.keyValueReaderWriter.EXPECT().GetByKey([]byte(key)).Return([]byte{5}, nil)

	block, err := s.nonceStore.GetNonce(big.NewInt(1))

	s.Nil(err)
	s.Equal(block, big.NewInt(5))
}
