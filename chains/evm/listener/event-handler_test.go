package listener_test

import (
	"errors"
	"testing"

	"github.com/ChainSafe/chainbridge-core/chains/evm/listener"
	mock_listener "github.com/ChainSafe/chainbridge-core/chains/evm/listener/mock"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

var (
	missingGetDepositRecordAbi = "[{\"inputs\":[],\"stateMutability\":\"nonpayable\",\"type\":\"constructor\"},{\"inputs\":[{\"internalType\":\"uint8\",\"name\":\"\",\"type\":\"uint8\"},{\"internalType\":\"uint64\",\"name\":\"\",\"type\":\"uint64\"}],\"name\":\"_depositRecords\",\"outputs\":[{\"internalType\":\"uint8\",\"name\":\"_destinationChainID\",\"type\":\"uint8\"},{\"internalType\":\"address\",\"name\":\"_depositer\",\"type\":\"address\"},{\"internalType\":\"bytes32\",\"name\":\"_resourceID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_metaData\",\"type\":\"bytes\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
	validAbi                   = "[{\"inputs\":[{\"internalType\":\"uint64\",\"name\":\"depositNonce\",\"type\":\"uint64\"},{\"internalType\":\"uint8\",\"name\":\"destId\",\"type\":\"uint8\"}],\"name\":\"getDepositRecord\",\"outputs\":[{\"components\":[{\"internalType\":\"address\",\"name\":\"_tokenAddress\",\"type\":\"address\"},{\"internalType\":\"uint8\",\"name\":\"_lenDestinationRecipientAddress\",\"type\":\"uint8\"},{\"internalType\":\"uint8\",\"name\":\"_destinationDomainID\",\"type\":\"uint8\"},{\"internalType\":\"bytes32\",\"name\":\"_resourceID\",\"type\":\"bytes32\"},{\"internalType\":\"bytes\",\"name\":\"_destinationRecipientAddress\",\"type\":\"bytes\"},{\"internalType\":\"address\",\"name\":\"_depositer\",\"type\":\"address\"},{\"internalType\":\"uint256\",\"name\":\"_amount\",\"type\":\"uint256\"}],\"internalType\":\"structERC20Handler.DepositRecord\",\"name\":\"\",\"type\":\"tuple\"}],\"stateMutability\":\"view\",\"type\":\"function\"}]"
)

type ListenerTestSuite struct {
	suite.Suite
	clientMock       *mock_listener.MockChainClient
	gomockController *gomock.Controller
}

func TestRunTestSuite(t *testing.T) {
	suite.Run(t, new(ListenerTestSuite))
}

func (s *ListenerTestSuite) SetupSuite()    {}
func (s *ListenerTestSuite) TearDownSuite() {}
func (s *ListenerTestSuite) SetupTest() {
	s.gomockController = gomock.NewController(s.T())
	s.clientMock = mock_listener.NewMockChainClient(s.gomockController)
}
func (s *ListenerTestSuite) TearDownTest() {}

func (s *ListenerTestSuite) TestGetDepositRecordInvalidDefinitionFormat() {
	_, err := listener.GetDepositRecord("invalid", 1, 1, &common.Address{}, s.clientMock)

	s.NotNil(err)
}

func (s *ListenerTestSuite) TestGetDepositRecordMissingGetDepositRecordMethod() {
	_, err := listener.GetDepositRecord(missingGetDepositRecordAbi, 1, 1, &common.Address{}, s.clientMock)

	s.NotNil(err)
}

func (s *ListenerTestSuite) TestGetDepositRecordFailedCallMsg() {
	s.clientMock.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte{}, errors.New("error"))

	_, err := listener.GetDepositRecord(validAbi, 1, 1, &common.Address{}, s.clientMock)

	s.NotNil(err)
	s.Equal(err.Error(), "error")
}

func (s *ListenerTestSuite) TestGetDepositRecordFailedUnpack() {
	s.clientMock.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Return([]byte{}, nil)

	_, err := listener.GetDepositRecord(validAbi, 1, 1, &common.Address{}, s.clientMock)

	s.NotNil(err)
}

func (s *ListenerTestSuite) TestGetDepositRecordSuccess() {
	validCallContractResponse := []byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 32, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 104, 41, 208, 218, 159, 3, 51, 135, 3, 227, 217, 252, 173, 109, 211, 8, 32, 47, 175, 70, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 227, 137, 214, 28, 17, 229, 254, 50, 236, 23, 53, 179, 205, 56, 198, 149, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 224, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 119, 51, 173, 63, 165, 58, 45, 96, 40, 16, 97, 10, 199, 250, 154, 194, 197, 8, 239, 173, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 13, 224, 182, 179, 167, 100, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 20, 119, 51, 173, 63, 165, 58, 45, 96, 40, 16, 97, 10, 199, 250, 154, 194, 197, 8, 239, 173, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	s.clientMock.EXPECT().CallContract(gomock.Any(), gomock.Any(), gomock.Any()).Return(validCallContractResponse, nil)

	rec, err := listener.GetDepositRecord(validAbi, 1, 1, &common.Address{}, s.clientMock)

	s.Nil(err)
	s.NotNil(rec)
}
