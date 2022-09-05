package evm_test

import (
	"context"
	"math/big"
	"testing"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmclient"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/evmtransaction"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/transactor/signAndSend"
	"github.com/ChainSafe/chainbridge-core/e2e/dummy"
	"github.com/ChainSafe/chainbridge-core/e2e/evm"
	substrateTypes "github.com/centrifuge/go-substrate-rpc-client/types"
	"github.com/ethereum/go-ethereum/common"

	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/bridge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/centrifuge"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/erc20"
	"github.com/ChainSafe/chainbridge-core/chains/evm/calls/contracts/erc721"
	"github.com/ChainSafe/chainbridge-core/chains/evm/cli/local"
	"github.com/ChainSafe/chainbridge-core/keystore"
	"github.com/ethereum/go-ethereum"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/stretchr/testify/suite"
)

type TestClient interface {
	local.EVMClient
	LatestBlock() (*big.Int, error)
	CodeAt(ctx context.Context, contractAddress common.Address, block *big.Int) ([]byte, error)
	FetchEventLogs(ctx context.Context, contractAddress common.Address, event string, startBlock *big.Int, endBlock *big.Int) ([]types.Log, error)
	SubscribeFilterLogs(ctx context.Context, q ethereum.FilterQuery, ch chan<- types.Log) (ethereum.Subscription, error)
	TransactionByHash(ctx context.Context, hash common.Hash) (tx *types.Transaction, isPending bool, err error)
}

const ETHEndpoint1 = "ws://localhost:8546"
const ETHEndpoint2 = "ws://localhost:8548"

// Alice key is used by the relayer, Eve key is used as admin and depositter
func Test_EVM2EVM(t *testing.T) {
	config := local.BridgeConfig{
		BridgeAddr: common.HexToAddress("0xd606A00c1A39dA53EA7Bb3Ab570BBE40b156EB66"),

		Erc20Addr:        common.HexToAddress("0x75dF75bcdCa8eA2360c562b4aaDBAF3dfAf5b19b"),
		Erc20HandlerAddr: common.HexToAddress("0xb83065680e6AEc805774d8545516dF4e936F0dC0"),
		Erc20ResourceID:  calls.SliceTo32Bytes(common.LeftPadBytes([]byte{0}, 31)),

		Erc721HandlerAddr: common.HexToAddress("0x05C5AFACf64A6082D4933752FfB447AED63581b1"),
		Erc721Addr:        common.HexToAddress("0xb911DF90bCccd3D76a1d8f5fDcd32471e28Cc2c1"),
		Erc721ResourceID:  calls.SliceTo32Bytes(common.LeftPadBytes([]byte{2}, 31)),

		GenericHandlerAddr: common.HexToAddress("0x7573B1c6de00a73e98CDac5Cd2c4a252BdC87600"),
		GenericResourceID:  calls.SliceTo32Bytes(common.LeftPadBytes([]byte{1}, 31)),
		AssetStoreAddr:     common.HexToAddress("0x3cA3808176Ad060Ad80c4e08F30d85973Ef1d99e"),
	}

	ethClient1, err := evmclient.NewEVMClient(ETHEndpoint1, local.EveKp.PrivateKey())
	if err != nil {
		panic(err)
	}
	gasPricer1 := dummy.NewStaticGasPriceDeterminant(ethClient1, nil)

	ethClient2, err := evmclient.NewEVMClient(ETHEndpoint2, local.EveKp.PrivateKey())
	if err != nil {
		panic(err)
	}
	gasPricer2 := dummy.NewStaticGasPriceDeterminant(ethClient2, nil)

	suite.Run(
		t,
		NewEVM2EVMTestSuite(
			evmtransaction.NewTransaction,
			evmtransaction.NewTransaction,
			ethClient1,
			ethClient2,
			gasPricer1,
			gasPricer2,
			config,
			config,
		),
	)
}

func NewEVM2EVMTestSuite(
	fabric1, fabric2 calls.TxFabric,
	client1, client2 TestClient,
	gasPricer1, gasPricer2 calls.GasPricer,
	config1, config2 local.BridgeConfig,
) *IntegrationTestSuite {
	return &IntegrationTestSuite{
		fabric1:    fabric1,
		fabric2:    fabric2,
		client1:    client1,
		client2:    client2,
		gasPricer1: gasPricer1,
		gasPricer2: gasPricer2,
		config1:    config1,
		config2:    config2,
	}
}

