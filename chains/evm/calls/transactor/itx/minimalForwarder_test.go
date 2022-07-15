package itx_test

import (
	"errors"
	"math/big"
	"strings"
	"testing"

	"github.com/ChainSafe/sygma-core/chains/evm/calls/consts"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/contracts/forwarder"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/sygma-core/chains/evm/calls/transactor/itx"
	mock_forwarder "github.com/ChainSafe/sygma-core/chains/evm/calls/transactor/itx/mock"
	"github.com/ChainSafe/sygma-core/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type MinimalForwarderTestSuite struct {
	suite.Suite
	minimalForwarder  *itx.MinimalForwarder
	forwarderContract *mock_forwarder.MockForwarderContract
	nonceStore        *mock_forwarder.MockNonceStorer
	kp                *secp256k1.Keypair
}

func TestRunMinimalForwarderTestSuite(t *testing.T) {
	suite.Run(t, new(MinimalForwarderTestSuite))
}

func (s *MinimalForwarderTestSuite) SetupSuite()    {}
func (s *MinimalForwarderTestSuite) TearDownSuite() {}
func (s *MinimalForwarderTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	s.kp, _ = secp256k1.NewKeypairFromPrivateKey(common.Hex2Bytes("e8e0f5427111dee651e63a6f1029da6929ebf7d2d61cefaf166cebefdf2c012e"))
	s.forwarderContract = mock_forwarder.NewMockForwarderContract(gomockController)
	s.nonceStore = mock_forwarder.NewMockNonceStorer(gomockController)
	s.minimalForwarder = itx.NewMinimalForwarder(big.NewInt(5), s.kp, s.forwarderContract, s.nonceStore)
}
func (s *MinimalForwarderTestSuite) TearDownTest() {}

func (s *MinimalForwarderTestSuite) TestChainID() {
	chainID := s.minimalForwarder.ChainId()

	s.Equal(big.NewInt(5), chainID)
}

func (s *MinimalForwarderTestSuite) TestForwarderData_ValidData() {
	to := common.HexToAddress("0x04005C8A516292af163b1AFe3D855b9f4f4631B5")
	forwarderAddress := common.HexToAddress("0x5eDF97800a15E23F386785a2D486bA3E43545210")
	s.forwarderContract.EXPECT().ContractAddress().Return(&forwarderAddress)
	s.forwarderContract.EXPECT().PrepareExecute(gomock.Any(), gomock.Any()).DoAndReturn(func(
		forwardReq forwarder.ForwardRequest,
		sig []byte,
	) ([]byte, error) {
		a, _ := abi.JSON(strings.NewReader(consts.MinimalForwarderABI))
		return a.Pack("execute", forwardReq, sig)
	})

	data, err := s.minimalForwarder.ForwarderData(&to, []byte{}, transactor.TransactOptions{
		Value:    big.NewInt(0),
		GasLimit: 200000,
		Nonce:    big.NewInt(1),
	})

	expectedForwarderData := "47153f82000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000001200000000000000000000000007d0e20299178a8d0a8e7410726acc8e338119b8600000000000000000000000004005c8a516292af163b1afe3d855b9f4f4631b500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000030d40000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000410bdb633d2bf0583c34749616155595791c7b7ace433ffcbb1dad606f31a806ed1f8e783b71cd3253a44766db37956bc36cef943707a0f0ab567f1620f145307e1b00000000000000000000000000000000000000000000000000000000000000"
	s.Nil(err)
	s.Equal(common.Bytes2Hex(data), expectedForwarderData)
}

func (s *MinimalForwarderTestSuite) TestUnsafeNonce_ErrorFetchingFromStore() {
	s.nonceStore.EXPECT().GetNonce(big.NewInt(5)).Return(nil, errors.New("error"))

	_, err := s.minimalForwarder.UnsafeNonce()

	s.NotNil(err)
}

func (s *MinimalForwarderTestSuite) TestNextNonce_ErrorFetchingFromContract() {
	s.nonceStore.EXPECT().GetNonce(big.NewInt(5)).Return(big.NewInt(1), nil)
	s.forwarderContract.EXPECT().GetNonce(common.HexToAddress(s.kp.Address())).Return(nil, errors.New("error"))

	_, err := s.minimalForwarder.UnsafeNonce()

	s.NotNil(err)
}

func (s *MinimalForwarderTestSuite) TestNextNonce_ContractNonceHigher() {
	s.nonceStore.EXPECT().GetNonce(big.NewInt(5)).Return(big.NewInt(1), nil)
	s.forwarderContract.EXPECT().GetNonce(common.HexToAddress(s.kp.Address())).Return(big.NewInt(2), nil)

	nonce, err := s.minimalForwarder.UnsafeNonce()

	s.Nil(err)
	s.Equal(nonce, big.NewInt(2))
}

func (s *MinimalForwarderTestSuite) TestNextNonce_StoredNonceHigher() {
	s.nonceStore.EXPECT().GetNonce(big.NewInt(5)).Return(big.NewInt(3), nil)
	s.forwarderContract.EXPECT().GetNonce(common.HexToAddress(s.kp.Address())).Return(big.NewInt(2), nil)

	nonce, err := s.minimalForwarder.UnsafeNonce()

	s.Nil(err)
	s.Equal(nonce, big.NewInt(3))
}

func (s *MinimalForwarderTestSuite) TestUnsafeIncreaseNonce_NonceIcremented() {
	s.nonceStore.EXPECT().GetNonce(big.NewInt(5)).Return(big.NewInt(3), nil)
	s.forwarderContract.EXPECT().GetNonce(common.HexToAddress(s.kp.Address())).Return(big.NewInt(2), nil)
	nonce1, err := s.minimalForwarder.UnsafeNonce()
	s.Nil(err)
	s.Equal(nonce1, big.NewInt(3))

	s.minimalForwarder.UnsafeIncreaseNonce()
	nonce2, err := s.minimalForwarder.UnsafeNonce()

	s.Nil(err)
	s.Equal(nonce2, big.NewInt(4))
}

func (s *MinimalForwarderTestSuite) TestUnlockNonce_FailedStore_NonceUnlocked() {
	s.nonceStore.EXPECT().GetNonce(big.NewInt(5)).Return(big.NewInt(3), nil)
	s.forwarderContract.EXPECT().GetNonce(common.HexToAddress(s.kp.Address())).Return(big.NewInt(2), nil)
	oldNonce, err := s.minimalForwarder.UnsafeNonce()
	s.Nil(err)
	s.Equal(oldNonce, big.NewInt(3))

	s.minimalForwarder.LockNonce()

	s.nonceStore.EXPECT().StoreNonce(big.NewInt(5), big.NewInt(3)).Return(errors.New("error"))
	s.minimalForwarder.UnlockNonce()

	s.minimalForwarder.UnsafeIncreaseNonce()
	nonce, err := s.minimalForwarder.UnsafeNonce()

	s.Nil(err)
	s.Equal(nonce, big.NewInt(4))
}
