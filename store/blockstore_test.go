package store_test

import (
	"errors"
	"math/big"
	"testing"

	"github.com/ChainSafe/sygma-core/store"
	mock_store "github.com/ChainSafe/sygma-core/store/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/syndtr/goleveldb/leveldb"
)

type BlockStoreTestSuite struct {
	suite.Suite
	blockStore           *store.BlockStore
	keyValueReaderWriter *mock_store.MockKeyValueReaderWriter
}

func TestRunBlockStoreTestSuite(t *testing.T) {
	suite.Run(t, new(BlockStoreTestSuite))
}

func (s *BlockStoreTestSuite) SetupSuite()    {}
func (s *BlockStoreTestSuite) TearDownSuite() {}
func (s *BlockStoreTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	s.keyValueReaderWriter = mock_store.NewMockKeyValueReaderWriter(gomockController)
	s.blockStore = store.NewBlockStore(s.keyValueReaderWriter)
}
func (s *BlockStoreTestSuite) TearDownTest() {}

func (s *BlockStoreTestSuite) TestStoreBlock_FailedStore() {
	key := "chain:5:block"
	s.keyValueReaderWriter.EXPECT().SetByKey([]byte(key), []byte{1}).Return(errors.New("error"))

	err := s.blockStore.StoreBlock(big.NewInt(1), 5)

	s.NotNil(err)
}

func (s *BlockStoreTestSuite) TestStoreBlock_SuccessfulStore() {
	key := "chain:5:block"
	s.keyValueReaderWriter.EXPECT().SetByKey([]byte(key), []byte{1}).Return(nil)

	err := s.blockStore.StoreBlock(big.NewInt(1), 5)

	s.Nil(err)
}

func (s *BlockStoreTestSuite) TestGetLastStoredBlock_FailedFetch() {
	key := "chain:5:block"
	s.keyValueReaderWriter.EXPECT().GetByKey([]byte(key)).Return(nil, errors.New("error"))

	_, err := s.blockStore.GetLastStoredBlock(5)

	s.NotNil(err)
}

func (s *BlockStoreTestSuite) TestGetLastStoredBlock_BlockNotFound() {
	key := "chain:5:block"
	s.keyValueReaderWriter.EXPECT().GetByKey([]byte(key)).Return(nil, leveldb.ErrNotFound)

	block, err := s.blockStore.GetLastStoredBlock(5)

	s.Nil(err)
	s.Equal(block, big.NewInt(0))
}

func (s *BlockStoreTestSuite) TestGetLastStoredBlock_SuccessfulFetch() {
	key := "chain:5:block"
	s.keyValueReaderWriter.EXPECT().GetByKey([]byte(key)).Return([]byte{5}, nil)

	block, err := s.blockStore.GetLastStoredBlock(5)

	s.Nil(err)
	s.Equal(block, big.NewInt(5))
}

func (s *BlockStoreTestSuite) TestGetStartBlock_Latest() {
	block, err := s.blockStore.GetStartBlock(5, big.NewInt(1), true, false)

	s.Nil(err)
	s.Nil(block)
}

func (s *BlockStoreTestSuite) TestGetStartBlock_Fresh() {
	block, err := s.blockStore.GetStartBlock(5, big.NewInt(1), false, true)

	s.Nil(err)
	s.Equal(block, big.NewInt(1))
}

func (s *BlockStoreTestSuite) TestGetStartBlock_FailedFetch() {
	key := "chain:5:block"
	s.keyValueReaderWriter.EXPECT().GetByKey([]byte(key)).Return(nil, errors.New("error"))

	_, err := s.blockStore.GetStartBlock(5, big.NewInt(1), false, false)

	s.NotNil(err)
}

func (s *BlockStoreTestSuite) TestGetStartBlock_StartBlockGtLastStoredBlock() {
	key := "chain:5:block"
	s.keyValueReaderWriter.EXPECT().GetByKey([]byte(key)).Return([]byte{5}, nil)

	block, err := s.blockStore.GetStartBlock(5, big.NewInt(10), false, false)

	s.Nil(err)
	s.Equal(block, big.NewInt(10))
}

func (s *BlockStoreTestSuite) TestGetStartBlock_StartBlockLtLastStoredBlock() {
	key := "chain:5:block"
	s.keyValueReaderWriter.EXPECT().GetByKey([]byte(key)).Return([]byte{5}, nil)

	block, err := s.blockStore.GetStartBlock(5, big.NewInt(2), false, false)

	s.Nil(err)
	s.Equal(block, big.NewInt(5))
}
