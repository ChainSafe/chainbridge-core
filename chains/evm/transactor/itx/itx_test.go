package itx_test

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/ChainSafe/chainbridge-core/chains/evm/transactor"
	"github.com/ChainSafe/chainbridge-core/chains/evm/transactor/itx"
	mock_itx "github.com/ChainSafe/chainbridge-core/chains/evm/transactor/itx/mock"
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type TransactTestSuite struct {
	suite.Suite
	forwarder   *mock_itx.MockForwarder
	relayCaller *mock_itx.MockRelayCaller
	transactor  *itx.ITXTransactor
	kp          *secp256k1.Keypair
}

func TestRunTransactTestSuite(t *testing.T) {
	suite.Run(t, new(TransactTestSuite))
}

func (s *TransactTestSuite) SetupSuite()    {}
func (s *TransactTestSuite) TearDownSuite() {}
func (s *TransactTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	s.kp, _ = secp256k1.GenerateKeypair()
	s.forwarder = mock_itx.NewMockForwarder(gomockController)
	s.relayCaller = mock_itx.NewMockRelayCaller(gomockController)
	s.transactor = itx.NewITXTransactor(s.relayCaller, s.forwarder, s.kp)
	s.forwarder.EXPECT().ChainId().Return(uint8(5))
}
func (s *TransactTestSuite) TearDownTest() {}

func (s *TransactTestSuite) TestTransact_FailedFetchingForwarderData() {
	to := common.HexToAddress("0x04005C8A516292af163b1AFe3D855b9f4f4631B5")
	data := []byte{}
	opts := transactor.TransactOptions{
		ChainID: 5,
	}
	s.forwarder.EXPECT().ForwarderData(to, data, s.kp, opts).Return(nil, errors.New("error"))

	_, err := s.transactor.Transact(to, data, opts)

	s.NotNil(err)
}

func (s *TransactTestSuite) TestTransact_FailedSendTransaction() {
	to := common.HexToAddress("0x04005C8A516292af163b1AFe3D855b9f4f4631B5")
	data := []byte{}
	opts := transactor.TransactOptions{
		ChainID:  5,
		GasLimit: big.NewInt(200000),
	}
	s.forwarder.EXPECT().ForwarderData(to, data, s.kp, opts).Return([]byte{}, nil)
	s.forwarder.EXPECT().ForwarderAddress().Return(to)
	s.relayCaller.EXPECT().CallContext(
		context.Background(),
		gomock.Any(),
		"relay_sendTransaction",
		gomock.Any(),
		gomock.Any(),
	).Return(errors.New("error"))

	_, err := s.transactor.Transact(to, data, opts)

	s.NotNil(err)
}

func (s *TransactTestSuite) TestTransact_SuccessfulSend() {
	to := common.HexToAddress("0x04005C8A516292af163b1AFe3D855b9f4f4631B5")
	data := []byte{}
	opts := transactor.TransactOptions{
		ChainID:  5,
		GasLimit: big.NewInt(200000),
		Priority: "low",
	}
	s.forwarder.EXPECT().ForwarderData(to, data, s.kp, opts).Return([]byte{}, nil)
	s.forwarder.EXPECT().ForwarderAddress().Return(to)
	s.relayCaller.EXPECT().CallContext(
		context.Background(),
		gomock.Any(),
		"relay_sendTransaction",
		gomock.Any(),
		gomock.Any(),
	).Return(nil)

	hash, err := s.transactor.Transact(to, data, opts)

	s.Nil(err)
	s.NotNil(hash)
}
