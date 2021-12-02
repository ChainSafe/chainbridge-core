package evm

import (
	"context"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/client"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/erc20"
	"math/big"
	"time"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/erc721"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/local"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/evmgaspricer"
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/ChainSafe/chainbridge-core/keystore"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
)

type TestClient interface {
	local.E2EClient
	LatestBlock() (*big.Int, error)
	FetchEventLogs(ctx context.Context, contractAddress common.Address, event string, startBlock *big.Int, endBlock *big.Int) ([]types.Log, error)
}

func SetupEVM2EVMTestSuite(fabric1, fabric2 client.TxFabric, endpoint1, endpoint2 string, adminKey *secp256k1.Keypair) *IntegrationTestSuite {
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
	gasPricer          client.GasPricer
	bridgeAddr         common.Address
	erc20HandlerAddr   common.Address
	erc20ContractAddr  common.Address
	erc721HandlerAddr  common.Address
	erc721ContractAddr common.Address
	genericHandlerAddr common.Address
	assetStoreAddr     common.Address
	fabric1            client.TxFabric
	fabric2            client.TxFabric
	endpoint1          string
	endpoint2          string
	adminKey           *secp256k1.Keypair
	erc20RID           [32]byte
	erc721RID          [32]byte
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

	config, err := local.PrepareLocalEVME2EEnv(ethClient, s.fabric1, 1, big.NewInt(2), s.adminKey.CommonAddress())
	if err != nil {
		panic(err)
	}

	s.bridgeAddr = config.BridgeAddr
	s.erc20ContractAddr = config.Erc20Addr
	s.erc20HandlerAddr = config.Erc20HandlerAddr
	s.erc721ContractAddr = config.Erc721Addr
	s.erc721HandlerAddr = config.Erc721HandlerAddr
	s.genericHandlerAddr = config.GenericHandlerAddr
	s.assetStoreAddr = config.AssetStoreAddr

	s.gasPricer = evmgaspricer.NewStaticGasPriceDeterminant(s.client, nil)

	_, err = local.PrepareLocalEVME2EEnv(ethClient2, s.fabric2, 2, big.NewInt(2), s.adminKey.CommonAddress())
	if err != nil {
		panic(err)
	}

	s.erc20RID = client.SliceTo32Bytes(append(common.LeftPadBytes(config.Erc20Addr.Bytes(), 31), uint8(0)))
	s.genericRID = client.SliceTo32Bytes(append(common.LeftPadBytes(config.GenericHandlerAddr.Bytes(), 31), uint8(1)))
	s.erc721RID = client.SliceTo32Bytes(append(common.LeftPadBytes(config.Erc721Addr.Bytes(), 31), uint8(2)))
}
func (s *IntegrationTestSuite) TearDownSuite() {}
func (s *IntegrationTestSuite) SetupTest()     {}
func (s *IntegrationTestSuite) TearDownTest()  {}

func (s *IntegrationTestSuite) TestErc721Deposit() {
	s.NotEmpty(s.erc721ContractAddr)
	s.NotEmpty(s.erc721HandlerAddr)
	tokenId := big.NewInt(1)
	metadata := "metadata.url"

	dstAddr := keystore.TestKeyRing.EthereumKeys[keystore.BobKey].CommonAddress()

	txOptions := transactor.TransactOptions{}

	// erc721 contract for evm1
	transactor1 := transactor.NewSignAndSendTransactor(s.fabric1, s.gasPricer, s.client)
	erc721Contract1 := erc721.NewErc721Contract(s.client, s.erc721ContractAddr, transactor1)
	bridgeContract1 := bridge.NewBridgeContract(s.client, s.bridgeAddr, transactor1)

	// erc721 contract for evm2
	transactor2 := transactor.NewSignAndSendTransactor(s.fabric2, s.gasPricer, s.client2)
	erc721Contract2 := erc721.NewErc721Contract(s.client2, s.erc721ContractAddr, transactor2)

	// Mint token and give approval
	// This is done here so token only exists on evm1
	_, err := erc721Contract1.Mint(tokenId, metadata, s.adminKey.CommonAddress(), txOptions)
	s.Nil(err, "Mint failed")
	_, err = erc721Contract1.Approve(tokenId, s.erc721HandlerAddr, txOptions)
	s.Nil(err, "Approve failed")

	// Check on evm1 if initial owner is admin
	initialOwner, err := erc721Contract1.Owner(tokenId)
	s.Nil(err)
	s.Equal(initialOwner.String(), s.adminKey.CommonAddress().String())

	// Check on evm2 token doesn't exist
	_, err = erc721Contract2.Owner(tokenId)
	s.Error(err)

	_, err = bridgeContract1.Erc721Deposit(
		tokenId, metadata, dstAddr, s.erc721RID, 2, transactor.TransactOptions{},
	)
	s.Nil(err)

	//Wait 120 seconds for relayer vote
	time.Sleep(120 * time.Second)

	// Check on evm1 that token is burned
	_, err = erc721Contract1.Owner(tokenId)
	s.Error(err)

	// Check on evm2 that token is minted to destination address
	owner, err := erc721Contract2.Owner(tokenId)
	s.Nil(err)
	s.Equal(dstAddr.String(), owner.String())
}

func (s *IntegrationTestSuite) TestErc20Deposit() {
	dstAddr := keystore.TestKeyRing.EthereumKeys[keystore.BobKey].CommonAddress()

	transactor1 := transactor.NewSignAndSendTransactor(s.fabric1, s.gasPricer, s.client)
	erc20Contract1 := erc20.NewERC20Contract(s.client, s.erc20ContractAddr, transactor1)
	bridgeContract1 := bridge.NewBridgeContract(s.client, s.bridgeAddr, transactor1)

	senderBalBefore, err := erc20Contract1.GetBalance(local.EveKp.CommonAddress())
	s.Nil(err)
	destBalanceBefore, err := calls.GetERC20Balance(s.client2, s.erc20ContractAddr, dstAddr)
	s.Nil(err)

	amountToDeposit := big.NewInt(1000000)
	data := calls.ConstructErc20DepositData(dstAddr.Bytes(), amountToDeposit)
	_, err = calls.Deposit(s.client, s.fabric1, s.gasPricer, s.bridgeAddr, s.erc20RID, 2, data)
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

// func (s *IntegrationTestSuite) TestGenericDeposit() {
// 	hash, _ := substrateTypes.GetHash(substrateTypes.NewI64(int64(1)))
// 	data := calls.ConstructGenericDepositData(hash[:])
// 	_, err := calls.Deposit(
// 		s.client,
// 		s.fabric1,
// 		s.gasPricer,
// 		s.bridgeAddr,
// 		s.genericRID,
// 		2,
// 		data,
// 	)
// 	s.Nil(err)

// 	time.Sleep(120 * time.Second)

// 	// Asset hash sent is stored in centrifuge asset store contract
// 	exists, err := calls.IsCentrifugeAssetStored(
// 		s.client2,
// 		s.assetStoreAddr,
// 		hash,
// 	)
// 	s.Nil(err)
// 	s.Equal(true, exists)
// }