type IntegrationTestSuite struct {
	suite.Suite
	client1    TestClient
	client2    TestClient
	gasPricer1 calls.GasPricer
	gasPricer2 calls.GasPricer
	fabric1    calls.TxFabric
	fabric2    calls.TxFabric
	config1    local.BridgeConfig
	config2    local.BridgeConfig
}

// SetupSuite waits until all contracts are deployed
func (s *IntegrationTestSuite) SetupSuite() {
	err := evm.WaitUntilBridgeReady(s.client2, s.config2.BridgeAddr)
	if err != nil {
		panic(err)
	}
}

func (s *IntegrationTestSuite) Test_Erc20Deposit() {
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

	depositTxHash, err := bridgeContract1.Erc20Deposit(dstAddr, amountToDeposit, s.config1.Erc20ResourceID, 2, transactor.TransactOptions{
		Priority: uint8(2), // fast
	})
	s.Nil(err)

	depositTx, _, err := s.client1.TransactionByHash(context.Background(), *depositTxHash)
	s.Nil(err)
	// check gas price of deposit tx - 140 gwei
	s.Equal(big.NewInt(140000000000), depositTx.GasPrice())

	err = evm.WaitUntilProposalExecuted(s.client2, s.config2.BridgeAddr)
	s.Nil(err)

	senderBalAfter, err := erc20Contract1.GetBalance(s.client1.From())
	s.Nil(err)
	s.Equal(-1, senderBalAfter.Cmp(senderBalBefore))

	destBalanceAfter, err := erc20Contract2.GetBalance(dstAddr)
	s.Nil(err)
	//Balance has increased
	s.Equal(1, destBalanceAfter.Cmp(destBalanceBefore))
}

func (s *IntegrationTestSuite) Test_Erc721Deposit() {
	tokenId := big.NewInt(1)
	metadata := "metadata.url"

	dstAddr := keystore.TestKeyRing.EthereumKeys[keystore.BobKey].CommonAddress()

	txOptions := transactor.TransactOptions{
		Priority: uint8(2), // fast
	}

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

	depositTxHash, err := bridgeContract1.Erc721Deposit(
		tokenId, metadata, dstAddr, s.config1.Erc721ResourceID, 2, transactor.TransactOptions{},
	)
	s.Nil(err)

	depositTx, _, err := s.client1.TransactionByHash(context.Background(), *depositTxHash)
	s.Nil(err)
	// check gas price of deposit tx - 50 gwei (slow)
	s.Equal(big.NewInt(50000000000), depositTx.GasPrice())

	err = evm.WaitUntilProposalExecuted(s.client2, s.config2.BridgeAddr)
	s.Nil(err)

	// Check on evm1 that token is burned
	_, err = erc721Contract1.Owner(tokenId)
	s.Error(err)

	// Check on evm2 that token is minted to destination address
	owner, err := erc721Contract2.Owner(tokenId)
	s.Nil(err)
	s.Equal(dstAddr.String(), owner.String())
}

func (s *IntegrationTestSuite) Test_GenericDeposit() {
	transactor1 := signAndSend.NewSignAndSendTransactor(s.fabric1, s.gasPricer1, s.client1)
	transactor2 := signAndSend.NewSignAndSendTransactor(s.fabric2, s.gasPricer2, s.client2)

	bridgeContract1 := bridge.NewBridgeContract(s.client1, s.config1.BridgeAddr, transactor1)
	assetStoreContract2 := centrifuge.NewAssetStoreContract(s.client2, s.config2.AssetStoreAddr, transactor2)

	hash, _ := substrateTypes.GetHash(substrateTypes.NewI64(int64(1)))

	depositTxHash, err := bridgeContract1.GenericDeposit(hash[:], s.config1.GenericResourceID, 2, transactor.TransactOptions{
		Priority: uint8(0), // slow
	})
	s.Nil(err)

	depositTx, _, err := s.client1.TransactionByHash(context.Background(), *depositTxHash)
	s.Nil(err)
	// check gas price of deposit tx - 140 gwei
	s.Equal(big.NewInt(50000000000), depositTx.GasPrice())

	err = evm.WaitUntilProposalExecuted(s.client2, s.config2.BridgeAddr)
	s.Nil(err)
	// Asset hash sent is stored in centrifuge asset store contract
	exists, err := assetStoreContract2.IsCentrifugeAssetStored(hash)
	s.Nil(err)
	s.Equal(true, exists)
}
