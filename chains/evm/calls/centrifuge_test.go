package calls_test

import (
	"errors"
	"testing"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	mock_utils "github.com/ChainSafe/chainbridge-core/chains/evm/calls/mock"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type IsCentrifugeAssetStoredTestSuite struct {
	suite.Suite
	gomockController *gomock.Controller
	clientMock       *mock_utils.MockContractCallerClient
}

func TestRunIsCentrifugeAssetStoredTestSuite(t *testing.T) {
	suite.Run(t, new(IsCentrifugeAssetStoredTestSuite))
}

func (s *IsCentrifugeAssetStoredTestSuite) SetupSuite()    {}
func (s *IsCentrifugeAssetStoredTestSuite) TearDownSuite() {}
func (s *IsCentrifugeAssetStoredTestSuite) SetupTest() {
	s.gomockController = gomock.NewController(s.T())
	s.clientMock = mock_utils.NewMockContractCallerClient(s.gomockController)
}
func (s *IsCentrifugeAssetStoredTestSuite) TearDownTest() {}

func (s *IsCentrifugeAssetStoredTestSuite) TestCallContractFails() {
	s.clientMock.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte{}, errors.New("error"))

	isStored, err := calls.IsCentrifugeAssetStored(s.clientMock, common.Address{}, [32]byte{})

	s.NotNil(err)
	s.Equal(isStored, false)
}

func (s *IsCentrifugeAssetStoredTestSuite) TestUnpackingInvalidOutput() {
	s.clientMock.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte("invalid"), nil)

	isStored, err := calls.IsCentrifugeAssetStored(s.clientMock, common.Address{}, [32]byte{})

	s.Nil(err)
	s.Equal(isStored, false)
}

func (s *IsCentrifugeAssetStoredTestSuite) TestEmptyOutput() {
	s.clientMock.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte{}, nil)

	isStored, err := calls.IsCentrifugeAssetStored(s.clientMock, common.Address{}, [32]byte{})

	s.Nil(err)
	s.Equal(isStored, false)
}

func (s *IsCentrifugeAssetStoredTestSuite) TestValidStoredAsset() {
	response := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1}
	s.clientMock.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Return(response, nil)

	isStored, err := calls.IsCentrifugeAssetStored(s.clientMock, common.Address{}, [32]byte{})

	s.Nil(err)
	s.Equal(isStored, true)
}
