package evm

import (
	"context"
	"math/big"
	"time"

	substrateTypes "github.com/centrifuge/go-substrate-rpc-client/types"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/consts"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/ChainSafe/chainbridge-core/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
)

type TestClient interface {
	calls.ChainClient
	LatestBlock() (*big.Int, error)
	FetchEventLogs(ctx context.Context, contractAddress common.Address, event string, startBlock *big.Int, endBlock *big.Int) ([]types.Log, error)
}

func PreSetupTestSuite(fabric1, fabric2 calls.TxFabric, endpoint1, endpoint2 string, adminKey *secp256k1.Keypair) *IntegrationTestSuite {
	return &IntegrationTestSuite{
		fabric1:   fabric1,
		fabric2:   fabric2,
		endpoint1: endpoint1,
		endpoint2: endpoint2,
		adminKey:  adminKey,
	}
}

type IntegrationTestSuite struct {
	suite.Suite
	client             TestClient
	client2            TestClient
	bridgeAddr         common.Address
	erc20HandlerAddr   common.Address
	erc20ContractAddr  common.Address
	genericHandlerAddr common.Address
	assetStoreAddr     common.Address
	fabric1            calls.TxFabric
	fabric2            calls.TxFabric
	endpoint1          string
	endpoint2          string
	adminKey           *secp256k1.Keypair
	erc20RID           [32]byte
	genericRID         [32]byte
}

func (s *IntegrationTestSuite) SetupSuite() {
	ethClient, err := evmclient.NewEVMClientFromParams(s.endpoint1, s.adminKey.PrivateKey(), big.NewInt(consts.DefaultGasPrice))
	s.Nil(err)
	s.client = ethClient

	ethClient2, err := evmclient.NewEVMClientFromParams(s.endpoint2, s.adminKey.PrivateKey(), big.NewInt(consts.DefaultGasPrice))
	s.Nil(err)
	s.client2 = ethClient2

	b, err := ethClient.LatestBlock()
	s.Nil(err)

	log.Debug().Msgf("Latest block %s", b.String())

	bridgeAddr, erc20Addr, erc20HandlerAddr, assetStoreAddr, genericHandlerAddr, err := PrepareEVME2EEnv(ethClient, s.fabric1, 1, big.NewInt(1), s.adminKey.CommonAddress())
	s.Nil(err)

	s.bridgeAddr = bridgeAddr
	s.erc20ContractAddr = erc20Addr
	s.erc20HandlerAddr = erc20HandlerAddr
	s.genericHandlerAddr = genericHandlerAddr
	s.assetStoreAddr = assetStoreAddr
	_, _, _, _, _, err = PrepareEVME2EEnv(ethClient2, s.fabric2, 2, big.NewInt(1), s.adminKey.CommonAddress())
	s.Nil(err)

	s.erc20RID = calls.SliceTo32Bytes(append(common.LeftPadBytes(genericHandlerAddr.Bytes(), 31), 1))
	s.genericRID = calls.SliceTo32Bytes(append(common.LeftPadBytes(genericHandlerAddr.Bytes(), 31), 1))
}
func (s *IntegrationTestSuite) TearDownSuite() {}
func (s *IntegrationTestSuite) SetupTest()     {}
func (s *IntegrationTestSuite) TearDownTest()  {}

func (s *IntegrationTestSuite) TestErc20Deposit() {
	dstAddr := keystore.TestKeyRing.EthereumKeys[keystore.BobKey].CommonAddress()
	senderBalBefore, err := calls.GetERC20Balance(s.client, s.erc20ContractAddr, EveKp.CommonAddress())
	s.Nil(err)
	destBalanceBefore, err := calls.GetERC20Balance(s.client2, s.erc20ContractAddr, dstAddr)
	s.Nil(err)

	amountToDeposit := big.NewInt(1000000)
	data := calls.ConstructErc20DepositData(dstAddr.Bytes(), amountToDeposit)
	err = calls.Deposit(s.client, s.fabric1, s.bridgeAddr, s.erc20RID, 2, data)
	s.Nil(err)

	//Wait 120 seconds for relayer vote
	time.Sleep(120 * time.Second)

	senderBalAfter, err := calls.GetERC20Balance(s.client, s.erc20ContractAddr, s.adminKey.CommonAddress())
	s.Nil(err)
	s.Equal(-1, senderBalAfter.Cmp(senderBalBefore))

	destBalanceAfter, err := calls.GetERC20Balance(s.client2, s.erc20ContractAddr, dstAddr)
	s.Nil(err)
	//Balance has increased
	s.Equal(1, destBalanceAfter.Cmp(destBalanceBefore))
}

func (s *IntegrationTestSuite) TestGenericDeposit() {
	b, err := s.client2.LatestBlock()
	s.Nil(err)

	hash, _ := substrateTypes.GetHash(substrateTypes.NewI64(int64(1)))
	data := calls.ConstructGenericDepositData(hash[:])
	err = calls.Deposit(
		s.client,
		s.fabric1,
		s.bridgeAddr,
		s.genericRID,
		2,
		data,
	)
	s.Nil(err)

	time.Sleep(120 * time.Second)

	exists, err := calls.IsCentrifugeAssetStored(
		s.client2,
		s.assetStoreAddr,
		hash,
	)
	s.Nil(err)
	s.Equal(true, exists)
}
