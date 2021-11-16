package evm

import (
	"context"
	"math/big"
	"time"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/ChainSafe/chainbridge-core/keystore"

	substrateTypes "github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
)

type TestClient interface {
	E2EClient
	LatestBlock() (*big.Int, error)
	FetchEventLogs(ctx context.Context, contractAddress common.Address, event string, startBlock *big.Int, endBlock *big.Int) ([]types.Log, error)
}

func SetupEVM2EVEMTestSuite(fabric1, fabric2 calls.TxFabric, endpoint1, endpoint2 string, adminKey *secp256k1.Keypair) *IntegrationTestSuite {
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
	gasPricer          calls.GasPricer
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
	ethClient, err := evmclient.NewEVMClientFromParams(s.endpoint1, s.adminKey.PrivateKey())
	if err != nil {
		panic(err)
	}
	s.client = ethClient

	ethClient2, err := evmclient.NewEVMClientFromParams(s.endpoint2, s.adminKey.PrivateKey())
	if err != nil {
		panic(err)
	}
	s.client2 = ethClient2

	b, err := ethClient.LatestBlock()
	if err != nil {
		panic(err)
	}

	log.Debug().Msgf("Latest block %s", b.String())

	config, err := PrepareEVME2EEnv(ethClient, s.fabric1, 1, big.NewInt(1), s.adminKey.CommonAddress())
	if err != nil {
		panic(err)
	}

	s.bridgeAddr = config.bridgeAddr
	s.erc20ContractAddr = config.erc20Addr
	s.erc20HandlerAddr = config.erc20HandlerAddr
	s.genericHandlerAddr = config.genericHandlerAddr
	s.assetStoreAddr = config.assetStoreAddr

	_, err = PrepareEVME2EEnv(ethClient2, s.fabric2, 2, big.NewInt(1), s.adminKey.CommonAddress())
	if err != nil {
		panic(err)
	}

	s.gasPricer = evmgaspricer.NewStaticGasPriceDeterminant(s.client, nil)

	s.erc20RID = calls.SliceTo32Bytes(append(common.LeftPadBytes(config.erc20Addr.Bytes(), 31), uint8(0)))
	s.genericRID = calls.SliceTo32Bytes(append(common.LeftPadBytes(config.genericHandlerAddr.Bytes(), 31), uint8(1)))
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
	err = calls.Deposit(s.client, s.fabric1, s.gasPricer, s.bridgeAddr, s.erc20RID, 2, data)
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
	hash, _ := substrateTypes.GetHash(substrateTypes.NewI64(int64(1)))
	data := calls.ConstructGenericDepositData(hash[:])
	err := calls.Deposit(
		s.client,
		s.fabric1,
		s.gasPricer,
		s.bridgeAddr,
		s.genericRID,
		2,
		data,
	)
	s.Nil(err)

	time.Sleep(120 * time.Second)

	// Asset hash sent is stored in centrifuge asset store contract
	exists, err := calls.IsCentrifugeAssetStored(
		s.client2,
		s.assetStoreAddr,
		hash,
	)
	s.Nil(err)
	s.Equal(true, exists)
}
