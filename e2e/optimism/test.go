package optimism

import (
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/chainbridge-core/chains/optimism/optimismclient"
	"github.com/ChainSafe/chainbridge-core/e2e/evm"
	substrateTypes "github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/centrifuge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/erc20"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/erc721"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmgaspricer"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/local"
	"github.com/ChainSafe/chainbridge-core/crypto/secp256k1"
	"github.com/ChainSafe/chainbridge-core/keystore"

	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/suite"
)

func SetupEVM2OptimismTestSuite(fabric1, fabric2 calls.TxFabric, endpoint1, endpoint2 string, verifierEndpoint string, ethAdminKey *secp256k1.Keypair, optimismAdminKey *secp256k1.Keypair) *IntegrationTestSuite {
	return &IntegrationTestSuite{
		fabric1:          fabric1,
		fabric2:          fabric2,
		endpoint1:        endpoint1,
		endpoint2:        endpoint2,
		verifierEndpoint: verifierEndpoint,
		ethAdminKey:      ethAdminKey,
		optimismAdminKey: optimismAdminKey,
	}
}

type IntegrationTestSuite struct {
	suite.Suite
	client              evm.TestClient
	client2             evm.TestClient
	gasPricer           calls.GasPricer
	bridgeAddr          common.Address
	bridgeAddr2         common.Address
	erc20HandlerAddr    common.Address
	erc20HandlerAddr2   common.Address
	erc20ContractAddr   common.Address
	erc20ContractAddr2  common.Address
	erc721HandlerAddr   common.Address
	erc721HandlerAddr2  common.Address
	erc721ContractAddr  common.Address
	erc721ContractAddr2 common.Address
	genericHandlerAddr  common.Address
	genericHandlerAddr2 common.Address
	assetStoreAddr      common.Address
	assetStoreAddr2     common.Address
	fabric1             calls.TxFabric
	fabric2             calls.TxFabric
	endpoint1           string
	endpoint2           string
	verifierEndpoint    string
	ethAdminKey         *secp256k1.Keypair
	optimismAdminKey    *secp256k1.Keypair
	erc20RID            [32]byte
	erc721RID           [32]byte
	genericRID          [32]byte
}

func (s *IntegrationTestSuite) SetupSuite() {
	ethClient, err := evmclient.NewEVMClientFromParams(s.endpoint1, s.ethAdminKey.PrivateKey())
	if err != nil {
		panic(err)
	}
	s.client = ethClient

	optimismClient, err := optimismclient.NewOptimismClientFromParams(s.endpoint2, s.optimismAdminKey.PrivateKey(), true, s.verifierEndpoint)
	if err != nil {
		panic(err)
	}
	s.client2 = optimismClient

	b, err := ethClient.LatestBlock()
	if err != nil {
		panic(err)
	}

	log.Debug().Msgf("Latest block %s", b.String())

	config, err := local.PrepareLocalEVME2EEnv(ethClient, s.fabric1, 1, big.NewInt(2), s.ethAdminKey.CommonAddress(), nil)
	if err != nil {
		panic(err)
	}

	log.Debug().Msg("Prepared client 1 environment")

	s.bridgeAddr = config.BridgeAddr
	s.erc20ContractAddr = config.Erc20Addr
	s.erc20HandlerAddr = config.Erc20HandlerAddr
	s.erc721ContractAddr = config.Erc721Addr
	s.erc721HandlerAddr = config.Erc721HandlerAddr
	s.genericHandlerAddr = config.GenericHandlerAddr
	s.assetStoreAddr = config.AssetStoreAddr

	s.gasPricer = evmgaspricer.NewStaticGasPriceDeterminant(s.client, nil)

	s.erc20RID = calls.SliceTo32Bytes(append(common.LeftPadBytes(config.Erc20Addr.Bytes(), 31), uint8(0)))
	s.genericRID = calls.SliceTo32Bytes(append(common.LeftPadBytes(config.GenericHandlerAddr.Bytes(), 31), uint8(1)))
	s.erc721RID = calls.SliceTo32Bytes(append(common.LeftPadBytes(config.Erc721Addr.Bytes(), 31), uint8(2)))

	cfg2, err := local.PrepareLocalEVME2EEnv(optimismClient, s.fabric2, 2, big.NewInt(2), s.optimismAdminKey.CommonAddress(), &config)
	if err != nil {
		panic(err)
	}

	log.Debug().Msg("Prepared client 2 environment")

	s.bridgeAddr2 = cfg2.BridgeAddr
	s.erc20ContractAddr2 = cfg2.Erc20Addr
	s.erc20HandlerAddr2 = cfg2.Erc20HandlerAddr
	s.erc721ContractAddr2 = cfg2.Erc721Addr
	s.erc721HandlerAddr2 = cfg2.Erc721HandlerAddr
	s.genericHandlerAddr2 = cfg2.GenericHandlerAddr
	s.assetStoreAddr2 = cfg2.AssetStoreAddr

	log.Debug().Msgf("Integration test suite: %v", s.erc20HandlerAddr.String())
}
func (s *IntegrationTestSuite) TearDownSuite() {}
func (s *IntegrationTestSuite) SetupTest()     {}
func (s *IntegrationTestSuite) TearDownTest()  {}

