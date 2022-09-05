package centrifuge_test

import (
	"errors"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/centrifuge"
	mock_calls "github.com/ChainSafe/chainbridge-core/chains/evm/calls/mock"
	mock_transactor "github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor/mock"
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type IsCentrifugeAssetStoredTestSuite struct {
	suite.Suite
	gomockController                   *gomock.Controller
	mockContractCallerDispatcherClient *mock_calls.MockContractCallerDispatcher
	mockTransactor                     *mock_transactor.MockTransactor
	assetStoreContractAddress          common.Address
	assetStoreContract                 *centrifuge.AssetStoreContract
}

func TestRunIsCentrifugeAssetStoredTestSuite(t *testing.T) {
	suite.Run(t, new(IsCentrifugeAssetStoredTestSuite))
}

func (s *IsCentrifugeAssetStoredTestSuite) SetupSuite()    {}
func (s *IsCentrifugeAssetStoredTestSuite) TearDownSuite() {}
func (s *IsCentrifugeAssetStoredTestSuite) SetupTest() {
	s.gomockController = gomock.NewController(s.T())
	s.mockContractCallerDispatcherClient = mock_calls.NewMockContractCallerDispatcher(s.gomockController)
	s.mockTransactor = mock_transactor.NewMockTransactor(s.gomockController)
	s.assetStoreContractAddress = common.HexToAddress("0x9A0E6F91E6031C08326764655432f8F9c180fBa0")
	s.assetStoreContract = centrifuge.NewAssetStoreContract(
		s.mockContractCallerDispatcherClient, s.assetStoreContractAddress, s.mockTransactor,
	)
}
func (s *IsCentrifugeAssetStoredTestSuite) TearDownTest() {}

func (s *IsCentrifugeAssetStoredTestSuite) TestCallContractFails() {
	s.mockContractCallerDispatcherClient.EXPECT().CallContract(
		gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte{}, errors.New("error"))
	s.mockContractCallerDispatcherClient.EXPECT().From().Times(1).Return(common.Address{})
	isStored, err := s.assetStoreContract.IsCentrifugeAssetStored([32]byte{})

	s.NotNil(err)
	s.Equal(isStored, false)
}

func (s *IsCentrifugeAssetStoredTestSuite) TestUnpackingInvalidOutput() {
	s.mockContractCallerDispatcherClient.EXPECT().CallContract(
		gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte("invalid"), nil)
	s.mockContractCallerDispatcherClient.EXPECT().From().Times(1).Return(common.Address{})
	isStored, err := s.assetStoreContract.IsCentrifugeAssetStored([32]byte{})

	s.NotNil(err)
	s.Equal(isStored, false)
}

func (s *IsCentrifugeAssetStoredTestSuite) TestEmptyOutput() {
	s.mockContractCallerDispatcherClient.EXPECT().CallContract(
		gomock.Any(), gomock.Any(), gomock.Any(),
	).Return([]byte{}, nil)
	s.mockContractCallerDispatcherClient.EXPECT().CodeAt(
		gomock.Any(), gomock.Any(), gomock.Any(),
	).Return(nil, errors.New("error"))
	s.mockContractCallerDispatcherClient.EXPECT().From().Times(1).Return(common.Address{})

	isStored, err := s.assetStoreContract.IsCentrifugeAssetStored([32]byte{})

	s.NotNil(err)
	s.Equal(isStored, false)
}

func (s *IsCentrifugeAssetStoredTestSuite) TestValidStoredAsset() {
	response := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	s.mockContractCallerDispatcherClient.EXPECT().CallContract(
		gomock.Any(), gomock.Any(), gomock.Any()).Return(response, nil)
	s.mockContractCallerDispatcherClient.EXPECT().From().Times(1).Return(common.Address{})

	isStored, err := s.assetStoreContract.IsCentrifugeAssetStored([32]byte{})

	s.Nil(err)
	s.Equal(isStored, true)
}
