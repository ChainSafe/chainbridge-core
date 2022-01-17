package evm

import (
	"context"
	"math/big"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor/signAndSend"
	substrateTypes "github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/centrifuge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/erc20"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/erc721"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmgaspricer"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/local"
	"github.com/ChainSafe/chainbridge-core/keystore"
	"github.com/ethereum/go-ethereum"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/suite"
)

type TestClient interface {
	local.E2EClient
	LatestBlock() (*big.Int, error)
	FetchEventLogs(ctx context.Context, contractAddress common.Address, event string, startBlock *big.Int, endBlock *big.Int) ([]types.Log, error)
	SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
}

func SetupEVM2EVMTestSuite(fabric1, fabric2 calls.TxFabric, client1, client2 TestClient, relayerAddresses1, relayerAddresses2 []common.Address) *IntegrationTestSuite {
	return &IntegrationTestSuite{
		fabric1:           fabric1,
		fabric2:           fabric2,
		client1:           client1,
		client2:           client2,
		relayerAddresses1: relayerAddresses1,
		relayerAddresses2: relayerAddresses2,
	}
}

type IntegrationTestSuite struct {
	suite.Suite
	relayerAddresses1 []common.Address
	relayerAddresses2 []common.Address
	client1           TestClient
	client2           TestClient
	gasPricer1        calls.GasPricer
	gasPricer2        calls.GasPricer
	fabric1           calls.TxFabric
	fabric2           calls.TxFabric
	erc20RID          [32]byte
	erc721RID         [32]byte
	genericRID        [32]byte
	config1           local.EVME2EConfig
	config2           local.EVME2EConfig
}

func (s *IntegrationTestSuite) SetupSuite() {
	config1, err := local.PrepareLocalEVME2EEnv(s.client1, s.fabric1, 1, big.NewInt(2), s.client1.From(), s.relayerAddresses1)
	if err != nil {
		panic(err)
	}
	s.config1 = config1

	config2, err := local.PrepareLocalEVME2EEnv(s.client2, s.fabric2, 2, big.NewInt(2), s.client2.From(), s.relayerAddresses2)
	if err != nil {
		panic(err)
	}
	s.config2 = config2

	s.erc20RID = calls.SliceTo32Bytes(common.LeftPadBytes([]byte{0}, 31))
	s.genericRID = calls.SliceTo32Bytes(common.LeftPadBytes([]byte{1}, 31))
	s.erc721RID = calls.SliceTo32Bytes(common.LeftPadBytes([]byte{2}, 31))
	s.gasPricer1 = evmgaspricer.NewStaticGasPriceDeterminant(s.client1, nil)
	s.gasPricer2 = evmgaspricer.NewStaticGasPriceDeterminant(s.client2, nil)
}
func (s *IntegrationTestSuite) TearDownSuite() {}
func (s *IntegrationTestSuite) SetupTest()     {}
func (s *IntegrationTestSuite) TearDownTest()  {}

func (s *IntegrationTestSuite) TestErc20Deposit() {
	dstAddr := keystore.TestKeyRing.EthereumKeys[keystore.BobKey].CommonAddress()

	transactor1 := signAndSend.NewSignAndSendTransactor(s.fabric1, s.gasPricer1, s.client1)
	erc20Contract1 := erc20.NewERC20Contract(s.client1, s.config1.Erc20Addr, transactor1)
	bridgeContract1 := bridge.NewBridgeContract(s.client1, s.config1.BridgeAddr, transactor1)

	transactor2 := signAndSend.NewSignAndSendTransactor(s.fabric2, s.gasPricer2, s.client2)
	erc20Contract2 := erc20.NewERC20Contract(s.client2, s.config2.Erc20Addr, transactor2)

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

	err = WaitForProposalExecuted(s.client2, s.config2.BridgeAddr)
	s.Nil(err)

	senderBalAfter, err := erc20Contract1.GetBalance(s.client1.From())
	s.Nil(err)
	s.Equal(-1, senderBalAfter.Cmp(senderBalBefore))

	destBalanceAfter, err := erc20Contract2.GetBalance(dstAddr)
	s.Nil(err)
	//Balance has increased
	s.Equal(1, destBalanceAfter.Cmp(destBalanceBefore))
}

func (s *IntegrationTestSuite) TestErc721Deposit() {
	tokenId := big.NewInt(1)
	metadata := "metadata.url"

	dstAddr := keystore.TestKeyRing.EthereumKeys[keystore.BobKey].CommonAddress()

	txOptions := transactor.TransactOptions{}

	// erc721 contract for evm1
	transactor1 := signAndSend.NewSignAndSendTransactor(s.fabric1, s.gasPricer1, s.client1)
	erc721Contract1 := erc721.NewErc721Contract(s.client1, s.config1.Erc721Addr, transactor1)
	bridgeContract1 := bridge.NewBridgeContract(s.client1, s.config1.BridgeAddr, transactor1)

	// erc721 contract for evm2
	transactor2 := signAndSend.NewSignAndSendTransactor(s.fabric2, s.gasPricer2, s.client2)
	erc721Contract2 := erc721.NewErc721Contract(s.client2, s.config2.Erc721Addr, transactor2)

	// Mint token and give approval
	// This is done here so token only exists on evm1
	_, err := erc721Contract1.Mint(tokenId, metadata, s.client1.From(), txOptions)
	s.Nil(err, "Mint failed")
	_, err = erc721Contract1.Approve(tokenId, s.config1.Erc721HandlerAddr, txOptions)
	s.Nil(err, "Approve failed")

	// Check on evm1 if initial owner is admin
	initialOwner, err := erc721Contract1.Owner(tokenId)
	s.Nil(err)
	s.Equal(initialOwner.String(), s.client1.From().String())

	// Check on evm2 token doesn't exist
	_, err = erc721Contract2.Owner(tokenId)
	s.Error(err)

	_, err = bridgeContract1.Erc721Deposit(
		tokenId, metadata, dstAddr, s.erc721RID, 2, transactor.TransactOptions{},
	)
	s.Nil(err)

	err = WaitForProposalExecuted(s.client2, s.config2.BridgeAddr)
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
	transactor1 := signAndSend.NewSignAndSendTransactor(s.fabric1, s.gasPricer1, s.client1)
	transactor2 := signAndSend.NewSignAndSendTransactor(s.fabric2, s.gasPricer2, s.client2)

	bridgeContract1 := bridge.NewBridgeContract(s.client1, s.config1.BridgeAddr, transactor1)
	assetStoreContract2 := centrifuge.NewAssetStoreContract(s.client2, s.config2.AssetStoreAddr, transactor2)

	hash, _ := substrateTypes.GetHash(substrateTypes.NewI64(int64(1)))

	_, err := bridgeContract1.GenericDeposit(hash[:], s.genericRID, 2, transactor.TransactOptions{})
	if err != nil {
		return
	}
	s.Nil(err)

	err = WaitForProposalExecuted(s.client2, s.config2.BridgeAddr)
	s.Nil(err)
	// Asset hash sent is stored in centrifuge asset store contract
	exists, err := assetStoreContract2.IsCentrifugeAssetStored(hash)
	s.Nil(err)
	s.Equal(true, exists)
}
