package itx_test

import (
	"errors"
	"math/big"
	"strings"
	"testing"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/forwarder"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor/itx"
	mock_forwarder "github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor/itx/mock"
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
)

type GsnForwarderTestSuite struct {
	suite.Suite
	gsnForwarder      *itx.GsnForwarder
	forwarderContract *mock_forwarder.MockForwarderContract
	nonceStore        *mock_forwarder.MockNonceStorer
	kp                *secp256k1.Keypair
}

func TestRunGsnForwarderTestSuite(t *testing.T) {
	suite.Run(t, new(GsnForwarderTestSuite))
}

func (s *GsnForwarderTestSuite) SetupSuite()    {}
func (s *GsnForwarderTestSuite) TearDownSuite() {}
func (s *GsnForwarderTestSuite) SetupTest() {
	gomockController := gomock.NewController(s.T())
	s.kp, _ = secp256k1.NewKeypairFromPrivateKey(common.Hex2Bytes("e8e0f5427111dee651e63a6f1029da6929ebf7d2d61cefaf166cebefdf2c012e"))
	s.forwarderContract = mock_forwarder.NewMockForwarderContract(gomockController)
	s.nonceStore = mock_forwarder.NewMockNonceStorer(gomockController)
	s.gsnForwarder = itx.NewGsnForwarder(big.NewInt(5), s.kp, s.forwarderContract, s.nonceStore)
}
func (s *GsnForwarderTestSuite) TearDownTest() {}

func (s *GsnForwarderTestSuite) TestChainID() {
	chainID := s.gsnForwarder.ChainId()

	s.Equal(big.NewInt(5), chainID)
}

func (s *GsnForwarderTestSuite) TestForwarderData_ValidData() {
	to := common.HexToAddress("0x04005C8A516292af163b1AFe3D855b9f4f4631B5")
	forwarderAddress := common.HexToAddress("0x5eDF97800a15E23F386785a2D486bA3E43545210")
	s.forwarderContract.EXPECT().ContractAddress().Return(&forwarderAddress)
	s.forwarderContract.EXPECT().ExecuteData(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).DoAndReturn(func(
		forwardReq forwarder.ForwardRequest,
		domainSeparator *[32]byte,
		typeHash *[32]byte,
		suffixData []byte,
		sig []byte,
	) ([]byte, error) {
		a, _ := abi.JSON(strings.NewReader(consts.GsnForwarderABI))
		return a.Pack("execute", forwardReq, domainSeparator, typeHash, suffixData, sig)
	})

	data, err := s.gsnForwarder.ForwarderData(to, []byte{}, transactor.TransactOptions{
		Value:    big.NewInt(0),
		GasLimit: 200000,
		Nonce:    big.NewInt(1),
	})

	expectedForwarderData := "e024dc7f00000000000000000000000000000000000000000000000000000000000000a04197985b310dcad44e73e243d8d416aae7c1b472d440d3dd15194d86ac46b2152510fc5e187085770200b027d9f2cc4b930768f3b2bd81daafb71ffeb53d21c400000000000000000000000000000000000000000000000000000000000001a000000000000000000000000000000000000000000000000000000000000001c00000000000000000000000007d0e20299178a8d0a8e7410726acc8e338119b8600000000000000000000000004005c8a516292af163b1afe3d855b9f4f4631b500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000030d40000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000e0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000004122f2f1b32ed813e77b30a5aee5c45250a2a46f8c85a57e8d4e760f042d00a7b56241793a18b791c4fb633aec7ae8252b377e9567892e4fcbbf8a27bd223ebc331b00000000000000000000000000000000000000000000000000000000000000"
	s.Nil(err)
	s.Equal(common.Bytes2Hex(data), expectedForwarderData)
}

func (s *GsnForwarderTestSuite) TestUnsafeNonce_ErrorFetchingFromStore() {
	s.nonceStore.EXPECT().GetNonce(big.NewInt(5)).Return(nil, errors.New("error"))

	_, err := s.gsnForwarder.UnsafeNonce()

	s.NotNil(err)
}

func (s *GsnForwarderTestSuite) TestNextNonce_ErrorFetchingFromContract() {
	s.nonceStore.EXPECT().GetNonce(big.NewInt(5)).Return(big.NewInt(1), nil)
	s.forwarderContract.EXPECT().GetNonce(common.HexToAddress(s.kp.Address())).Return(nil, errors.New("error"))

	_, err := s.gsnForwarder.UnsafeNonce()

	s.NotNil(err)
}

func (s *GsnForwarderTestSuite) TestNextNonce_ContractNonceHigher() {
	s.nonceStore.EXPECT().GetNonce(big.NewInt(5)).Return(big.NewInt(1), nil)
	s.forwarderContract.EXPECT().GetNonce(common.HexToAddress(s.kp.Address())).Return(big.NewInt(2), nil)

	nonce, err := s.gsnForwarder.UnsafeNonce()

	s.Nil(err)
	s.Equal(nonce, big.NewInt(2))
}

func (s *GsnForwarderTestSuite) TestNextNonce_StoredNonceHigher() {
	s.nonceStore.EXPECT().GetNonce(big.NewInt(5)).Return(big.NewInt(3), nil)
	s.forwarderContract.EXPECT().GetNonce(common.HexToAddress(s.kp.Address())).Return(big.NewInt(2), nil)

	nonce, err := s.gsnForwarder.UnsafeNonce()

	s.Nil(err)
	s.Equal(nonce, big.NewInt(3))
}

func (s *GsnForwarderTestSuite) TestUnsafeIncreaseNonce_NonceIcremented() {
	s.nonceStore.EXPECT().GetNonce(big.NewInt(5)).Return(big.NewInt(3), nil)
	s.forwarderContract.EXPECT().GetNonce(common.HexToAddress(s.kp.Address())).Return(big.NewInt(2), nil)
	nonce1, err := s.gsnForwarder.UnsafeNonce()
	s.Nil(err)
	s.Equal(nonce1, big.NewInt(3))

	s.gsnForwarder.UnsafeIncreaseNonce()
	nonce2, err := s.gsnForwarder.UnsafeNonce()

	s.Nil(err)
	s.Equal(nonce2, big.NewInt(4))
}

func (s *GsnForwarderTestSuite) TestUnlockNonce_FailedStore_NonceUnlocked() {
	s.nonceStore.EXPECT().GetNonce(big.NewInt(5)).Return(big.NewInt(3), nil)
	s.forwarderContract.EXPECT().GetNonce(common.HexToAddress(s.kp.Address())).Return(big.NewInt(2), nil)
	oldNonce, err := s.gsnForwarder.UnsafeNonce()
	s.Nil(err)
	s.Equal(oldNonce, big.NewInt(3))

	s.gsnForwarder.LockNonce()

	s.nonceStore.EXPECT().StoreNonce(big.NewInt(5), big.NewInt(3)).Return(errors.New("error"))
	s.gsnForwarder.UnlockNonce()

	s.gsnForwarder.UnsafeIncreaseNonce()
	nonce, err := s.gsnForwarder.UnsafeNonce()

	s.Nil(err)
	s.Equal(nonce, big.NewInt(4))
}