func (s *IntegrationTestSuite) TestErc20Deposit() {
	dstAddr := keystore.TestKeyRing.EthereumKeys[keystore.BobKey].CommonAddress()

	transactor1 := transactor.NewSignAndSendTransactor(s.fabric1, s.gasPricer, s.client)
	erc20Contract1 := erc20.NewERC20Contract(s.client, s.erc20ContractAddr, transactor1)
	bridgeContract1 := bridge.NewBridgeContract(s.client, s.bridgeAddr, transactor1)

	transactor2 := transactor.NewSignAndSendTransactor(s.fabric2, s.gasPricer, s.client2)
	erc20Contract2 := erc20.NewERC20Contract(s.client2, s.erc20ContractAddr2, transactor2)

	senderBalBefore, err := erc20Contract1.GetBalance(local.EveKp.CommonAddress())
	s.Nil(err)
	destBalanceBefore, err := erc20Contract2.GetBalance(dstAddr)
	s.Nil(err)

	amountToDeposit := big.NewInt(1000000)
	_, err = bridgeContract1.Erc20Deposit(dstAddr, amountToDeposit, s.erc20RID, 2, transactor.TransactOptions{})
	if err != nil {
		return
	}
	s.Nil(err)

	err = evm.WaitForProposalExecuted(s.client2, s.bridgeAddr2)
	s.Nil(err)

	senderBalAfter, err := erc20Contract1.GetBalance(s.ethAdminKey.CommonAddress())
	s.Nil(err)
	s.Equal(-1, senderBalAfter.Cmp(senderBalBefore))

	destBalanceAfter, err := erc20Contract2.GetBalance(dstAddr)
	s.Nil(err)
	//Balance has increased
	s.Equal(1, destBalanceAfter.Cmp(destBalanceBefore))
}

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
	erc721Contract2 := erc721.NewErc721Contract(s.client2, s.erc721ContractAddr2, transactor2)

	// Mint token and give approval
	// This is done here so token only exists on evm1
	_, err := erc721Contract1.Mint(tokenId, metadata, s.ethAdminKey.CommonAddress(), txOptions)
	s.Nil(err, "Mint failed")
	_, err = erc721Contract1.Approve(tokenId, s.erc721HandlerAddr, txOptions)
	s.Nil(err, "Approve failed")

	// Check on evm1 if initial owner is admin
	initialOwner, err := erc721Contract1.Owner(tokenId)
	s.Nil(err)
	s.Equal(initialOwner.String(), s.ethAdminKey.CommonAddress().String())

	// Check on evm2 token doesn't exist
	_, err = erc721Contract2.Owner(tokenId)
	s.Error(err)

	_, err = bridgeContract1.Erc721Deposit(
		tokenId, metadata, dstAddr, s.erc721RID, 2, transactor.TransactOptions{},
	)
	s.Nil(err)

	err = evm.WaitForProposalExecuted(s.client2, s.bridgeAddr2)
	s.Nil(err)
	// Check on evm1 that token is burned
	_, err = erc721Contract1.Owner(tokenId)
	s.Error(err)

	// Check on evm2 that token is minted to destination address
	owner, err := erc721Contract2.Owner(tokenId)
	s.Nil(err)
	s.Equal(dstAddr.String(), owner.String())
}

func (s *IntegrationTestSuite) TestGenericDeposit() {
	transactor1 := transactor.NewSignAndSendTransactor(s.fabric1, s.gasPricer, s.client)
	transactor2 := transactor.NewSignAndSendTransactor(s.fabric2, s.gasPricer, s.client2)

	bridgeContract1 := bridge.NewBridgeContract(s.client, s.bridgeAddr, transactor1)
	assetStoreContract2 := centrifuge.NewAssetStoreContract(s.client2, s.assetStoreAddr2, transactor2)

	hash, _ := substrateTypes.GetHash(substrateTypes.NewI64(int64(1)))

	_, err := bridgeContract1.GenericDeposit(hash[:], s.genericRID, 2, transactor.TransactOptions{})
	if err != nil {
		return
	}
	s.Nil(err)

	err = evm.WaitForProposalExecuted(s.client2, s.bridgeAddr2)
	s.Nil(err)
	// Asset hash sent is stored in centrifuge asset store contract
	exists, err := assetStoreContract2.IsCentrifugeAssetStored(hash)
	s.Nil(err)
	s.Equal(true, exists)
}
